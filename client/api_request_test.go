package client

import (
	"encoding/json"
	"testing"
)

// TestNewApiRequest_WithString 测试使用字符串创建 API 请求
func TestNewApiRequest_WithString(t *testing.T) {
	req, err := NewApiRequest("market_state", `{"market":"US"}`)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	if req.Method != "market_state" {
		t.Errorf("期望 method=market_state，实际为 %s", req.Method)
	}
	if req.BizContent != `{"market":"US"}` {
		t.Errorf("期望 biz_content={\"market\":\"US\"}，实际为 %s", req.BizContent)
	}
}

// TestNewApiRequest_WithStruct 测试使用结构体创建 API 请求
func TestNewApiRequest_WithStruct(t *testing.T) {
	params := map[string]string{"market": "US", "symbol": "AAPL"}
	req, err := NewApiRequest("brief", params)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	if req.Method != "brief" {
		t.Errorf("期望 method=brief，实际为 %s", req.Method)
	}
	// 验证 biz_content 是有效 JSON
	if !json.Valid([]byte(req.BizContent)) {
		t.Errorf("biz_content 不是有效 JSON: %s", req.BizContent)
	}
	// 反序列化验证内容
	var parsed map[string]string
	json.Unmarshal([]byte(req.BizContent), &parsed)
	if parsed["market"] != "US" {
		t.Errorf("market 不匹配: %s", parsed["market"])
	}
	if parsed["symbol"] != "AAPL" {
		t.Errorf("symbol 不匹配: %s", parsed["symbol"])
	}
}

// TestNewApiRequest_WithNil 测试使用 nil 创建 API 请求
func TestNewApiRequest_WithNil(t *testing.T) {
	req, err := NewApiRequest("market_state", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	if req.BizContent != "{}" {
		t.Errorf("nil 参数应生成空 JSON 对象，实际为 %s", req.BizContent)
	}
}
