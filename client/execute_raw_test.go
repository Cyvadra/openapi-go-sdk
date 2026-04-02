package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExecuteRaw_Success 测试 ExecuteRaw 成功调用
func TestExecuteRaw_Success(t *testing.T) {
	expectedResp := `{"code":0,"message":"success","data":{"market":"US","status":"Trading"},"timestamp":1700000000}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedResp)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	result, err := client.ExecuteRaw("market_state", `{"market":"US"}`)
	if err != nil {
		t.Fatalf("ExecuteRaw 失败: %v", err)
	}
	if result != expectedResp {
		t.Errorf("响应不匹配:\n期望: %s\n实际: %s", expectedResp, result)
	}
}

// TestExecuteRaw_EmptyApiMethod 测试空 API 方法名
func TestExecuteRaw_EmptyApiMethod(t *testing.T) {
	cfg := newTestConfig(t, "http://localhost:9999")
	client := NewHttpClient(cfg)

	_, err := client.ExecuteRaw("", `{"market":"US"}`)
	if err == nil {
		t.Fatal("空 api_method 应返回错误")
	}
}

// TestExecuteRaw_InvalidJSON 测试无效 JSON 参数
func TestExecuteRaw_InvalidJSON(t *testing.T) {
	cfg := newTestConfig(t, "http://localhost:9999")
	client := NewHttpClient(cfg)

	_, err := client.ExecuteRaw("market_state", "not-valid-json")
	if err == nil {
		t.Fatal("无效 JSON 应返回错误")
	}
}

// TestExecuteRaw_RequestConstruction 测试 ExecuteRaw 请求参数构造
func TestExecuteRaw_RequestConstruction(t *testing.T) {
	var receivedParams map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedParams)
		fmt.Fprint(w, `{"code":0,"message":"success","data":null,"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	bizJSON := `{"symbol":"AAPL","market":"US"}`
	_, err := client.ExecuteRaw("brief", bizJSON)
	if err != nil {
		t.Fatalf("ExecuteRaw 失败: %v", err)
	}

	// 验证公共参数
	if receivedParams["method"] != "brief" {
		t.Errorf("method 不匹配: %s", receivedParams["method"])
	}
	if receivedParams["biz_content"] != bizJSON {
		t.Errorf("biz_content 不匹配: %s", receivedParams["biz_content"])
	}
	if receivedParams["tiger_id"] != "test_tiger_id" {
		t.Errorf("tiger_id 不匹配: %s", receivedParams["tiger_id"])
	}
	if _, ok := receivedParams["sign"]; !ok {
		t.Error("缺少 sign 参数")
	}
	if _, ok := receivedParams["timestamp"]; !ok {
		t.Error("缺少 timestamp 参数")
	}
}

// TestExecuteRaw_RawResponsePassthrough 测试原始响应透传
func TestExecuteRaw_RawResponsePassthrough(t *testing.T) {
	// 服务器返回各种格式的 JSON
	testCases := []string{
		`{"code":0,"message":"ok","data":[1,2,3]}`,
		`{"code":0,"message":"ok","data":{"nested":{"deep":"value"}}}`,
		`{"code":0,"message":"ok","data":null}`,
		`{"code":1000,"message":"error","data":null}`,
	}

	for _, expected := range testCases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))

		cfg := newTestConfig(t, server.URL)
		client := NewHttpClient(cfg)

		result, _ := client.ExecuteRaw("test_method", `{}`)
		if result != expected {
			t.Errorf("响应不匹配:\n期望: %s\n实际: %s", expected, result)
		}
		server.Close()
	}
}
