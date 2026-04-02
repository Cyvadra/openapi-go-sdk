// Package model 定义老虎证券 OpenAPI 的数据模型和枚举类型。
package model

// Market 市场枚举
type Market string

const (
	MarketAll Market = "ALL"
	MarketUS  Market = "US"
	MarketHK  Market = "HK"
	MarketCN  Market = "CN"
	MarketSG  Market = "SG"
)

// SecurityType 证券类型枚举
type SecurityType string

const (
	SecTypeAll  SecurityType = "ALL"
	SecTypeSTK  SecurityType = "STK"
	SecTypeOPT  SecurityType = "OPT"
	SecTypeWAR  SecurityType = "WAR"
	SecTypeIOPT SecurityType = "IOPT"
	SecTypeFUT  SecurityType = "FUT"
	SecTypeFOP  SecurityType = "FOP"
	SecTypeCASH SecurityType = "CASH"
	SecTypeMLEG SecurityType = "MLEG"
	SecTypeFUND SecurityType = "FUND"
)

// Currency 货币枚举
type Currency string

const (
	CurrencyAll Currency = "ALL"
	CurrencyUSD Currency = "USD"
	CurrencyHKD Currency = "HKD"
	CurrencyCNH Currency = "CNH"
	CurrencySGD Currency = "SGD"
)

// OrderType 订单类型枚举
type OrderType string

const (
	OrderTypeMKT    OrderType = "MKT"
	OrderTypeLMT    OrderType = "LMT"
	OrderTypeSTP    OrderType = "STP"
	OrderTypeSTPLMT OrderType = "STP_LMT"
	OrderTypeTRAIL  OrderType = "TRAIL"
	OrderTypeAM     OrderType = "AM"
	OrderTypeAL     OrderType = "AL"
	OrderTypeTWAP   OrderType = "TWAP"
	OrderTypeVWAP   OrderType = "VWAP"
	OrderTypeOCA    OrderType = "OCA"
)

// OrderStatus 订单状态枚举
type OrderStatus string

const (
	OrderStatusPendingNew      OrderStatus = "PendingNew"
	OrderStatusInitial         OrderStatus = "Initial"
	OrderStatusSubmitted       OrderStatus = "Submitted"
	OrderStatusPartiallyFilled OrderStatus = "PartiallyFilled"
	OrderStatusFilled          OrderStatus = "Filled"
	OrderStatusCancelled       OrderStatus = "Cancelled"
	OrderStatusPendingCancel   OrderStatus = "PendingCancel"
	OrderStatusInactive        OrderStatus = "Inactive"
	OrderStatusInvalid         OrderStatus = "Invalid"
)

// BarPeriod K 线周期枚举
type BarPeriod string

const (
	BarPeriodDay   BarPeriod = "day"
	BarPeriodWeek  BarPeriod = "week"
	BarPeriodMonth BarPeriod = "month"
	BarPeriodYear  BarPeriod = "year"
	BarPeriod1Min  BarPeriod = "1min"
	BarPeriod5Min  BarPeriod = "5min"
	BarPeriod15Min BarPeriod = "15min"
	BarPeriod30Min BarPeriod = "30min"
	BarPeriod60Min BarPeriod = "60min"
)

// Language 语言枚举
type Language string

const (
	LanguageZhCN Language = "zh_CN"
	LanguageZhTW Language = "zh_TW"
	LanguageEnUS Language = "en_US"
)

// QuoteRight 复权类型枚举
type QuoteRight string

const (
	QuoteRightBr QuoteRight = "br" // 前复权
	QuoteRightNr QuoteRight = "nr" // 不复权
)

// License 牌照类型枚举
type License string

const (
	LicenseTBNZ License = "TBNZ"
	LicenseTBSG License = "TBSG"
	LicenseTBHK License = "TBHK"
	LicenseTBAU License = "TBAU"
	LicenseTBUS License = "TBUS"
)

// TimeInForce 订单有效期枚举
type TimeInForce string

const (
	TimeInForceDAY TimeInForce = "DAY"
	TimeInForceGTC TimeInForce = "GTC"
	TimeInForceOPG TimeInForce = "OPG"
)
