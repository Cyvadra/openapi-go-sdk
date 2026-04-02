package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"pgregory.net/rapid"
)

// 测试 Contract JSON 序列化 round-trip（单元测试）
func TestContractJSONRoundTrip(t *testing.T) {
	original := Contract{
		ContractId: 1234,
		Symbol:     "AAPL",
		SecType:    string(SecTypeSTK),
		Currency:   "USD",
		Exchange:   "SMART",
		Market:     "US",
		Tradeable:  true,
		Name:       "Apple Inc",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded Contract
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
	}
}

// 测试 Contract JSON 字段名与 API JSON 一致
func TestContractJSONFieldNames(t *testing.T) {
	c := Contract{
		ContractId: 100,
		Symbol:     "AAPL",
		SecType:    "STK",
		Currency:   "USD",
		Exchange:   "SMART",
		Expiry:     "20250620",
		Strike:     150.0,
		Right:      "CALL",
		Multiplier: 100.0,
		Identifier: "AAPL 250620C00150000",
		Name:       "Apple Inc",
		Market:     "US",
		Tradeable:  true,
		Conid:      12345,
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 解析为 map 检查字段名
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析为 map 失败: %v", err)
	}

	// 验证关键字段名保持 API 原始名
	expectedFields := []string{
		"contractId", "symbol", "secType", "currency", "exchange",
		"expiry", "strike", "right", "multiplier", "identifier",
		"name", "market", "tradeable", "conid",
	}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("缺少 API 字段名 %q，JSON: %s", field, string(data))
		}
	}

	// 验证不存在被改名的字段
	forbiddenFields := []string{"putCall", "put_call", "contractID", "contract_id", "trade"}
	for _, field := range forbiddenFields {
		if _, ok := m[field]; ok {
			t.Errorf("不应存在改名后的字段 %q，JSON: %s", field, string(data))
		}
	}
}

// 测试 Contract omitempty：零值字段不出现在 JSON 中
func TestContractOmitEmpty(t *testing.T) {
	c := Contract{
		Symbol:  "AAPL",
		SecType: "STK",
	}

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析为 map 失败: %v", err)
	}

	// symbol 和 secType 必须存在
	if _, ok := m["symbol"]; !ok {
		t.Error("symbol 字段应始终存在")
	}
	if _, ok := m["secType"]; !ok {
		t.Error("secType 字段应始终存在")
	}

	// 可选字段零值不应出现
	optionalFields := []string{"contractId", "currency", "exchange", "expiry", "strike", "right", "multiplier"}
	for _, field := range optionalFields {
		if _, ok := m[field]; ok {
			t.Errorf("零值可选字段 %q 不应出现在 JSON 中", field)
		}
	}
}

// 测试从 API JSON 反序列化 Contract
func TestContractFromAPIJSON(t *testing.T) {
	apiJSON := `{
		"contractId": 9876,
		"symbol": "TSLA",
		"secType": "STK",
		"currency": "USD",
		"exchange": "SMART",
		"market": "US",
		"tradeable": true,
		"right": "CALL",
		"conid": 54321,
		"name": "Tesla Inc"
	}`

	var c Contract
	if err := json.Unmarshal([]byte(apiJSON), &c); err != nil {
		t.Fatalf("反序列化 API JSON 失败: %v", err)
	}

	if c.ContractId != 9876 {
		t.Errorf("ContractId = %d, want 9876", c.ContractId)
	}
	if c.Symbol != "TSLA" {
		t.Errorf("Symbol = %q, want TSLA", c.Symbol)
	}
	if c.Right != "CALL" {
		t.Errorf("Right = %q, want CALL（字段名保持 right 不改为 putCall）", c.Right)
	}
	if c.Conid != 54321 {
		t.Errorf("Conid = %d, want 54321（字段名保持 conid 不改为 contractId）", c.Conid)
	}
	if c.Tradeable != true {
		t.Errorf("Tradeable = %v, want true（字段名保持 tradeable 不改为 trade）", c.Tradeable)
	}
}

// 2.3.1 Property 7 属性测试：Contract JSON round-trip
// Feature: multi-language-sdks, Property 7: 数据模型 JSON 序列化 round-trip
// **Validates: Requirements 7.1, 7.7**
func TestContractJSONRoundTripProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := Contract{
			ContractId: rapid.Int64Range(0, 999999).Draw(t, "contractId"),
			Symbol:     rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			SecType:    rapid.SampledFrom([]string{"STK", "OPT", "FUT", "WAR", "CASH", "FUND"}).Draw(t, "secType"),
			Currency:   rapid.SampledFrom([]string{"USD", "HKD", "CNH", "SGD"}).Draw(t, "currency"),
			Exchange:   rapid.SampledFrom([]string{"SMART", "NYSE", "NASDAQ", "SEHK", ""}).Draw(t, "exchange"),
			Expiry:     rapid.SampledFrom([]string{"", "20250620", "20251219"}).Draw(t, "expiry"),
			Strike:     rapid.Float64Range(0, 10000).Draw(t, "strike"),
			Right:      rapid.SampledFrom([]string{"", "CALL", "PUT"}).Draw(t, "right"),
			Multiplier: rapid.Float64Range(0, 1000).Draw(t, "multiplier"),
			Identifier: rapid.SampledFrom([]string{"", "AAPL 250620C00150000"}).Draw(t, "identifier"),
			Name:       rapid.SampledFrom([]string{"", "Apple Inc", "Tesla Inc"}).Draw(t, "name"),
			Market:     rapid.SampledFrom([]string{"", "US", "HK", "CN", "SG"}).Draw(t, "market"),
			Tradeable:  rapid.Bool().Draw(t, "tradeable"),
			Conid:      rapid.Int64Range(0, 999999).Draw(t, "conid"),
		}

		// 序列化
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		// 反序列化
		var decoded Contract
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		// 验证 round-trip 等价
		if !reflect.DeepEqual(decoded, original) {
			t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
		}

		// 验证 JSON 字段名词根与 API 一致
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("解析为 map 失败: %v", err)
		}

		// right 字段不应被改名为 putCall
		if _, ok := m["putCall"]; ok {
			t.Error("JSON 中不应出现 putCall 字段，应保持 right")
		}
		// tradeable 字段不应被改名为 trade
		if _, ok := m["trade"]; ok {
			t.Error("JSON 中不应出现 trade 字段，应保持 tradeable")
		}
		// conid 字段不应被改名为 contractId（注意 contractId 是另一个独立字段）
		if original.Conid != 0 {
			if _, ok := m["conid"]; !ok {
				t.Error("JSON 中应包含 conid 字段")
			}
		}
	})
}
