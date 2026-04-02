package push

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tigerfintech/openapi-go-sdk/config"
)

// ===== 测试辅助工具 =====

// generateTestPrivateKey 生成测试用 RSA PKCS#1 私钥 PEM
func generateTestPrivateKey(t *testing.T) string {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})
	return string(privPEM)
}

// newTestConfig 创建测试用 ClientConfig（跳过校验）
func newTestConfig(t *testing.T) *config.ClientConfig {
	t.Helper()
	return &config.ClientConfig{
		TigerID:    "test_tiger_id",
		PrivateKey: generateTestPrivateKey(t),
		Account:    "test_account",
		Language:   "zh_CN",
		Timeout:    15 * time.Second,
		ServerURL:  "https://openapi.tigerfintech.com/gateway",
	}
}

// mockWSServer 创建一个模拟 WebSocket 服务器
// handler 处理每个 WebSocket 连接
func mockWSServer(t *testing.T, handler func(conn *websocket.Conn)) *httptest.Server {
	t.Helper()
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("WebSocket 升级失败: %v", err)
			return
		}
		defer conn.Close()
		handler(conn)
	}))
	return server
}

// wsURL 将 http:// 转换为 ws://
func wsURL(server *httptest.Server) string {
	return "ws" + strings.TrimPrefix(server.URL, "http")
}

// mockDialer 测试用 WebSocket 拨号器
type mockDialer struct {
	url string
}

func (d *mockDialer) Dial(urlStr string, timeout time.Duration) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: timeout,
	}
	conn, _, err := dialer.Dial(d.url, nil)
	return conn, err
}

// ===== 20.1 连接和认证测试 =====

// TestPushClient_NewPushClient 测试创建 PushClient
func TestPushClient_NewPushClient(t *testing.T) {
	cfg := newTestConfig(t)
	client := NewPushClient(cfg)

	if client == nil {
		t.Fatal("NewPushClient 不应返回 nil")
	}
	if client.State() != StateDisconnected {
		t.Errorf("初始状态应为 StateDisconnected，实际为 %v", client.State())
	}
	if client.pushURL != defaultPushURL {
		t.Errorf("默认推送地址应为 %s，实际为 %s", defaultPushURL, client.pushURL)
	}
	if !client.autoReconnect {
		t.Error("默认应启用自动重连")
	}
}

// TestPushClient_Options 测试 PushClient 配置选项
func TestPushClient_Options(t *testing.T) {
	cfg := newTestConfig(t)
	client := NewPushClient(cfg,
		WithPushURL("wss://custom.example.com"),
		WithHeartbeatInterval(20*time.Second),
		WithReconnectInterval(10*time.Second),
		WithAutoReconnect(false),
		WithConnectTimeout(60*time.Second),
	)

	if client.pushURL != "wss://custom.example.com" {
		t.Errorf("自定义推送地址未生效")
	}
	if client.heartbeatInterval != 20*time.Second {
		t.Errorf("自定义心跳间隔未生效")
	}
	if client.reconnectInterval != 10*time.Second {
		t.Errorf("自定义重连间隔未生效")
	}
	if client.autoReconnect {
		t.Error("禁用自动重连未生效")
	}
	if client.connectTimeout != 60*time.Second {
		t.Errorf("自定义连接超时未生效")
	}
}

// TestPushClient_ConnectAndAuthenticate 测试连接和认证流程
func TestPushClient_ConnectAndAuthenticate(t *testing.T) {
	var receivedMsg PushMessage
	var msgMu sync.Mutex
	authReceived := make(chan struct{})

	server := mockWSServer(t, func(conn *websocket.Conn) {
		// 读取认证消息
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Logf("读取消息失败: %v", err)
			return
		}
		msgMu.Lock()
		json.Unmarshal(data, &receivedMsg)
		msgMu.Unlock()
		close(authReceived)

		// 保持连接直到客户端断开
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	err := client.Connect()
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer client.Disconnect()

	// 等待服务端收到认证消息
	select {
	case <-authReceived:
	case <-time.After(3 * time.Second):
		t.Fatal("等待认证消息超时")
	}

	msgMu.Lock()
	defer msgMu.Unlock()

	// 验证认证消息
	if receivedMsg.Type != MsgTypeConnect {
		t.Errorf("认证消息类型应为 %s，实际为 %s", MsgTypeConnect, receivedMsg.Type)
	}

	// 解析认证数据
	var req ConnectRequest
	json.Unmarshal(receivedMsg.Data, &req)
	if req.TigerID != "test_tiger_id" {
		t.Errorf("认证消息中 TigerID 应为 test_tiger_id，实际为 %s", req.TigerID)
	}
	if req.Sign == "" {
		t.Error("认证消息中签名不应为空")
	}
	if req.Version != "2.0" {
		t.Errorf("认证消息中版本应为 2.0，实际为 %s", req.Version)
	}

	// 验证连接状态
	if client.State() != StateConnected {
		t.Errorf("连接后状态应为 StateConnected，实际为 %v", client.State())
	}
}

// TestPushClient_Disconnect 测试断开连接
func TestPushClient_Disconnect(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	client.Connect()
	err := client.Disconnect()
	if err != nil {
		t.Fatalf("断开连接失败: %v", err)
	}

	if client.State() != StateDisconnected {
		t.Errorf("断开后状态应为 StateDisconnected，实际为 %v", client.State())
	}
}

// TestPushClient_DisconnectWhenNotConnected 测试未连接时断开不报错
func TestPushClient_DisconnectWhenNotConnected(t *testing.T) {
	cfg := newTestConfig(t)
	client := NewPushClient(cfg)

	err := client.Disconnect()
	if err != nil {
		t.Fatalf("未连接时断开不应报错: %v", err)
	}
}

// TestPushClient_ConnectWhenAlreadyConnected 测试重复连接应报错
func TestPushClient_ConnectWhenAlreadyConnected(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	client.Connect()
	defer client.Disconnect()

	err := client.Connect()
	if err == nil {
		t.Fatal("重复连接应返回错误")
	}
}

// ===== 20.3 Protobuf 序列化测试（简化为 JSON） =====

// TestPushMessage_Serialize 测试消息序列化
func TestPushMessage_Serialize(t *testing.T) {
	quoteData := &QuoteData{
		Symbol:      "AAPL",
		LatestPrice: 150.25,
		Volume:      1000000,
		Timestamp:   1700000000,
	}

	msg, err := NewPushMessage(MsgTypeQuote, SubjectQuote, quoteData)
	if err != nil {
		t.Fatalf("创建推送消息失败: %v", err)
	}

	data, err := msg.Serialize()
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 反序列化验证
	msg2, err := DeserializeMessage(data)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if msg2.Type != MsgTypeQuote {
		t.Errorf("消息类型应为 %s，实际为 %s", MsgTypeQuote, msg2.Type)
	}
	if msg2.Subject != SubjectQuote {
		t.Errorf("消息主题应为 %s，实际为 %s", SubjectQuote, msg2.Subject)
	}

	var restored QuoteData
	if err := json.Unmarshal(msg2.Data, &restored); err != nil {
		t.Fatalf("解析行情数据失败: %v", err)
	}
	if restored.Symbol != "AAPL" {
		t.Errorf("Symbol 应为 AAPL，实际为 %s", restored.Symbol)
	}
	if restored.LatestPrice != 150.25 {
		t.Errorf("LatestPrice 应为 150.25，实际为 %f", restored.LatestPrice)
	}
}

// TestPushMessage_SerializeRoundTrip_Tick 测试逐笔成交消息 round-trip
func TestPushMessage_SerializeRoundTrip_Tick(t *testing.T) {
	original := &TickData{
		Symbol:    "TSLA",
		Price:     250.50,
		Volume:    500,
		Type:      "BUY",
		Timestamp: 1700000001,
	}

	msg, _ := NewPushMessage(MsgTypeTick, SubjectTick, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored TickData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Symbol != original.Symbol || restored.Price != original.Price ||
		restored.Volume != original.Volume || restored.Type != original.Type {
		t.Errorf("TickData round-trip 失败: 原始=%+v, 恢复=%+v", original, restored)
	}
}

// TestPushMessage_SerializeRoundTrip_Depth 测试深度行情消息 round-trip
func TestPushMessage_SerializeRoundTrip_Depth(t *testing.T) {
	original := &DepthData{
		Symbol: "AAPL",
		Asks: []PriceLevel{
			{Price: 150.10, Volume: 100, Count: 5},
			{Price: 150.20, Volume: 200, Count: 3},
		},
		Bids: []PriceLevel{
			{Price: 150.00, Volume: 150, Count: 4},
			{Price: 149.90, Volume: 300, Count: 6},
		},
	}

	msg, _ := NewPushMessage(MsgTypeDepth, SubjectDepth, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored DepthData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Symbol != "AAPL" {
		t.Errorf("Symbol 应为 AAPL，实际为 %s", restored.Symbol)
	}
	if len(restored.Asks) != 2 || len(restored.Bids) != 2 {
		t.Errorf("Asks/Bids 长度不匹配")
	}
	if restored.Asks[0].Price != 150.10 || restored.Bids[0].Price != 150.00 {
		t.Errorf("价格档位数据不匹配")
	}
}

// TestPushMessage_SerializeRoundTrip_Kline 测试 K 线消息 round-trip
func TestPushMessage_SerializeRoundTrip_Kline(t *testing.T) {
	original := &KlineData{
		Symbol: "GOOG", Open: 140.0, High: 145.0, Low: 139.0,
		Close: 143.5, Volume: 2000000, Timestamp: 1700000002,
	}

	msg, _ := NewPushMessage(MsgTypeKline, SubjectKline, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored KlineData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Symbol != original.Symbol || restored.Open != original.Open ||
		restored.Close != original.Close || restored.Volume != original.Volume {
		t.Errorf("KlineData round-trip 失败")
	}
}

// TestPushMessage_SerializeRoundTrip_Order 测试订单消息 round-trip
func TestPushMessage_SerializeRoundTrip_Order(t *testing.T) {
	original := &OrderData{
		Account: "acc123", ID: 1001, OrderID: 2001, Symbol: "AAPL",
		Action: "BUY", OrderType: "LMT", Quantity: 100,
		LimitPrice: 150.0, Status: "Filled", Filled: 100, AvgFillPrice: 149.95,
	}

	msg, _ := NewPushMessage(MsgTypeOrder, SubjectOrder, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored OrderData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Account != original.Account || restored.ID != original.ID ||
		restored.Symbol != original.Symbol || restored.Status != original.Status {
		t.Errorf("OrderData round-trip 失败")
	}
}

// TestPushMessage_SerializeRoundTrip_Asset 测试资产消息 round-trip
func TestPushMessage_SerializeRoundTrip_Asset(t *testing.T) {
	original := &AssetData{
		Account: "acc123", NetLiquidation: 100000.50,
		CashBalance: 50000.25, BuyingPower: 200000.0, Currency: "USD",
	}

	msg, _ := NewPushMessage(MsgTypeAsset, SubjectAsset, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored AssetData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Account != original.Account || restored.NetLiquidation != original.NetLiquidation ||
		restored.CashBalance != original.CashBalance {
		t.Errorf("AssetData round-trip 失败")
	}
}

// TestPushMessage_SerializeRoundTrip_Position 测试持仓消息 round-trip
func TestPushMessage_SerializeRoundTrip_Position(t *testing.T) {
	original := &PositionData{
		Account: "acc123", Symbol: "AAPL", SecType: "STK",
		Quantity: 100, AverageCost: 145.50, MarketPrice: 150.25,
		MarketValue: 15025.0, UnrealizedPnl: 475.0,
	}

	msg, _ := NewPushMessage(MsgTypePosition, SubjectPosition, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored PositionData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Symbol != original.Symbol || restored.Quantity != original.Quantity ||
		restored.AverageCost != original.AverageCost {
		t.Errorf("PositionData round-trip 失败")
	}
}

// TestPushMessage_SerializeRoundTrip_Transaction 测试成交消息 round-trip
func TestPushMessage_SerializeRoundTrip_Transaction(t *testing.T) {
	original := &TransactionData{
		Account: "acc123", ID: 3001, OrderID: 2001, Symbol: "AAPL",
		Action: "BUY", Price: 149.95, Quantity: 100, Timestamp: 1700000003,
	}

	msg, _ := NewPushMessage(MsgTypeTransaction, SubjectTransaction, original)
	data, _ := msg.Serialize()
	msg2, _ := DeserializeMessage(data)

	var restored TransactionData
	json.Unmarshal(msg2.Data, &restored)

	if restored.Symbol != original.Symbol || restored.Price != original.Price ||
		restored.Quantity != original.Quantity {
		t.Errorf("TransactionData round-trip 失败")
	}
}

// TestDeserializeMessage_Invalid 测试反序列化无效数据
func TestDeserializeMessage_Invalid(t *testing.T) {
	_, err := DeserializeMessage([]byte("not json"))
	if err == nil {
		t.Fatal("反序列化无效数据应返回错误")
	}
}

// TestNewPushMessage_NilData 测试创建无数据的消息
func TestNewPushMessage_NilData(t *testing.T) {
	msg, err := NewPushMessage(MsgTypeHeartbeat, "", nil)
	if err != nil {
		t.Fatalf("创建心跳消息失败: %v", err)
	}
	if msg.Data != nil {
		t.Error("心跳消息不应有数据")
	}
}

// ===== 20.5 订阅/退订测试 =====

// TestPushClient_SubscribeQuote 测试订阅行情
func TestPushClient_SubscribeQuote(t *testing.T) {
	var receivedMsgs []PushMessage
	var mu sync.Mutex
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var msg PushMessage
			json.Unmarshal(data, &msg)
			mu.Lock()
			receivedMsgs = append(receivedMsgs, msg)
			mu.Unlock()
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	// 订阅行情
	err := client.SubscribeQuote([]string{"AAPL", "TSLA"})
	if err != nil {
		t.Fatalf("订阅行情失败: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// 验证订阅状态
	subs := client.GetSubscriptions()
	if symbols, ok := subs[SubjectQuote]; !ok {
		t.Error("应有 quote 订阅记录")
	} else if len(symbols) != 2 {
		t.Errorf("应订阅 2 个标的，实际 %d 个", len(symbols))
	}
}

// TestPushClient_UnsubscribeQuote 测试退订行情
func TestPushClient_UnsubscribeQuote(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	client.SubscribeQuote([]string{"AAPL", "TSLA", "GOOG"})
	client.UnsubscribeQuote([]string{"TSLA"})

	time.Sleep(100 * time.Millisecond)

	subs := client.GetSubscriptions()
	if symbols, ok := subs[SubjectQuote]; !ok {
		t.Error("应有 quote 订阅记录")
	} else if len(symbols) != 2 {
		t.Errorf("退订后应剩 2 个标的，实际 %d 个", len(symbols))
	}
}

// TestPushClient_SubscribeMultipleSubjects 测试订阅多种行情
func TestPushClient_SubscribeMultipleSubjects(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	client.SubscribeQuote([]string{"AAPL"})
	client.SubscribeTick([]string{"TSLA"})
	client.SubscribeDepth([]string{"GOOG"})
	client.SubscribeOption([]string{"AAPL"})
	client.SubscribeFuture([]string{"ES"})
	client.SubscribeKline([]string{"AAPL"})

	subs := client.GetSubscriptions()
	if len(subs) != 6 {
		t.Errorf("应有 6 种订阅，实际 %d 种", len(subs))
	}
}

// TestPushClient_UnsubscribeAll 测试退订全部（传 nil）
func TestPushClient_UnsubscribeAll(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	client.SubscribeQuote([]string{"AAPL", "TSLA"})
	client.UnsubscribeQuote(nil) // 退订全部

	subs := client.GetSubscriptions()
	if _, ok := subs[SubjectQuote]; ok {
		t.Error("退订全部后不应有 quote 订阅记录")
	}
}

// ===== 20.7 账户推送测试 =====

// TestPushClient_SubscribeAsset 测试订阅资产变动
func TestPushClient_SubscribeAsset(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	err := client.SubscribeAsset("")
	if err != nil {
		t.Fatalf("订阅资产失败: %v", err)
	}

	acctSubs := client.GetAccountSubscriptions()
	found := false
	for _, s := range acctSubs {
		if s == SubjectAsset {
			found = true
		}
	}
	if !found {
		t.Error("应有 asset 账户订阅记录")
	}
}

// TestPushClient_SubscribeAllAccountPush 测试订阅所有账户推送
func TestPushClient_SubscribeAllAccountPush(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	client.SubscribeAsset("")
	client.SubscribePosition("")
	client.SubscribeOrder("")
	client.SubscribeTransaction("")

	acctSubs := client.GetAccountSubscriptions()
	if len(acctSubs) != 4 {
		t.Errorf("应有 4 种账户订阅，实际 %d 种", len(acctSubs))
	}
}

// TestPushClient_UnsubscribeAccountPush 测试退订账户推送
func TestPushClient_UnsubscribeAccountPush(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}
	client.Connect()
	defer client.Disconnect()

	client.SubscribeAsset("")
	client.SubscribePosition("")
	client.UnsubscribeAsset()
	client.UnsubscribePosition()

	acctSubs := client.GetAccountSubscriptions()
	if len(acctSubs) != 0 {
		t.Errorf("退订后不应有账户订阅，实际 %d 种", len(acctSubs))
	}
}

// ===== 20.9 回调和重连测试 =====

// TestPushClient_ConnectCallback 测试连接成功回调
func TestPushClient_ConnectCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var connected int32
	client.SetCallbacks(Callbacks{
		OnConnect: func() {
			atomic.StoreInt32(&connected, 1)
		},
	})

	client.Connect()
	defer client.Disconnect()

	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&connected) != 1 {
		t.Error("连接成功回调未触发")
	}
}

// TestPushClient_DisconnectCallback 测试断开连接回调
func TestPushClient_DisconnectCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var disconnected int32
	client.SetCallbacks(Callbacks{
		OnDisconnect: func() {
			atomic.StoreInt32(&disconnected, 1)
		},
	})

	client.Connect()
	client.Disconnect()

	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&disconnected) != 1 {
		t.Error("断开连接回调未触发")
	}
}

// TestPushClient_QuoteCallback 测试行情推送回调
func TestPushClient_QuoteCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		// 读取认证消息
		conn.ReadMessage()

		// 发送行情推送
		quoteData := &QuoteData{Symbol: "AAPL", LatestPrice: 155.0}
		msg, _ := NewPushMessage(MsgTypeQuote, SubjectQuote, quoteData)
		data, _ := msg.Serialize()
		conn.WriteMessage(websocket.TextMessage, data)

		// 保持连接
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var receivedQuote *QuoteData
	var quoteMu sync.Mutex
	client.SetCallbacks(Callbacks{
		OnQuote: func(data *QuoteData) {
			quoteMu.Lock()
			receivedQuote = data
			quoteMu.Unlock()
		},
	})

	client.Connect()
	defer client.Disconnect()

	// 等待消息到达
	time.Sleep(300 * time.Millisecond)

	quoteMu.Lock()
	defer quoteMu.Unlock()
	if receivedQuote == nil {
		t.Fatal("行情回调未触发")
	}
	if receivedQuote.Symbol != "AAPL" {
		t.Errorf("Symbol 应为 AAPL，实际为 %s", receivedQuote.Symbol)
	}
	if receivedQuote.LatestPrice != 155.0 {
		t.Errorf("LatestPrice 应为 155.0，实际为 %f", receivedQuote.LatestPrice)
	}
}

// TestPushClient_OrderCallback 测试订单推送回调
func TestPushClient_OrderCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		conn.ReadMessage()
		orderData := &OrderData{
			Account: "acc123", Symbol: "AAPL", Status: "Filled",
		}
		msg, _ := NewPushMessage(MsgTypeOrder, SubjectOrder, orderData)
		data, _ := msg.Serialize()
		conn.WriteMessage(websocket.TextMessage, data)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var received *OrderData
	var mu sync.Mutex
	client.SetCallbacks(Callbacks{
		OnOrder: func(data *OrderData) {
			mu.Lock()
			received = data
			mu.Unlock()
		},
	})

	client.Connect()
	defer client.Disconnect()
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if received == nil {
		t.Fatal("订单回调未触发")
	}
	if received.Status != "Filled" {
		t.Errorf("Status 应为 Filled，实际为 %s", received.Status)
	}
}

// TestPushClient_KickoutCallback 测试被踢出回调
func TestPushClient_KickoutCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		conn.ReadMessage()
		msg, _ := NewPushMessage(MsgTypeKickout, "", "另一设备登录")
		data, _ := msg.Serialize()
		conn.WriteMessage(websocket.TextMessage, data)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var kickoutMsg string
	var mu sync.Mutex
	client.SetCallbacks(Callbacks{
		OnKickout: func(message string) {
			mu.Lock()
			kickoutMsg = message
			mu.Unlock()
		},
	})

	client.Connect()
	defer client.Disconnect()
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if kickoutMsg != "另一设备登录" {
		t.Errorf("踢出消息应为 '另一设备登录'，实际为 '%s'", kickoutMsg)
	}
}

// TestPushClient_ErrorCallback 测试错误回调
func TestPushClient_ErrorCallback(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		conn.ReadMessage()
		msg, _ := NewPushMessage(MsgTypeError, "", "服务端内部错误")
		data, _ := msg.Serialize()
		conn.WriteMessage(websocket.TextMessage, data)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var errReceived error
	var mu sync.Mutex
	client.SetCallbacks(Callbacks{
		OnError: func(err error) {
			mu.Lock()
			errReceived = err
			mu.Unlock()
		},
	})

	client.Connect()
	defer client.Disconnect()
	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if errReceived == nil {
		t.Fatal("错误回调未触发")
	}
}

// TestPushClient_MultipleCallbacks 测试多种回调同时注册
func TestPushClient_MultipleCallbacks(t *testing.T) {
	server := mockWSServer(t, func(conn *websocket.Conn) {
		conn.ReadMessage()
		// 发送多种消息
		q, _ := NewPushMessage(MsgTypeQuote, SubjectQuote, &QuoteData{Symbol: "AAPL"})
		d1, _ := q.Serialize()
		conn.WriteMessage(websocket.TextMessage, d1)

		tk, _ := NewPushMessage(MsgTypeTick, SubjectTick, &TickData{Symbol: "TSLA", Price: 250.0})
		d2, _ := tk.Serialize()
		conn.WriteMessage(websocket.TextMessage, d2)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	})
	defer server.Close()

	cfg := newTestConfig(t)
	client := NewPushClient(cfg, WithAutoReconnect(false))
	client.dialer = &mockDialer{url: wsURL(server)}

	var quoteCount, tickCount int32
	client.SetCallbacks(Callbacks{
		OnQuote: func(data *QuoteData) { atomic.AddInt32(&quoteCount, 1) },
		OnTick:  func(data *TickData) { atomic.AddInt32(&tickCount, 1) },
	})

	client.Connect()
	defer client.Disconnect()
	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&quoteCount) != 1 {
		t.Errorf("行情回调应触发 1 次，实际 %d 次", quoteCount)
	}
	if atomic.LoadInt32(&tickCount) != 1 {
		t.Errorf("逐笔回调应触发 1 次，实际 %d 次", tickCount)
	}
}

// TestPushClient_SubscriptionStateManagement 测试订阅状态管理的纯逻辑
func TestPushClient_SubscriptionStateManagement(t *testing.T) {
	cfg := newTestConfig(t)
	client := NewPushClient(cfg)

	// 直接测试内部订阅状态管理（不需要连接）
	client.addSubscription(SubjectQuote, []string{"AAPL", "TSLA"})
	client.addSubscription(SubjectTick, []string{"GOOG"})

	subs := client.GetSubscriptions()
	if len(subs) != 2 {
		t.Errorf("应有 2 种订阅，实际 %d 种", len(subs))
	}

	// 追加订阅
	client.addSubscription(SubjectQuote, []string{"GOOG"})
	subs = client.GetSubscriptions()
	if len(subs[SubjectQuote]) != 3 {
		t.Errorf("quote 应有 3 个标的，实际 %d 个", len(subs[SubjectQuote]))
	}

	// 部分退订
	client.removeSubscription(SubjectQuote, []string{"TSLA"})
	subs = client.GetSubscriptions()
	if len(subs[SubjectQuote]) != 2 {
		t.Errorf("退订后 quote 应有 2 个标的，实际 %d 个", len(subs[SubjectQuote]))
	}

	// 全部退订
	client.removeSubscription(SubjectQuote, nil)
	subs = client.GetSubscriptions()
	if _, ok := subs[SubjectQuote]; ok {
		t.Error("全部退订后不应有 quote 记录")
	}
}
