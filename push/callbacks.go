package push

// QuoteCallback 行情推送回调
type QuoteCallback func(data *QuoteData)

// TickCallback 逐笔成交推送回调
type TickCallback func(data *TickData)

// DepthCallback 深度行情推送回调
type DepthCallback func(data *DepthData)

// KlineCallback K 线推送回调
type KlineCallback func(data *KlineData)

// AssetCallback 资产变动推送回调
type AssetCallback func(data *AssetData)

// PositionCallback 持仓变动推送回调
type PositionCallback func(data *PositionData)

// OrderCallback 订单状态推送回调
type OrderCallback func(data *OrderData)

// TransactionCallback 成交明细推送回调
type TransactionCallback func(data *TransactionData)

// ConnectCallback 连接成功回调
type ConnectCallback func()

// DisconnectCallback 断开连接回调
type DisconnectCallback func()

// ErrorCallback 错误回调
type ErrorCallback func(err error)

// KickoutCallback 被踢出回调（同一 tiger_id 在其他地方登录）
type KickoutCallback func(message string)

// SubscribeCallback 订阅成功回调
type SubscribeCallback func(subject SubjectType, symbols []string)

// UnsubscribeCallback 退订成功回调
type UnsubscribeCallback func(subject SubjectType, symbols []string)

// Callbacks 所有回调函数的集合
type Callbacks struct {
	// 行情推送回调
	OnQuote    QuoteCallback
	OnTick     TickCallback
	OnDepth    DepthCallback
	OnOption   QuoteCallback
	OnFuture   QuoteCallback
	OnKline    KlineCallback
	OnStockTop QuoteCallback
	OnOptionTop QuoteCallback
	OnFullTick TickCallback
	OnQuoteBBO QuoteCallback

	// 账户推送回调
	OnAsset       AssetCallback
	OnPosition    PositionCallback
	OnOrder       OrderCallback
	OnTransaction TransactionCallback

	// 连接状态回调
	OnConnect      ConnectCallback
	OnDisconnect   DisconnectCallback
	OnError        ErrorCallback
	OnKickout      KickoutCallback
	OnSubscribe    SubscribeCallback
	OnUnsubscribe  UnsubscribeCallback
}
