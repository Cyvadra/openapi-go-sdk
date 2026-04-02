package model

import (
	"testing"
)

// 测试 StockContract 构造函数
func TestStockContract(t *testing.T) {
	c := StockContract("AAPL", "USD")
	if c.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", c.Symbol)
	}
	if c.SecType != string(SecTypeSTK) {
		t.Errorf("SecType = %q, want STK", c.SecType)
	}
	if c.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", c.Currency)
	}
}

// 测试 OptionContract 通过 identifier 构造
func TestOptionContract(t *testing.T) {
	c := OptionContract("AAPL 250620C00150000")
	if c.Identifier != "AAPL 250620C00150000" {
		t.Errorf("Identifier = %q, want AAPL 250620C00150000", c.Identifier)
	}
	if c.SecType != string(SecTypeOPT) {
		t.Errorf("SecType = %q, want OPT", c.SecType)
	}
}

// 测试 OptionContractBySymbol 通过各字段构造
func TestOptionContractBySymbol(t *testing.T) {
	c := OptionContractBySymbol("AAPL", "20250620", 150.0, "CALL", "USD")
	if c.Symbol != "AAPL" {
		t.Errorf("Symbol = %q, want AAPL", c.Symbol)
	}
	if c.SecType != string(SecTypeOPT) {
		t.Errorf("SecType = %q, want OPT", c.SecType)
	}
	if c.Expiry != "20250620" {
		t.Errorf("Expiry = %q, want 20250620", c.Expiry)
	}
	if c.Strike != 150.0 {
		t.Errorf("Strike = %f, want 150.0", c.Strike)
	}
	if c.Right != "CALL" {
		t.Errorf("Right = %q, want CALL", c.Right)
	}
	if c.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", c.Currency)
	}
}

// 测试 FutureContract 构造函数
func TestFutureContract(t *testing.T) {
	c := FutureContract("ES", "USD", "20251219")
	if c.Symbol != "ES" {
		t.Errorf("Symbol = %q, want ES", c.Symbol)
	}
	if c.SecType != string(SecTypeFUT) {
		t.Errorf("SecType = %q, want FUT", c.SecType)
	}
	if c.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", c.Currency)
	}
	if c.Expiry != "20251219" {
		t.Errorf("Expiry = %q, want 20251219", c.Expiry)
	}
}

// 测试 CashContract 构造函数
func TestCashContract(t *testing.T) {
	c := CashContract("USD.HKD")
	if c.Symbol != "USD.HKD" {
		t.Errorf("Symbol = %q, want USD.HKD", c.Symbol)
	}
	if c.SecType != string(SecTypeCASH) {
		t.Errorf("SecType = %q, want CASH", c.SecType)
	}
}

// 测试 FundContract 构造函数
func TestFundContract(t *testing.T) {
	c := FundContract("SPY", "USD")
	if c.Symbol != "SPY" {
		t.Errorf("Symbol = %q, want SPY", c.Symbol)
	}
	if c.SecType != string(SecTypeFUND) {
		t.Errorf("SecType = %q, want FUND", c.SecType)
	}
	if c.Currency != "USD" {
		t.Errorf("Currency = %q, want USD", c.Currency)
	}
}

// 测试 WarrantContract 构造函数
func TestWarrantContract(t *testing.T) {
	c := WarrantContract("00700", "HKD", "20251219", 400.0, "CALL")
	if c.Symbol != "00700" {
		t.Errorf("Symbol = %q, want 00700", c.Symbol)
	}
	if c.SecType != string(SecTypeWAR) {
		t.Errorf("SecType = %q, want WAR", c.SecType)
	}
	if c.Currency != "HKD" {
		t.Errorf("Currency = %q, want HKD", c.Currency)
	}
	if c.Expiry != "20251219" {
		t.Errorf("Expiry = %q, want 20251219", c.Expiry)
	}
	if c.Strike != 400.0 {
		t.Errorf("Strike = %f, want 400.0", c.Strike)
	}
	if c.Right != "CALL" {
		t.Errorf("Right = %q, want CALL", c.Right)
	}
}
