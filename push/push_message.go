// Package push 提供 WebSocket 推送客户端，支持实时行情和账户推送。
package push

import "encoding/json"

// MessageType 推送消息类型
type MessageType string

const (
	// 连接相关
	MsgTypeConnect    MessageType = "connect"
	MsgTypeConnected  MessageType = "connected"
	MsgTypeDisconnect MessageType = "disconnect"
	MsgTypeHeartbeat  MessageType = "heartbeat"
	MsgTypeKickout    MessageType = "kickout"
	MsgTypeError      MessageType = "error"

	// 订阅相关
	MsgTypeSubscribe   MessageType = "subscribe"
	MsgTypeUnsubscribe MessageType = "unsubscribe"

	// 行情推送
	MsgTypeQuote      MessageType = "quote"
	MsgTypeTick       MessageType = "tick"
	MsgTypeDepth      MessageType = "depth"
	MsgTypeOption     MessageType = "option"
	MsgTypeFuture     MessageType = "future"
	MsgTypeKline      MessageType = "kline"
	MsgTypeStockTop   MessageType = "stock_top"
	MsgTypeOptionTop  MessageType = "option_top"
	MsgTypeFullTick   MessageType = "full_tick"
	MsgTypeQuoteBBO   MessageType = "quote_bbo"

	// 账户推送
	MsgTypeAsset       MessageType = "asset"
	MsgTypePosition    MessageType = "position"
	MsgTypeOrder       MessageType = "order"
	MsgTypeTransaction MessageType = "transaction"
)

// SubjectType 订阅主题类型
type SubjectType string

const (
	SubjectQuote       SubjectType = "quote"
	SubjectTick        SubjectType = "tick"
	SubjectDepth       SubjectType = "depth"
	SubjectOption      SubjectType = "option"
	SubjectFuture      SubjectType = "future"
	SubjectKline       SubjectType = "kline"
	SubjectStockTop    SubjectType = "stock_top"
	SubjectOptionTop   SubjectType = "option_top"
	SubjectFullTick    SubjectType = "full_tick"
	SubjectQuoteBBO    SubjectType = "quote_bbo"
	SubjectAsset       SubjectType = "asset"
	SubjectPosition    SubjectType = "position"
	SubjectOrder       SubjectType = "order"
	SubjectTransaction SubjectType = "transaction"
)

// PushMessage 推送消息的通用结构（简化的 JSON 格式替代 Protobuf）
type PushMessage struct {
	Type    MessageType     `json:"type"`
	Subject SubjectType     `json:"subject,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// ConnectRequest 连接认证请求
type ConnectRequest struct {
	TigerID   string `json:"tigerId"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// SubscribeRequest 订阅/退订请求
type SubscribeRequest struct {
	Subject SubjectType `json:"subject"`
	Symbols []string    `json:"symbols,omitempty"`
	Account string      `json:"account,omitempty"`
	Market  string      `json:"market,omitempty"`
}

// QuoteData 行情推送数据
type QuoteData struct {
	Symbol     string  `json:"symbol"`
	LatestPrice float64 `json:"latestPrice,omitempty"`
	PreClose   float64 `json:"preClose,omitempty"`
	Open       float64 `json:"open,omitempty"`
	High       float64 `json:"high,omitempty"`
	Low        float64 `json:"low,omitempty"`
	Volume     int64   `json:"volume,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	Timestamp  int64   `json:"timestamp,omitempty"`
}

// TickData 逐笔成交推送数据
type TickData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Volume    int64   `json:"volume"`
	Type      string  `json:"type,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// DepthData 深度行情推送数据
type DepthData struct {
	Symbol string       `json:"symbol"`
	Asks   []PriceLevel `json:"asks,omitempty"`
	Bids   []PriceLevel `json:"bids,omitempty"`
}

// PriceLevel 价格档位
type PriceLevel struct {
	Price  float64 `json:"price"`
	Volume int64   `json:"volume"`
	Count  int     `json:"count,omitempty"`
}

// KlineData K 线推送数据
type KlineData struct {
	Symbol    string  `json:"symbol"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// AssetData 资产推送数据
type AssetData struct {
	Account        string  `json:"account"`
	NetLiquidation float64 `json:"netLiquidation,omitempty"`
	EquityWithLoan float64 `json:"equityWithLoan,omitempty"`
	CashBalance    float64 `json:"cashBalance,omitempty"`
	BuyingPower    float64 `json:"buyingPower,omitempty"`
	Currency       string  `json:"currency,omitempty"`
}

// PositionData 持仓推送数据
type PositionData struct {
	Account       string  `json:"account"`
	Symbol        string  `json:"symbol"`
	SecType       string  `json:"secType,omitempty"`
	Quantity      int     `json:"quantity"`
	AverageCost   float64 `json:"averageCost,omitempty"`
	MarketPrice   float64 `json:"marketPrice,omitempty"`
	MarketValue   float64 `json:"marketValue,omitempty"`
	UnrealizedPnl float64 `json:"unrealizedPnl,omitempty"`
}

// OrderData 订单推送数据
type OrderData struct {
	Account      string  `json:"account"`
	ID           int64   `json:"id,omitempty"`
	OrderID      int64   `json:"orderId,omitempty"`
	Symbol       string  `json:"symbol"`
	Action       string  `json:"action,omitempty"`
	OrderType    string  `json:"orderType,omitempty"`
	Quantity     int     `json:"quantity,omitempty"`
	LimitPrice   float64 `json:"limitPrice,omitempty"`
	Status       string  `json:"status,omitempty"`
	Filled       int     `json:"filled,omitempty"`
	AvgFillPrice float64 `json:"avgFillPrice,omitempty"`
}

// TransactionData 成交推送数据
type TransactionData struct {
	Account   string  `json:"account"`
	ID        int64   `json:"id,omitempty"`
	OrderID   int64   `json:"orderId,omitempty"`
	Symbol    string  `json:"symbol"`
	Action    string  `json:"action,omitempty"`
	Price     float64 `json:"price,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

// Serialize 将推送消息序列化为 JSON 字节
func (m *PushMessage) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// DeserializeMessage 从 JSON 字节反序列化推送消息
func DeserializeMessage(data []byte) (*PushMessage, error) {
	var msg PushMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// NewPushMessage 创建一个推送消息
func NewPushMessage(msgType MessageType, subject SubjectType, data interface{}) (*PushMessage, error) {
	msg := &PushMessage{
		Type:    msgType,
		Subject: subject,
	}
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		msg.Data = raw
	}
	return msg, nil
}
