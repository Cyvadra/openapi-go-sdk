package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 13: Generic execute 请求构造正确性
// **Validates: Requirements 15.3, 15.8**
// 对于任意有效的 API 方法名和有效的 biz_content JSON 字符串，调用通用 execute 方法时，
// 构造的 HTTP 请求参数应包含正确的公共参数，且 method 字段等于传入的 api_method，
// biz_content 字段等于传入的 request_json。
func TestProperty13_GenericExecuteRequestConstruction(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机 API 方法名
		methods := []string{
			"market_state", "brief", "kline", "timeline",
			"trade_tick", "quote_depth", "place_order",
			"positions", "assets", "orders", "contract",
		}
		methodIdx := rapid.IntRange(0, len(methods)-1).Draw(t, "methodIdx")
		apiMethod := methods[methodIdx]

		// 生成随机 biz_content JSON
		numKeys := rapid.IntRange(0, 3).Draw(t, "numKeys")
		bizMap := make(map[string]string)
		for i := 0; i < numKeys; i++ {
			key := rapid.StringMatching(`[a-z]{2,10}`).Draw(t, "key")
			value := rapid.StringMatching(`[a-zA-Z0-9]{1,15}`).Draw(t, "value")
			bizMap[key] = value
		}
		bizJSON, _ := json.Marshal(bizMap)
		requestJSON := string(bizJSON)

		// 捕获服务器收到的参数
		var receivedParams map[string]string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewDecoder(r.Body).Decode(&receivedParams)
			fmt.Fprint(w, `{"code":0,"message":"ok","data":null}`)
		}))
		defer server.Close()

		cfg := newTestConfigWithURL(server.URL)
		client := NewHttpClient(cfg)

		_, err := client.ExecuteRaw(apiMethod, requestJSON)
		if err != nil {
			t.Fatalf("ExecuteRaw 失败: %v", err)
		}

		// 验证 method 字段等于传入的 apiMethod
		if receivedParams["method"] != apiMethod {
			t.Fatalf("method 不匹配: 期望 %s，实际 %s", apiMethod, receivedParams["method"])
		}

		// 验证 biz_content 字段等于传入的 requestJSON
		if receivedParams["biz_content"] != requestJSON {
			t.Fatalf("biz_content 不匹配: 期望 %s，实际 %s", requestJSON, receivedParams["biz_content"])
		}

		// 验证必要的公共参数存在
		requiredKeys := []string{"tiger_id", "method", "charset", "sign_type", "timestamp", "version", "biz_content", "sign"}
		for _, key := range requiredKeys {
			if _, ok := receivedParams[key]; !ok {
				t.Fatalf("缺少必要参数: %s", key)
			}
		}

		// 验证 tiger_id 正确
		if receivedParams["tiger_id"] != "test_tiger_id" {
			t.Fatalf("tiger_id 不匹配: %s", receivedParams["tiger_id"])
		}

		// 验证 sign 非空
		if receivedParams["sign"] == "" {
			t.Fatal("sign 不应为空")
		}
	})
}

// Feature: multi-language-sdks, Property 14: Generic execute 响应原始透传
// **Validates: Requirements 15.4**
// 对于任意服务器返回的 JSON 响应字符串，通用 execute 方法返回的字符串应与服务器返回的原始 JSON 完全一致。
func TestProperty14_GenericExecuteRawPassthrough(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机响应 JSON
		code := rapid.IntRange(0, 5000).Draw(t, "code")
		message := rapid.StringMatching(`[a-zA-Z0-9 ]{1,30}`).Draw(t, "message")
		timestamp := rapid.Int64Range(1000000000, 2000000000).Draw(t, "timestamp")

		// 构造服务器响应
		serverResp := fmt.Sprintf(
			`{"code":%d,"message":"%s","data":{"key":"val"},"timestamp":%d}`,
			code, message, timestamp,
		)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, serverResp)
		}))
		defer server.Close()

		cfg := newTestConfigWithURL(server.URL)
		client := NewHttpClient(cfg)

		result, _ := client.ExecuteRaw("test_method", `{}`)

		// 验证返回的字符串与服务器返回的原始 JSON 完全一致
		if result != serverResp {
			t.Fatalf("响应不匹配:\n期望: %s\n实际: %s", serverResp, result)
		}
	})
}

// Feature: multi-language-sdks, Property 15: Generic execute 无效 JSON 拒绝
// **Validates: Requirements 15.5, 15.6**
// 对于任意非法 JSON 字符串，调用通用 execute 方法时应返回参数错误，不发送任何 HTTP 请求。
func TestProperty15_GenericExecuteRejectInvalidJSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成不是有效 JSON 的字符串
		invalidJSON := rapid.StringMatching(`[a-zA-Z]{5,30}`).Draw(t, "invalidJSON")

		// 跳过碰巧是有效 JSON 的情况
		if json.Valid([]byte(invalidJSON)) {
			return
		}

		// 创建一个不应被调用的服务器
		httpCalled := false
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpCalled = true
			fmt.Fprint(w, `{"code":0}`)
		}))
		defer server.Close()

		cfg := newTestConfigWithURL(server.URL)
		client := NewHttpClient(cfg)

		_, err := client.ExecuteRaw("test_method", invalidJSON)

		// 应返回错误
		if err == nil {
			t.Fatal("无效 JSON 应返回错误")
		}

		// 不应发送 HTTP 请求
		if httpCalled {
			t.Fatal("无效 JSON 不应发送 HTTP 请求")
		}
	})
}

// TestProperty15_EmptyApiMethodReject 测试空 API 方法名拒绝
func TestProperty15_EmptyApiMethodReject(t *testing.T) {
	httpCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpCalled = true
		fmt.Fprint(w, `{"code":0}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	_, err := client.ExecuteRaw("", `{"market":"US"}`)
	if err == nil {
		t.Fatal("空 api_method 应返回错误")
	}
	if httpCalled {
		t.Fatal("空 api_method 不应发送 HTTP 请求")
	}
}
