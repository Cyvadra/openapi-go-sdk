package model

// Position 持仓模型，字段名词根保持与 API JSON 一致。
type Position struct {
	// 交易账户
	Account string `json:"account,omitempty"`
	// 标的代码
	Symbol string `json:"symbol,omitempty"`
	// 证券类型
	SecType string `json:"secType,omitempty"`
	// 市场
	Market string `json:"market,omitempty"`
	// 货币
	Currency string `json:"currency,omitempty"`
	// 持仓数量（API 返回字段名为 position）
	Position int64 `json:"position,omitempty"`
	// 平均成本
	AverageCost float64 `json:"averageCost,omitempty"`
	// 市值
	MarketValue float64 `json:"marketValue,omitempty"`
	// 已实现盈亏
	RealizedPnl float64 `json:"realizedPnl,omitempty"`
	// 未实现盈亏
	UnrealizedPnl float64 `json:"unrealizedPnl,omitempty"`
	// 未实现盈亏百分比
	UnrealizedPnlPercent float64 `json:"unrealizedPnlPercent,omitempty"`
	// 合约 ID
	ContractId int64 `json:"contractId,omitempty"`
	// 合约标识符
	Identifier string `json:"identifier,omitempty"`
	// 合约名称
	Name string `json:"name,omitempty"`
	// 最新价格
	LatestPrice float64 `json:"latestPrice,omitempty"`
	// 合约乘数
	Multiplier float64 `json:"multiplier,omitempty"`
}
