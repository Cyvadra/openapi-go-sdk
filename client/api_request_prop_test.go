package client

import (
	"encoding/json"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 11: API 请求构造正确性
// **Validates: Requirements 4.1-4.12, 5.1-5.12**
// 对于任意有效的业务参数和 API 方法名，构造的请求对象应包含正确的 method 字段，
// 且 biz_content 字段为业务参数的 JSON 序列化结果。
func TestProperty11_ApiRequestConstruction(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机 API 方法名
		methods := []string{
			"market_state", "brief", "kline", "timeline",
			"trade_tick", "quote_depth", "place_order",
			"modify_order", "cancel_order", "positions",
			"assets", "orders", "contract", "contracts",
		}
		methodIdx := rapid.IntRange(0, len(methods)-1).Draw(t, "methodIdx")
		method := methods[methodIdx]

		// 生成随机业务参数
		numKeys := rapid.IntRange(0, 5).Draw(t, "numKeys")
		bizParams := make(map[string]string)
		for i := 0; i < numKeys; i++ {
			key := rapid.StringMatching(`[a-z_]{1,15}`).Draw(t, "key")
			value := rapid.StringMatching(`[a-zA-Z0-9]{1,20}`).Draw(t, "value")
			bizParams[key] = value
		}

		// 创建请求
		req, err := NewApiRequest(method, bizParams)
		if err != nil {
			t.Fatalf("创建请求失败: %v", err)
		}

		// 验证 method 字段
		if req.Method != method {
			t.Fatalf("method 不匹配: 期望 %s，实际 %s", method, req.Method)
		}

		// 验证 biz_content 是有效 JSON
		if !json.Valid([]byte(req.BizContent)) {
			t.Fatalf("biz_content 不是有效 JSON: %s", req.BizContent)
		}

		// 反序列化验证内容一致性
		var parsed map[string]string
		if err := json.Unmarshal([]byte(req.BizContent), &parsed); err != nil {
			t.Fatalf("反序列化 biz_content 失败: %v", err)
		}
		for k, v := range bizParams {
			if parsed[k] != v {
				t.Fatalf("参数 %s 不匹配: 期望 %s，实际 %s", k, v, parsed[k])
			}
		}
	})
}
