package push

import (
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 12: Protobuf 序列化 round-trip
// （简化为 JSON 序列化 round-trip）
// 对于任意有效的推送数据对象，使用 JSON 序列化后再反序列化，
// 得到的对象应与原始对象等价。
// **Validates: Requirements 6.12**

func TestProperty_QuoteData_SerializeRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := &QuoteData{
			Symbol:      rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			LatestPrice: rapid.Float64Range(0, 100000).Draw(t, "latestPrice"),
			PreClose:    rapid.Float64Range(0, 100000).Draw(t, "preClose"),
			Open:        rapid.Float64Range(0, 100000).Draw(t, "open"),
			High:        rapid.Float64Range(0, 100000).Draw(t, "high"),
			Low:         rapid.Float64Range(0, 100000).Draw(t, "low"),
			Volume:      rapid.Int64Range(0, 1000000000).Draw(t, "volume"),
			Amount:      rapid.Float64Range(0, 1e12).Draw(t, "amount"),
			Timestamp:   rapid.Int64Range(0, 2000000000).Draw(t, "timestamp"),
		}

		msg, err := NewPushMessage(MsgTypeQuote, SubjectQuote, original)
		if err != nil {
			t.Fatalf("创建消息失败: %v", err)
		}

		data, err := msg.Serialize()
		if err != nil {
			t.Fatalf("序列化失败: %v", err)
		}

		msg2, err := DeserializeMessage(data)
		if err != nil {
			t.Fatalf("反序列化失败: %v", err)
		}

		var restored QuoteData
		if err := json.Unmarshal(msg2.Data, &restored); err != nil {
			t.Fatalf("解析数据失败: %v", err)
		}

		if restored.Symbol != original.Symbol {
			t.Errorf("Symbol: 期望 %s, 实际 %s", original.Symbol, restored.Symbol)
		}
		if restored.Volume != original.Volume {
			t.Errorf("Volume: 期望 %d, 实际 %d", original.Volume, restored.Volume)
		}
		if restored.Timestamp != original.Timestamp {
			t.Errorf("Timestamp: 期望 %d, 实际 %d", original.Timestamp, restored.Timestamp)
		}
	})
}

func TestProperty_TickData_SerializeRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := &TickData{
			Symbol:    rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			Price:     rapid.Float64Range(0, 100000).Draw(t, "price"),
			Volume:    rapid.Int64Range(0, 1000000).Draw(t, "volume"),
			Type:      rapid.SampledFrom([]string{"BUY", "SELL", ""}).Draw(t, "type"),
			Timestamp: rapid.Int64Range(0, 2000000000).Draw(t, "timestamp"),
		}

		msg, err := NewPushMessage(MsgTypeTick, SubjectTick, original)
		if err != nil {
			t.Fatalf("创建消息失败: %v", err)
		}
		data, _ := msg.Serialize()
		msg2, _ := DeserializeMessage(data)

		var restored TickData
		json.Unmarshal(msg2.Data, &restored)

		if restored.Symbol != original.Symbol {
			t.Errorf("Symbol 不匹配")
		}
		if restored.Price != original.Price {
			t.Errorf("Price 不匹配")
		}
		if restored.Volume != original.Volume {
			t.Errorf("Volume 不匹配")
		}
	})
}

func TestProperty_OrderData_SerializeRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := &OrderData{
			Account:   rapid.StringMatching(`[a-z0-9]{5,10}`).Draw(t, "account"),
			ID:        rapid.Int64Range(1, 1000000).Draw(t, "id"),
			OrderID:   rapid.Int64Range(1, 1000000).Draw(t, "orderId"),
			Symbol:    rapid.StringMatching(`[A-Z]{1,5}`).Draw(t, "symbol"),
			Action:    rapid.SampledFrom([]string{"BUY", "SELL"}).Draw(t, "action"),
			OrderType: rapid.SampledFrom([]string{"MKT", "LMT", "STP"}).Draw(t, "orderType"),
			Quantity:  rapid.IntRange(1, 10000).Draw(t, "quantity"),
			Status:    rapid.SampledFrom([]string{"Submitted", "Filled", "Cancelled"}).Draw(t, "status"),
		}

		msg, _ := NewPushMessage(MsgTypeOrder, SubjectOrder, original)
		data, _ := msg.Serialize()
		msg2, _ := DeserializeMessage(data)

		var restored OrderData
		json.Unmarshal(msg2.Data, &restored)

		if restored.Account != original.Account || restored.ID != original.ID ||
			restored.Symbol != original.Symbol || restored.Status != original.Status ||
			restored.Quantity != original.Quantity {
			t.Errorf("OrderData round-trip 失败: 原始=%+v, 恢复=%+v", original, restored)
		}
	})
}

func TestProperty_AssetData_SerializeRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		original := &AssetData{
			Account:        rapid.StringMatching(`[a-z0-9]{5,10}`).Draw(t, "account"),
			NetLiquidation: rapid.Float64Range(0, 1e8).Draw(t, "netLiq"),
			CashBalance:    rapid.Float64Range(0, 1e8).Draw(t, "cash"),
			BuyingPower:    rapid.Float64Range(0, 1e8).Draw(t, "bp"),
			Currency:       rapid.SampledFrom([]string{"USD", "HKD", "CNH"}).Draw(t, "currency"),
		}

		msg, _ := NewPushMessage(MsgTypeAsset, SubjectAsset, original)
		data, _ := msg.Serialize()
		msg2, _ := DeserializeMessage(data)

		var restored AssetData
		json.Unmarshal(msg2.Data, &restored)

		if restored.Account != original.Account || restored.Currency != original.Currency {
			t.Errorf("AssetData round-trip 失败")
		}
	})
}
