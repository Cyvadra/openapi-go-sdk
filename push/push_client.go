package push

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tigerfintech/openapi-go-sdk/config"
	"github.com/tigerfintech/openapi-go-sdk/signer"
)

const (
	// 默认推送服务器地址
	defaultPushURL = "wss://openapi-push.tigerfintech.com"
	// 默认心跳间隔
	defaultHeartbeatInterval = 10 * time.Second
	// 默认重连间隔
	defaultReconnectInterval = 5 * time.Second
	// 最大重连间隔
	maxReconnectInterval = 60 * time.Second
	// 默认连接超时
	defaultConnectTimeout = 30 * time.Second
)

// ConnectionState 连接状态
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
)

// PushClientOption PushClient 配置选项
type PushClientOption func(*PushClient)

// WithPushURL 设置推送服务器地址
func WithPushURL(url string) PushClientOption {
	return func(c *PushClient) { c.pushURL = url }
}

// WithHeartbeatInterval 设置心跳间隔
func WithHeartbeatInterval(d time.Duration) PushClientOption {
	return func(c *PushClient) { c.heartbeatInterval = d }
}

// WithReconnectInterval 设置初始重连间隔
func WithReconnectInterval(d time.Duration) PushClientOption {
	return func(c *PushClient) { c.reconnectInterval = d }
}

// WithAutoReconnect 设置是否自动重连
func WithAutoReconnect(auto bool) PushClientOption {
	return func(c *PushClient) { c.autoReconnect = auto }
}

// WithConnectTimeout 设置连接超时
func WithConnectTimeout(d time.Duration) PushClientOption {
	return func(c *PushClient) { c.connectTimeout = d }
}

// PushClient WebSocket 推送客户端
type PushClient struct {
	config            *config.ClientConfig
	pushURL           string
	heartbeatInterval time.Duration
	reconnectInterval time.Duration
	connectTimeout    time.Duration
	autoReconnect     bool

	// WebSocket 连接
	conn  *websocket.Conn
	state ConnectionState

	// 回调
	callbacks Callbacks

	// 订阅状态管理
	subscriptions map[SubjectType]map[string]bool // subject -> symbols set
	accountSubs   map[SubjectType]bool            // 账户级别订阅

	// 并发控制
	mu      sync.RWMutex
	stopCh  chan struct{}
	doneCh  chan struct{}
	writeMu sync.Mutex

	// 用于测试的 dialer（可注入）
	dialer WebSocketDialer
}

// WebSocketDialer WebSocket 拨号器接口，方便测试注入
type WebSocketDialer interface {
	Dial(urlStr string, timeout time.Duration) (*websocket.Conn, error)
}

// defaultDialer 默认的 WebSocket 拨号器
type defaultDialer struct{}

func (d *defaultDialer) Dial(urlStr string, timeout time.Duration) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("解析推送服务器地址失败: %w", err)
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("WebSocket 连接失败: %w", err)
	}
	return conn, nil
}

// NewPushClient 创建推送客户端
func NewPushClient(cfg *config.ClientConfig, opts ...PushClientOption) *PushClient {
	c := &PushClient{
		config:            cfg,
		pushURL:           defaultPushURL,
		heartbeatInterval: defaultHeartbeatInterval,
		reconnectInterval: defaultReconnectInterval,
		connectTimeout:    defaultConnectTimeout,
		autoReconnect:     true,
		state:             StateDisconnected,
		subscriptions:     make(map[SubjectType]map[string]bool),
		accountSubs:       make(map[SubjectType]bool),
		dialer:            &defaultDialer{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetCallbacks 设置回调函数集合
func (c *PushClient) SetCallbacks(cb Callbacks) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.callbacks = cb
}

// State 获取当前连接状态
func (c *PushClient) State() ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// Connect 连接到推送服务器并进行认证
func (c *PushClient) Connect() error {
	c.mu.Lock()
	if c.state != StateDisconnected {
		c.mu.Unlock()
		return fmt.Errorf("客户端已连接或正在连接中")
	}
	c.state = StateConnecting
	c.stopCh = make(chan struct{})
	c.doneCh = make(chan struct{})
	c.mu.Unlock()

	// 建立 WebSocket 连接
	conn, err := c.dialer.Dial(c.pushURL, c.connectTimeout)
	if err != nil {
		c.mu.Lock()
		c.state = StateDisconnected
		c.mu.Unlock()
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	// 发送认证消息
	if err := c.authenticate(); err != nil {
		conn.Close()
		c.mu.Lock()
		c.state = StateDisconnected
		c.conn = nil
		c.mu.Unlock()
		return fmt.Errorf("认证失败: %w", err)
	}

	c.mu.Lock()
	c.state = StateConnected
	c.mu.Unlock()

	// 启动消息读取和心跳协程
	go c.readLoop()
	go c.heartbeatLoop()

	// 触发连接成功回调
	c.mu.RLock()
	cb := c.callbacks.OnConnect
	c.mu.RUnlock()
	if cb != nil {
		cb()
	}

	return nil
}

// authenticate 发送认证消息
func (c *PushClient) authenticate() error {
	ts := time.Now().Format("2006-01-02 15:04:05")
	signContent := c.config.TigerID
	sign, err := signer.SignWithRSA(c.config.PrivateKey, signContent)
	if err != nil {
		return fmt.Errorf("签名失败: %w", err)
	}

	req := ConnectRequest{
		TigerID:   c.config.TigerID,
		Sign:      sign,
		Timestamp: ts,
		Version:   "2.0",
	}

	msg, err := NewPushMessage(MsgTypeConnect, "", &req)
	if err != nil {
		return err
	}

	return c.sendMessage(msg)
}

// Disconnect 断开连接
func (c *PushClient) Disconnect() error {
	c.mu.Lock()
	if c.state == StateDisconnected {
		c.mu.Unlock()
		return nil
	}
	c.state = StateDisconnected
	if c.stopCh != nil {
		close(c.stopCh)
	}
	conn := c.conn
	c.conn = nil
	c.mu.Unlock()

	var err error
	if conn != nil {
		err = conn.Close()
	}

	// 等待协程退出
	if c.doneCh != nil {
		select {
		case <-c.doneCh:
		case <-time.After(5 * time.Second):
		}
	}

	// 触发断开回调
	c.mu.RLock()
	cb := c.callbacks.OnDisconnect
	c.mu.RUnlock()
	if cb != nil {
		cb()
	}

	return err
}

// sendMessage 发送消息到 WebSocket
func (c *PushClient) sendMessage(msg *PushMessage) error {
	data, err := msg.Serialize()
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("WebSocket 连接未建立")
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// readLoop 消息读取循环
func (c *PushClient) readLoop() {
	defer func() {
		select {
		case <-c.doneCh:
		default:
			close(c.doneCh)
		}
	}()

	for {
		select {
		case <-c.stopCh:
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			// 检查是否是主动关闭
			select {
			case <-c.stopCh:
				return
			default:
			}

			// 触发错误回调
			c.mu.RLock()
			errCb := c.callbacks.OnError
			c.mu.RUnlock()
			if errCb != nil {
				errCb(err)
			}

			// 尝试自动重连
			c.mu.RLock()
			autoReconnect := c.autoReconnect
			c.mu.RUnlock()
			if autoReconnect {
				go c.reconnect()
			}
			return
		}

		c.handleMessage(data)
	}
}

// heartbeatLoop 心跳保活循环
func (c *PushClient) heartbeatLoop() {
	ticker := time.NewTicker(c.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			msg := &PushMessage{Type: MsgTypeHeartbeat}
			if err := c.sendMessage(msg); err != nil {
				// 心跳发送失败，可能连接已断开
				return
			}
		}
	}
}

// reconnect 自动重连
func (c *PushClient) reconnect() {
	c.mu.Lock()
	c.state = StateDisconnected
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()

	interval := c.reconnectInterval
	for {
		select {
		case <-c.stopCh:
			return
		default:
		}

		time.Sleep(interval)

		if err := c.Connect(); err != nil {
			// 指数退避
			interval = interval * 2
			if interval > maxReconnectInterval {
				interval = maxReconnectInterval
			}
			continue
		}

		// 重连成功，恢复订阅
		c.resubscribe()
		return
	}
}

// resubscribe 重连后恢复之前的订阅
func (c *PushClient) resubscribe() {
	c.mu.RLock()
	subs := make(map[SubjectType][]string)
	for subject, symbols := range c.subscriptions {
		list := make([]string, 0, len(symbols))
		for s := range symbols {
			list = append(list, s)
		}
		subs[subject] = list
	}
	acctSubs := make(map[SubjectType]bool)
	for k, v := range c.accountSubs {
		acctSubs[k] = v
	}
	c.mu.RUnlock()

	// 恢复行情订阅
	for subject, symbols := range subs {
		c.subscribe(subject, symbols, "", "")
	}

	// 恢复账户订阅
	for subject := range acctSubs {
		c.subscribe(subject, nil, c.config.Account, "")
	}
}

// handleMessage 处理收到的消息
func (c *PushClient) handleMessage(data []byte) {
	msg, err := DeserializeMessage(data)
	if err != nil {
		c.mu.RLock()
		errCb := c.callbacks.OnError
		c.mu.RUnlock()
		if errCb != nil {
			errCb(fmt.Errorf("反序列化消息失败: %w", err))
		}
		return
	}

	c.mu.RLock()
	cb := c.callbacks
	c.mu.RUnlock()

	switch msg.Type {
	case MsgTypeKickout:
		if cb.OnKickout != nil {
			var message string
			json.Unmarshal(msg.Data, &message)
			cb.OnKickout(message)
		}
	case MsgTypeError:
		if cb.OnError != nil {
			var message string
			json.Unmarshal(msg.Data, &message)
			cb.OnError(fmt.Errorf("服务端错误: %s", message))
		}
	case MsgTypeQuote:
		if cb.OnQuote != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnQuote(&d)
			}
		}
	case MsgTypeTick:
		if cb.OnTick != nil {
			var d TickData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnTick(&d)
			}
		}
	case MsgTypeDepth:
		if cb.OnDepth != nil {
			var d DepthData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnDepth(&d)
			}
		}
	case MsgTypeOption:
		if cb.OnOption != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnOption(&d)
			}
		}
	case MsgTypeFuture:
		if cb.OnFuture != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnFuture(&d)
			}
		}
	case MsgTypeKline:
		if cb.OnKline != nil {
			var d KlineData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnKline(&d)
			}
		}
	case MsgTypeAsset:
		if cb.OnAsset != nil {
			var d AssetData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnAsset(&d)
			}
		}
	case MsgTypePosition:
		if cb.OnPosition != nil {
			var d PositionData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnPosition(&d)
			}
		}
	case MsgTypeOrder:
		if cb.OnOrder != nil {
			var d OrderData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnOrder(&d)
			}
		}
	case MsgTypeTransaction:
		if cb.OnTransaction != nil {
			var d TransactionData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnTransaction(&d)
			}
		}
	case MsgTypeStockTop:
		if cb.OnStockTop != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnStockTop(&d)
			}
		}
	case MsgTypeOptionTop:
		if cb.OnOptionTop != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnOptionTop(&d)
			}
		}
	case MsgTypeFullTick:
		if cb.OnFullTick != nil {
			var d TickData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnFullTick(&d)
			}
		}
	case MsgTypeQuoteBBO:
		if cb.OnQuoteBBO != nil {
			var d QuoteData
			if json.Unmarshal(msg.Data, &d) == nil {
				cb.OnQuoteBBO(&d)
			}
		}
	}
}

// subscribe 内部订阅方法
func (c *PushClient) subscribe(subject SubjectType, symbols []string, account string, market string) error {
	req := SubscribeRequest{
		Subject: subject,
		Symbols: symbols,
		Account: account,
		Market:  market,
	}
	msg, err := NewPushMessage(MsgTypeSubscribe, subject, &req)
	if err != nil {
		return err
	}
	return c.sendMessage(msg)
}

// unsubscribe 内部退订方法
func (c *PushClient) unsubscribe(subject SubjectType, symbols []string) error {
	req := SubscribeRequest{
		Subject: subject,
		Symbols: symbols,
	}
	msg, err := NewPushMessage(MsgTypeUnsubscribe, subject, &req)
	if err != nil {
		return err
	}
	return c.sendMessage(msg)
}

// SubscribeQuote 订阅行情
func (c *PushClient) SubscribeQuote(symbols []string) error {
	if err := c.subscribe(SubjectQuote, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectQuote, symbols)
	return nil
}

// UnsubscribeQuote 退订行情
func (c *PushClient) UnsubscribeQuote(symbols []string) error {
	if err := c.unsubscribe(SubjectQuote, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectQuote, symbols)
	return nil
}

// SubscribeTick 订阅逐笔成交
func (c *PushClient) SubscribeTick(symbols []string) error {
	if err := c.subscribe(SubjectTick, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectTick, symbols)
	return nil
}

// UnsubscribeTick 退订逐笔成交
func (c *PushClient) UnsubscribeTick(symbols []string) error {
	if err := c.unsubscribe(SubjectTick, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectTick, symbols)
	return nil
}

// SubscribeDepth 订阅深度行情
func (c *PushClient) SubscribeDepth(symbols []string) error {
	if err := c.subscribe(SubjectDepth, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectDepth, symbols)
	return nil
}

// UnsubscribeDepth 退订深度行情
func (c *PushClient) UnsubscribeDepth(symbols []string) error {
	if err := c.unsubscribe(SubjectDepth, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectDepth, symbols)
	return nil
}

// SubscribeOption 订阅期权行情
func (c *PushClient) SubscribeOption(symbols []string) error {
	if err := c.subscribe(SubjectOption, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectOption, symbols)
	return nil
}

// UnsubscribeOption 退订期权行情
func (c *PushClient) UnsubscribeOption(symbols []string) error {
	if err := c.unsubscribe(SubjectOption, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectOption, symbols)
	return nil
}

// SubscribeFuture 订阅期货行情
func (c *PushClient) SubscribeFuture(symbols []string) error {
	if err := c.subscribe(SubjectFuture, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectFuture, symbols)
	return nil
}

// UnsubscribeFuture 退订期货行情
func (c *PushClient) UnsubscribeFuture(symbols []string) error {
	if err := c.unsubscribe(SubjectFuture, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectFuture, symbols)
	return nil
}

// SubscribeKline 订阅 K 线
func (c *PushClient) SubscribeKline(symbols []string) error {
	if err := c.subscribe(SubjectKline, symbols, "", ""); err != nil {
		return err
	}
	c.addSubscription(SubjectKline, symbols)
	return nil
}

// UnsubscribeKline 退订 K 线
func (c *PushClient) UnsubscribeKline(symbols []string) error {
	if err := c.unsubscribe(SubjectKline, symbols); err != nil {
		return err
	}
	c.removeSubscription(SubjectKline, symbols)
	return nil
}

// SubscribeAsset 订阅资产变动
func (c *PushClient) SubscribeAsset(account string) error {
	if account == "" {
		account = c.config.Account
	}
	if err := c.subscribe(SubjectAsset, nil, account, ""); err != nil {
		return err
	}
	c.mu.Lock()
	c.accountSubs[SubjectAsset] = true
	c.mu.Unlock()
	return nil
}

// UnsubscribeAsset 退订资产变动
func (c *PushClient) UnsubscribeAsset() error {
	if err := c.unsubscribe(SubjectAsset, nil); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.accountSubs, SubjectAsset)
	c.mu.Unlock()
	return nil
}

// SubscribePosition 订阅持仓变动
func (c *PushClient) SubscribePosition(account string) error {
	if account == "" {
		account = c.config.Account
	}
	if err := c.subscribe(SubjectPosition, nil, account, ""); err != nil {
		return err
	}
	c.mu.Lock()
	c.accountSubs[SubjectPosition] = true
	c.mu.Unlock()
	return nil
}

// UnsubscribePosition 退订持仓变动
func (c *PushClient) UnsubscribePosition() error {
	if err := c.unsubscribe(SubjectPosition, nil); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.accountSubs, SubjectPosition)
	c.mu.Unlock()
	return nil
}

// SubscribeOrder 订阅订单状态
func (c *PushClient) SubscribeOrder(account string) error {
	if account == "" {
		account = c.config.Account
	}
	if err := c.subscribe(SubjectOrder, nil, account, ""); err != nil {
		return err
	}
	c.mu.Lock()
	c.accountSubs[SubjectOrder] = true
	c.mu.Unlock()
	return nil
}

// UnsubscribeOrder 退订订单状态
func (c *PushClient) UnsubscribeOrder() error {
	if err := c.unsubscribe(SubjectOrder, nil); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.accountSubs, SubjectOrder)
	c.mu.Unlock()
	return nil
}

// SubscribeTransaction 订阅成交明细
func (c *PushClient) SubscribeTransaction(account string) error {
	if account == "" {
		account = c.config.Account
	}
	if err := c.subscribe(SubjectTransaction, nil, account, ""); err != nil {
		return err
	}
	c.mu.Lock()
	c.accountSubs[SubjectTransaction] = true
	c.mu.Unlock()
	return nil
}

// UnsubscribeTransaction 退订成交明细
func (c *PushClient) UnsubscribeTransaction() error {
	if err := c.unsubscribe(SubjectTransaction, nil); err != nil {
		return err
	}
	c.mu.Lock()
	delete(c.accountSubs, SubjectTransaction)
	c.mu.Unlock()
	return nil
}

// GetSubscriptions 获取当前订阅状态
func (c *PushClient) GetSubscriptions() map[SubjectType][]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[SubjectType][]string)
	for subject, symbols := range c.subscriptions {
		list := make([]string, 0, len(symbols))
		for s := range symbols {
			list = append(list, s)
		}
		result[subject] = list
	}
	return result
}

// GetAccountSubscriptions 获取账户级别订阅状态
func (c *PushClient) GetAccountSubscriptions() []SubjectType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]SubjectType, 0, len(c.accountSubs))
	for subject := range c.accountSubs {
		result = append(result, subject)
	}
	return result
}

// addSubscription 添加订阅记录
func (c *PushClient) addSubscription(subject SubjectType, symbols []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.subscriptions[subject] == nil {
		c.subscriptions[subject] = make(map[string]bool)
	}
	for _, s := range symbols {
		c.subscriptions[subject][s] = true
	}
}

// removeSubscription 移除订阅记录
func (c *PushClient) removeSubscription(subject SubjectType, symbols []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if symbols == nil {
		// 退订全部
		delete(c.subscriptions, subject)
		return
	}
	if m, ok := c.subscriptions[subject]; ok {
		for _, s := range symbols {
			delete(m, s)
		}
		if len(m) == 0 {
			delete(c.subscriptions, subject)
		}
	}
}
