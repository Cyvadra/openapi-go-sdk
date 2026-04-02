package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHttpClient_Execute_Success 测试成功的 API 请求
func TestHttpClient_Execute_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和 Content-Type
		if r.Method != "POST" {
			t.Errorf("期望 POST 方法，实际为 %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/json;charset=UTF-8" {
			t.Errorf("期望 Content-Type 为 application/json;charset=UTF-8，实际为 %s", ct)
		}
		// 验证 User-Agent
		ua := r.Header.Get("User-Agent")
		expected := "openapi-go-sdk-0.1.0"
		if ua != expected {
			t.Errorf("期望 User-Agent 为 %s，实际为 %s", expected, ua)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":0,"message":"success","data":{"market":"US","status":"Trading"},"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	req := &ApiRequest{Method: "market_state", BizContent: `{"market":"US"}`}
	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	if resp.Code != 0 {
		t.Errorf("期望 code=0，实际为 %d", resp.Code)
	}
	if resp.Message != "success" {
		t.Errorf("期望 message=success，实际为 %s", resp.Message)
	}
	if resp.Timestamp != 1700000000 {
		t.Errorf("期望 timestamp=1700000000，实际为 %d", resp.Timestamp)
	}
}

// TestHttpClient_Execute_ApiError 测试 API 返回错误码
func TestHttpClient_Execute_ApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"code":1000,"message":"公共参数错误","data":null,"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	req := &ApiRequest{Method: "market_state", BizContent: `{}`}
	resp, err := client.Execute(req)
	if err == nil {
		t.Fatal("期望返回错误，但没有")
	}
	tigerErr, ok := err.(*TigerError)
	if !ok {
		t.Fatalf("期望 TigerError 类型，实际为 %T", err)
	}
	if tigerErr.Code != 1000 {
		t.Errorf("期望错误码 1000，实际为 %d", tigerErr.Code)
	}
	if tigerErr.Category != CategoryCommonParam {
		t.Errorf("期望分类 common_param_error，实际为 %s", tigerErr.Category)
	}
	if resp == nil {
		t.Fatal("即使有错误，resp 也不应为 nil")
	}
}

// TestHttpClient_Execute_RequestParams 测试请求参数包含必要字段
func TestHttpClient_Execute_RequestParams(t *testing.T) {
	var receivedParams map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params map[string]string
		json.NewDecoder(r.Body).Decode(&params)
		receivedParams = params
		fmt.Fprint(w, `{"code":0,"message":"success","data":null,"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	client := NewHttpClient(cfg)

	req := &ApiRequest{Method: "market_state", BizContent: `{"market":"US"}`}
	_, err := client.Execute(req)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}

	requiredKeys := []string{"tiger_id", "method", "charset", "sign_type", "timestamp", "version", "biz_content", "sign"}
	for _, key := range requiredKeys {
		if _, ok := receivedParams[key]; !ok {
			t.Errorf("缺少必要参数: %s", key)
		}
	}
	if receivedParams["tiger_id"] != "test_tiger_id" {
		t.Errorf("tiger_id 不匹配: %s", receivedParams["tiger_id"])
	}
	if receivedParams["method"] != "market_state" {
		t.Errorf("method 不匹配: %s", receivedParams["method"])
	}
	if receivedParams["biz_content"] != `{"market":"US"}` {
		t.Errorf("biz_content 不匹配: %s", receivedParams["biz_content"])
	}
}

// TestHttpClient_AuthorizationHeader_WithToken 测试设置 Token 时请求携带 Authorization 头
func TestHttpClient_AuthorizationHeader_WithToken(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		fmt.Fprint(w, `{"code":0,"message":"success","data":null,"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	cfg.Token = "my_test_token_value"
	client := NewHttpClient(cfg)

	req := &ApiRequest{Method: "market_state", BizContent: `{}`}
	_, err := client.Execute(req)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	if authHeader != "my_test_token_value" {
		t.Errorf("期望 Authorization 头为 my_test_token_value，实际为 %q", authHeader)
	}
}

// TestHttpClient_AuthorizationHeader_WithoutToken 测试未设置 Token 时不携带 Authorization 头
func TestHttpClient_AuthorizationHeader_WithoutToken(t *testing.T) {
	var hasAuth bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasAuth = r.Header["Authorization"]
		fmt.Fprint(w, `{"code":0,"message":"success","data":null,"timestamp":1700000000}`)
	}))
	defer server.Close()

	cfg := newTestConfig(t, server.URL)
	cfg.Token = "" // 无 token
	client := NewHttpClient(cfg)

	req := &ApiRequest{Method: "market_state", BizContent: `{}`}
	_, err := client.Execute(req)
	if err != nil {
		t.Fatalf("Execute 失败: %v", err)
	}
	if hasAuth {
		t.Error("未设置 Token 时不应携带 Authorization 头")
	}
}
