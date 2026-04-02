package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"pgregory.net/rapid"
)

// 测试 Order JSON 序列化 round-trip（单元测试）
func TestOrderJSONRoundTrip(t *testing.T) {
	original := Order{
		Account:        "U123456",
		ID:             1001,
		OrderId:        2001,
		Action:         "BUY",
		OrderType:      string(OrderTypeLMT),
		TotalQuantity:  100,
		LimitPrice:     150.50,
		Status:         string(OrderStatusSubmitted),
		TimeInForce:    string(TimeInForceDAY),
		OutsideRth:     true,
		Symbol:         "AAPL",
		SecType:        string(SecTypeSTK),
		Market:         "US",
		Currency:       "USD",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
	}
}

// 测试 Order JSON 字段名与 API JSON 一致
func TestOrderJSONFieldNames(t *testing.T) {
	o := Order{
		Account:         "U123456",
		ID:              1001,
		OrderId:         2001,
		Action:          "BUY",
		OrderType:       "LMT",
		TotalQuantity:   100,
		LimitPrice:      150.50,
		AuxPrice:        145.00,
		TrailingPercent: 5.0,
		Status:          "Submitted",
		FilledQuantity:  50,
		AvgFillPrice:    150.25,
		TimeInForce:     "DAY",
		OutsideRth:      true,
		Symbol:          "AAPL",
		SecType:         "STK",
		Market:          "US",
		Currency:        "USD",
		Expiry:          "20250620",
		Strike:          "150.0",
		Right:           "CALL",
	}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("解析为 map 失败: %v", err)
	}

	// 验证关键字段名保持 API 原始名
	expectedFields := []string{
		"account", "id", "orderId", "action", "orderType",
		"totalQuantity", "limitPrice", "auxPrice", "trailingPercent",
		"status", "filledQuantity", "avgFillPrice", "timeInForce",
		"outsideRth", "symbol", "secType", "market", "currency",
		"expiry", "strike", "right",
	}
	for _, field := range expectedFields {
		if _, ok := m[field]; !ok {
			t.Errorf("缺少 API 字段名 %q，JSON: %s", field, string(data))
		}
	}

	// 验证 right 字段不被改名
	forbiddenFields := []string{"putCall", "put_call", "order_type", "limit_price"}
	for _, field := range forbiddenFields {
		if _, ok := m[field]; ok {
			t.Errorf("不应存在改名后的字段 %q", field)
		}
	}
}

// 测试 Order 含 OrderLeg 和 AlgoParams 的 JSON round-trip
func TestOrderWithLegsAndAlgoRoundTrip(t *testing.T) {
	original := Order{
		Account:       "U123456",
		Action:        "BUY",
		OrderType:     "LMT",
		TotalQuantity: 100,
		LimitPrice:    150.50,
		TimeInForce:   "DAY",
		Symbol:        "AAPL",
		SecType:       "STK",
		OrderLegs: []OrderLeg{
			{
				LegType:     "PROFIT",
				Price:       160.0,
				TimeInForce: "GTC",
			},
			{
				LegType:     "LOSS",
				Price:       140.0,
				TimeInForce: "GTC",
			},
		},
		AlgoParams: &AlgoParams{
			AlgoStrategy:      "TWAP",
			StartTime:         "09:30:00",
			EndTime:           "16:00:00",
			ParticipationRate: 0.1,
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	if !reflect.DeepEqual(decoded, original) {
		t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
	}
}

// 测试从 API JSON 反序列化 Order（使用真实 API 字段名）
func TestOrderFromAPIJSON(t *testing.T) {
	apiJSON := `{
		"account": "U123456",
		"id": 5001,
		"orderId": 6001,
		"action": "SELL",
		"orderType": "MKT",
		"totalQuantity": 50,
		"status": "Filled",
		"filledQuantity": 50,
		"avgFillPrice": 155.30,
		"timeInForce": "DAY",
		"outsideRth": false,
		"symbol": "TSLA",
		"secType": "STK",
		"market": "US",
		"right": "PUT",
		"commission": 1.99,
		"realizedPnl": 100.50,
		"name": "Tesla Inc",
		"identifier": "TSLA",
		"source": "openapi",
		"userMark": "test001"
	}`

	var o Order
	if err := json.Unmarshal([]byte(apiJSON), &o); err != nil {
		t.Fatalf("反序列化 API JSON 失败: %v", err)
	}

	if o.Account != "U123456" {
		t.Errorf("Account = %q, want U123456", o.Account)
	}
	if o.ID != 5001 {
		t.Errorf("ID = %d, want 5001", o.ID)
	}
	if o.TotalQuantity != 50 {
		t.Errorf("TotalQuantity = %d, want 50", o.TotalQuantity)
	}
	if o.FilledQuantity != 50 {
		t.Errorf("FilledQuantity = %d, want 50", o.FilledQuantity)
	}
	if o.AvgFillPrice != 155.30 {
		t.Errorf("AvgFillPrice = %f, want 155.30", o.AvgFillPrice)
	}
	if o.Right != "PUT" {
		t.Errorf("Right = %q, want PUT", o.Right)
	}
	if o.Commission != 1.99 {
		t.Errorf("Commission = %f, want 1.99", o.Commission)
	}
	if o.Name != "Tesla Inc" {
		t.Errorf("Name = %q, want Tesla Inc", o.Name)
	}
}

// Property 7 属性测试：Order JSON round-trip
func TestOrderJSONRoundTripProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := Order{
			Account:         rapid.StringMatching(`[A-Z][0-9]{6}`).Draw(t, "account"),
			ID:              rapid.Int64Range(0, 999999).Draw(t, "id"),
			OrderId:         rapid.Int64Range(0, 999999).Draw(t, "orderId"),
			Action:          rapid.SampledFrom([]string{"BUY", "SELL"}).Draw(t, "action"),
			OrderType:       rapid.SampledFrom([]string{"MKT", "LMT", "STP", "STP_LMT", "TRAIL"}).Draw(t, "orderType"),
			TotalQuantity:   rapid.Int64Range(1, 10000).Draw(t, "totalQuantity"),
			LimitPrice:      rapid.Float64Range(0, 10000).Draw(t, "limitPrice"),
			AuxPrice:        rapid.Float64Range(0, 10000).Draw(t, "auxPrice"),
			TrailingPercent: rapid.Float64Range(0, 100).Draw(t, "trailingPercent"),
			Status:          rapid.SampledFrom([]string{"PendingNew", "Submitted", "Filled", "Cancelled"}).Draw(t, "status"),
			FilledQuantity:  rapid.Int64Range(0, 10000).Draw(t, "filledQuantity"),
			AvgFillPrice:    rapid.Float64Range(0, 10000).Draw(t, "avgFillPrice"),
			TimeInForce:     rapid.SampledFrom([]string{"DAY", "GTC", "OPG"}).Draw(t, "timeInForce"),
			OutsideRth:      rapid.Bool().Draw(t, "outsideRth"),
			Symbol:          rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			SecType:         rapid.SampledFrom([]string{"STK", "OPT", "FUT"}).Draw(t, "secType"),
			Market:          rapid.SampledFrom([]string{"US", "HK", "CN", "SG"}).Draw(t, "market"),
			Currency:        rapid.SampledFrom([]string{"USD", "HKD", "CNH", "SGD"}).Draw(t, "currency"),
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		var decoded Order
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		if !reflect.DeepEqual(decoded, original) {
			t.Errorf("round-trip 不一致:\n原始: %+v\n解码: %+v", original, decoded)
		}

		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("解析为 map 失败: %v", err)
		}

		if _, ok := m["putCall"]; ok {
			t.Error("JSON 中不应出现 putCall 字段")
		}
		if _, ok := m["orderType"]; !ok {
			t.Error("JSON 中应包含 orderType 字段")
		}
	})
}
