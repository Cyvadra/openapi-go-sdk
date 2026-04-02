package model

import (
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
)

// 测试 Position JSON 序列化 round-trip（单元测试）
func TestPositionJSONRoundTrip(t *testing.T) {
	original := Position{
		Account:       "U123456",
		Symbol:        "AAPL",
		SecType:       "STK",
		Market:        "US",
		Currency:      "USD",
		Position:      100,
		AverageCost:   150.25,
		MarketValue:   15530.0,
		RealizedPnl:   200.0,
		UnrealizedPnl: 505.0,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded Position
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if decoded != original {
		t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
	}
}

// 测试 Position JSON 字段名与 API JSON 一致
func TestPositionJSONFieldNames(t *testing.T) {
	p := Position{
		Account:              "U123456",
		Symbol:               "AAPL",
		SecType:              "STK",
		Market:               "US",
		Currency:             "USD",
		Position:             100,
		AverageCost:          150.25,
		MarketValue:          15530.0,
		RealizedPnl:          200.0,
		UnrealizedPnl:        505.0,
		UnrealizedPnlPercent: 0.15,
		ContractId:           14,
		Identifier:           "AAPL",
		Name:                 "Apple Inc",
		LatestPrice:          155.30,
		Multiplier:           1.0,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析为 map 失败: %v", err)
	}

	expectedFields := []string{
		"account", "symbol", "secType", "market", "currency",
		"position", "averageCost", "marketValue",
		"realizedPnl", "unrealizedPnl", "unrealizedPnlPercent",
		"contractId", "identifier", "name", "latestPrice", "multiplier",
	}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("缺少 API 字段名 %q", field)
		}
	}
}

// 测试从 API JSON 反序列化 Position（使用真实 API 字段名）
func TestPositionFromAPIJSON(t *testing.T) {
	apiJSON := `{
		"account": "U789012",
		"symbol": "TSLA",
		"secType": "STK",
		"market": "US",
		"currency": "USD",
		"position": 50,
		"averageCost": 200.50,
		"marketValue": 10500.00,
		"realizedPnl": 0,
		"unrealizedPnl": 475.00,
		"unrealizedPnlPercent": -0.05,
		"contractId": 12345,
		"identifier": "TSLA",
		"name": "Tesla Inc",
		"latestPrice": 210.00,
		"multiplier": 1.0
	}`

	var p Position
	if err := json.Unmarshal([]byte(apiJSON), &p); err != nil {
		t.Fatalf("反序列化 API JSON 失败: %v", err)
	}

	if p.Account != "U789012" {
		t.Errorf("Account = %q, want U789012", p.Account)
	}
	if p.Position != 50 {
		t.Errorf("Position = %d, want 50", p.Position)
	}
	if p.AverageCost != 200.50 {
		t.Errorf("AverageCost = %f, want 200.50", p.AverageCost)
	}
	if p.UnrealizedPnl != 475.00 {
		t.Errorf("UnrealizedPnl = %f, want 475.00", p.UnrealizedPnl)
	}
	if p.ContractId != 12345 {
		t.Errorf("ContractId = %d, want 12345", p.ContractId)
	}
	if p.Name != "Tesla Inc" {
		t.Errorf("Name = %q, want Tesla Inc", p.Name)
	}
}

// Property 7 属性测试：Position JSON round-trip
func TestPositionJSONRoundTripProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := Position{
			Account:       rapid.StringMatching(`[A-Z][0-9]{6}`).Draw(t, "account"),
			Symbol:        rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			SecType:       rapid.SampledFrom([]string{"STK", "OPT", "FUT"}).Draw(t, "secType"),
			Market:        rapid.SampledFrom([]string{"US", "HK", "CN", "SG"}).Draw(t, "market"),
			Currency:      rapid.SampledFrom([]string{"USD", "HKD", "CNH", "SGD"}).Draw(t, "currency"),
			Position:      rapid.Int64Range(-10000, 10000).Draw(t, "position"),
			AverageCost:   rapid.Float64Range(0, 10000).Draw(t, "averageCost"),
			MarketValue:   rapid.Float64Range(-1000000, 1000000).Draw(t, "marketValue"),
			RealizedPnl:   rapid.Float64Range(-100000, 100000).Draw(t, "realizedPnl"),
			UnrealizedPnl: rapid.Float64Range(-100000, 100000).Draw(t, "unrealizedPnl"),
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		var decoded Position
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		if decoded != original {
			t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
		}
	})
}
