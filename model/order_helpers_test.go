package model

import (
	"testing"
)

// 测试 MarketOrder 构造函数
func TestMarketOrder(t *testing.T) {
	o := MarketOrder("U123456", "AAPL", "STK", "BUY", 100)
	if o.Account != "U123456" {
		t.Errorf("Account = %q, want U123456", o.Account)
	}
	if o.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", o.Symbol)
	}
	if o.SecType != "STK" {
		t.Errorf("SecType = %q, want STK", o.SecType)
	}
	if o.Action != "BUY" {
		t.Errorf("Action = %q, want BUY", o.Action)
	}
	if o.OrderType != string(OrderTypeMKT) {
		t.Errorf("OrderType = %q, want MKT", o.OrderType)
	}
	if o.TotalQuantity != 100 {
		t.Errorf("TotalQuantity = %d, want 100", o.TotalQuantity)
	}
	if o.TimeInForce != string(TimeInForceDAY) {
		t.Errorf("TimeInForce = %q, want DAY", o.TimeInForce)
	}
}

// 测试 LimitOrder 构造函数
func TestLimitOrder(t *testing.T) {
	o := LimitOrder("U123456", "AAPL", "STK", "BUY", 100, 150.50)
	if o.OrderType != string(OrderTypeLMT) {
		t.Errorf("OrderType = %q, want LMT", o.OrderType)
	}
	if o.LimitPrice != 150.50 {
		t.Errorf("LimitPrice = %f, want 150.50", o.LimitPrice)
	}
	if o.TimeInForce != string(TimeInForceDAY) {
		t.Errorf("TimeInForce = %q, want DAY", o.TimeInForce)
	}
}

// 测试 StopOrder 构造函数
func TestStopOrder(t *testing.T) {
	o := StopOrder("U123456", "AAPL", "STK", "SELL", 100, 140.00)
	if o.OrderType != string(OrderTypeSTP) {
		t.Errorf("OrderType = %q, want STP", o.OrderType)
	}
	if o.AuxPrice != 140.00 {
		t.Errorf("AuxPrice = %f, want 140.00", o.AuxPrice)
	}
}

// 测试 StopLimitOrder 构造函数
func TestStopLimitOrder(t *testing.T) {
	o := StopLimitOrder("U123456", "AAPL", "STK", "SELL", 100, 145.00, 140.00)
	if o.OrderType != string(OrderTypeSTPLMT) {
		t.Errorf("OrderType = %q, want STP_LMT", o.OrderType)
	}
	if o.LimitPrice != 145.00 {
		t.Errorf("LimitPrice = %f, want 145.00", o.LimitPrice)
	}
	if o.AuxPrice != 140.00 {
		t.Errorf("AuxPrice = %f, want 140.00", o.AuxPrice)
	}
}

// 测试 TrailOrder 构造函数
func TestTrailOrder(t *testing.T) {
	o := TrailOrder("U123456", "AAPL", "STK", "SELL", 100, 5.0)
	if o.OrderType != string(OrderTypeTRAIL) {
		t.Errorf("OrderType = %q, want TRAIL", o.OrderType)
	}
	if o.TrailingPercent != 5.0 {
		t.Errorf("TrailingPercent = %f, want 5.0", o.TrailingPercent)
	}
}

// 测试 AuctionLimitOrder 构造函数
func TestAuctionLimitOrder(t *testing.T) {
	o := AuctionLimitOrder("U123456", "AAPL", "STK", "BUY", 100, 150.00)
	if o.OrderType != string(OrderTypeAL) {
		t.Errorf("OrderType = %q, want AL", o.OrderType)
	}
	if o.LimitPrice != 150.00 {
		t.Errorf("LimitPrice = %f, want 150.00", o.LimitPrice)
	}
}

// 测试 AuctionMarketOrder 构造函数
func TestAuctionMarketOrder(t *testing.T) {
	o := AuctionMarketOrder("U123456", "AAPL", "STK", "BUY", 100)
	if o.OrderType != string(OrderTypeAM) {
		t.Errorf("OrderType = %q, want AM", o.OrderType)
	}
}

// 测试 AlgoOrder 构造函数
func TestAlgoOrder(t *testing.T) {
	params := AlgoParams{
		AlgoStrategy:      "TWAP",
		StartTime:         "09:30:00",
		EndTime:           "16:00:00",
		ParticipationRate: 0.1,
	}
	o := AlgoOrder("U123456", "AAPL", "STK", "BUY", 1000, 150.00, "TWAP", params)
	if o.OrderType != string(OrderTypeTWAP) {
		t.Errorf("OrderType = %q, want TWAP", o.OrderType)
	}
	if o.LimitPrice != 150.00 {
		t.Errorf("LimitPrice = %f, want 150.00", o.LimitPrice)
	}
	if o.AlgoParams == nil {
		t.Fatal("AlgoParams 不应为 nil")
	}
	if o.AlgoParams.AlgoStrategy != "TWAP" {
		t.Errorf("AlgoStrategy = %q, want TWAP", o.AlgoParams.AlgoStrategy)
	}
	if o.AlgoParams.ParticipationRate != 0.1 {
		t.Errorf("ParticipationRate = %f, want 0.1", o.AlgoParams.ParticipationRate)
	}
}

// 测试 NewOrderLeg 构造函数
func TestNewOrderLeg(t *testing.T) {
	leg := NewOrderLeg("PROFIT", 160.00, "GTC")
	if leg.LegType != "PROFIT" {
		t.Errorf("LegType = %q, want PROFIT", leg.LegType)
	}
	if leg.Price != 160.00 {
		t.Errorf("Price = %f, want 160.00", leg.Price)
	}
	if leg.TimeInForce != "GTC" {
		t.Errorf("TimeInForce = %q, want GTC", leg.TimeInForce)
	}
}
