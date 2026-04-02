package model

// Contract 合约模型，字段名词根保持与 API JSON 一致。
// Go struct 字段用 PascalCase，json tag 用 API 原始 camelCase。
type Contract struct {
	// 合约 ID
	ContractId int64 `json:"contractId,omitempty"`
	// 标的代码（如 AAPL）
	Symbol string `json:"symbol"`
	// 证券类型（STK/OPT/FUT/WAR/CASH/FUND 等）
	SecType string `json:"secType"`
	// 货币
	Currency string `json:"currency,omitempty"`
	// 交易所
	Exchange string `json:"exchange,omitempty"`
	// 到期日（期权/期货）
	Expiry string `json:"expiry,omitempty"`
	// 行权价（期权）
	Strike float64 `json:"strike,omitempty"`
	// 看涨/看跌（PUT/CALL），保持 API 原始名 right，不改为 putCall
	Right string `json:"right,omitempty"`
	// 合约乘数
	Multiplier float64 `json:"multiplier,omitempty"`
	// 合约标识符
	Identifier string `json:"identifier,omitempty"`
	// 合约名称
	Name string `json:"name,omitempty"`
	// 市场
	Market string `json:"market,omitempty"`
	// 是否可交易，保持 API 原始名 tradeable，不改为 trade
	Tradeable bool `json:"tradeable,omitempty"`
	// 合约内部 ID，保持 API 原始名 conid，不改为 contractId
	Conid int64 `json:"conid,omitempty"`
	// 做空保证金比例
	ShortMargin float64 `json:"shortMargin,omitempty"`
	// 做空初始保证金比例
	ShortInitialMargin float64 `json:"shortInitialMargin,omitempty"`
	// 做空维持保证金比例
	ShortMaintenanceMargin float64 `json:"shortMaintenanceMargin,omitempty"`
	// 做多初始保证金
	LongInitialMargin float64 `json:"longInitialMargin,omitempty"`
	// 做多维持保证金
	LongMaintenanceMargin float64 `json:"longMaintenanceMargin,omitempty"`
	// 最小报价单位价格区间
	TickSizes []TickSize `json:"tickSizes,omitempty"`
	// 每手数量
	LotSize float64 `json:"lotSize,omitempty"`
}

// TickSize 最小报价单位价格区间
type TickSize struct {
	// 起始价格（字符串，如 "0"、"1"）
	Begin string `json:"begin,omitempty"`
	// 结束价格（字符串，如 "1"、"Infinity"）
	End string `json:"end,omitempty"`
	// 最小报价单位
	TickSize float64 `json:"tickSize,omitempty"`
	// 类型（CLOSED/OPEN）
	Type string `json:"type,omitempty"`
}
