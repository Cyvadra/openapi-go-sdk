package model

// MarketOrder 构造市价单
func MarketOrder(account, symbol, secType, action string, quantity int64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeMKT),
		TotalQuantity: quantity,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// LimitOrder 构造限价单
func LimitOrder(account, symbol, secType, action string, quantity int64, limitPrice float64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeLMT),
		TotalQuantity: quantity,
		LimitPrice:    limitPrice,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// StopOrder 构造止损单
func StopOrder(account, symbol, secType, action string, quantity int64, auxPrice float64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeSTP),
		TotalQuantity: quantity,
		AuxPrice:      auxPrice,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// StopLimitOrder 构造止损限价单
func StopLimitOrder(account, symbol, secType, action string, quantity int64, limitPrice, auxPrice float64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeSTPLMT),
		TotalQuantity: quantity,
		LimitPrice:    limitPrice,
		AuxPrice:      auxPrice,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// TrailOrder 构造跟踪止损单
func TrailOrder(account, symbol, secType, action string, quantity int64, trailingPercent float64) Order {
	return Order{
		Account:         account,
		Symbol:          symbol,
		SecType:         secType,
		Action:          action,
		OrderType:       string(OrderTypeTRAIL),
		TotalQuantity:   quantity,
		TrailingPercent: trailingPercent,
		TimeInForce:     string(TimeInForceDAY),
	}
}

// AuctionLimitOrder 构造竞价限价单
func AuctionLimitOrder(account, symbol, secType, action string, quantity int64, limitPrice float64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeAL),
		TotalQuantity: quantity,
		LimitPrice:    limitPrice,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// AuctionMarketOrder 构造竞价市价单
func AuctionMarketOrder(account, symbol, secType, action string, quantity int64) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     string(OrderTypeAM),
		TotalQuantity: quantity,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// AlgoOrder 构造算法订单（TWAP/VWAP）
func AlgoOrder(account, symbol, secType, action string, quantity int64, limitPrice float64, algoType string, params AlgoParams) Order {
	return Order{
		Account:       account,
		Symbol:        symbol,
		SecType:       secType,
		Action:        action,
		OrderType:     algoType,
		TotalQuantity: quantity,
		LimitPrice:    limitPrice,
		AlgoParams:    &params,
		TimeInForce:   string(TimeInForceDAY),
	}
}

// NewOrderLeg 构造附加订单（止盈/止损）
func NewOrderLeg(legType string, price float64, timeInForce string) OrderLeg {
	return OrderLeg{
		LegType:     legType,
		Price:       price,
		TimeInForce: timeInForce,
	}
}
