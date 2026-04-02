package model

// Order 订单模型，字段名词根保持与 API JSON 一致。
type Order struct {
	// 交易账户
	Account string `json:"account,omitempty"`
	// 全局订单 ID
	ID int64 `json:"id,omitempty"`
	// 账户自增订单号
	OrderId int64 `json:"orderId,omitempty"`
	// 买卖方向（BUY/SELL）
	Action string `json:"action,omitempty"`
	// 订单类型（MKT/LMT/STP/STP_LMT/TRAIL 等）
	OrderType string `json:"orderType,omitempty"`
	// 总数量（API 返回字段名为 totalQuantity）
	TotalQuantity int64 `json:"totalQuantity,omitempty"`
	// 限价
	LimitPrice float64 `json:"limitPrice,omitempty"`
	// 辅助价格（止损价）
	AuxPrice float64 `json:"auxPrice,omitempty"`
	// 跟踪止损百分比
	TrailingPercent float64 `json:"trailingPercent,omitempty"`
	// 订单状态
	Status string `json:"status,omitempty"`
	// 已成交数量（API 返回字段名为 filledQuantity）
	FilledQuantity int64 `json:"filledQuantity,omitempty"`
	// 平均成交价
	AvgFillPrice float64 `json:"avgFillPrice,omitempty"`
	// 有效期（DAY/GTC/OPG）
	TimeInForce string `json:"timeInForce,omitempty"`
	// 是否允许盘前盘后
	OutsideRth bool `json:"outsideRth,omitempty"`
	// 附加订单（止盈/止损）
	OrderLegs []OrderLeg `json:"orderLegs,omitempty"`
	// 算法参数
	AlgoParams *AlgoParams `json:"algoParams,omitempty"`
	// 股票代码
	Symbol string `json:"symbol,omitempty"`
	// 合约类型
	SecType string `json:"secType,omitempty"`
	// 市场
	Market string `json:"market,omitempty"`
	// 货币
	Currency string `json:"currency,omitempty"`
	// 到期日（期权/期货）
	Expiry string `json:"expiry,omitempty"`
	// 行权价（期权），API 返回为字符串
	Strike string `json:"strike,omitempty"`
	// 看涨/看跌（PUT/CALL），保持 API 原始名 right
	Right string `json:"right,omitempty"`
	// 合约标识符
	Identifier string `json:"identifier,omitempty"`
	// 合约名称
	Name string `json:"name,omitempty"`
	// 佣金
	Commission float64 `json:"commission,omitempty"`
	// 已实现盈亏
	RealizedPnl float64 `json:"realizedPnl,omitempty"`
	// 开仓时间（毫秒时间戳）
	OpenTime int64 `json:"openTime,omitempty"`
	// 更新时间（毫秒时间戳）
	UpdateTime int64 `json:"updateTime,omitempty"`
	// 最新时间（毫秒时间戳）
	LatestTime int64 `json:"latestTime,omitempty"`
	// 备注
	Remark string `json:"remark,omitempty"`
	// 订单来源
	Source string `json:"source,omitempty"`
	// 用户标记
	UserMark string `json:"userMark,omitempty"`
}

// OrderLeg 附加订单（止盈/止损）
type OrderLeg struct {
	// 附加订单类型（PROFIT/LOSS）
	LegType string `json:"legType,omitempty"`
	// 价格
	Price float64 `json:"price,omitempty"`
	// 有效期
	TimeInForce string `json:"timeInForce,omitempty"`
	// 数量
	Quantity int64 `json:"quantity,omitempty"`
}

// AlgoParams 算法订单参数
type AlgoParams struct {
	// 算法策略（TWAP/VWAP）
	AlgoStrategy string `json:"algoStrategy,omitempty"`
	// 开始时间
	StartTime string `json:"startTime,omitempty"`
	// 结束时间
	EndTime string `json:"endTime,omitempty"`
	// 参与率
	ParticipationRate float64 `json:"participationRate,omitempty"`
}
