package signer

import (
	"testing"
)

// TestGetSignContent_BasicOrder 测试基本的参数排序拼接
func TestGetSignContent_BasicOrder(t *testing.T) {
	params := map[string]string{
		"tiger_id":  "test123",
		"method":    "market_state",
		"charset":   "UTF-8",
		"sign_type": "RSA",
	}

	result := GetSignContent(params)
	expected := "charset=UTF-8&method=market_state&sign_type=RSA&tiger_id=test123"
	if result != expected {
		t.Errorf("期望 %q，实际 %q", expected, result)
	}
}

// TestGetSignContent_SingleParam 测试单个参数
func TestGetSignContent_SingleParam(t *testing.T) {
	params := map[string]string{
		"key": "value",
	}

	result := GetSignContent(params)
	expected := "key=value"
	if result != expected {
		t.Errorf("期望 %q，实际 %q", expected, result)
	}
}

// TestGetSignContent_EmptyMap 测试空参数映射
func TestGetSignContent_EmptyMap(t *testing.T) {
	params := map[string]string{}
	result := GetSignContent(params)
	if result != "" {
		t.Errorf("空参数映射应返回空字符串，实际 %q", result)
	}
}

// TestGetSignContent_AlphabeticalOrder 测试参数严格按字母序排列
func TestGetSignContent_AlphabeticalOrder(t *testing.T) {
	params := map[string]string{
		"zebra":    "z",
		"apple":    "a",
		"mango":    "m",
		"banana":   "b",
	}

	result := GetSignContent(params)
	expected := "apple=a&banana=b&mango=m&zebra=z"
	if result != expected {
		t.Errorf("期望 %q，实际 %q", expected, result)
	}
}

// TestGetSignContent_FullAPIParams 测试完整的 API 请求参数排序
func TestGetSignContent_FullAPIParams(t *testing.T) {
	params := map[string]string{
		"tiger_id":    "20150001",
		"method":      "market_state",
		"charset":     "UTF-8",
		"sign_type":   "RSA",
		"timestamp":   "2024-01-01 00:00:00",
		"version":     "3.0",
		"biz_content": `{"market":"US"}`,
	}

	result := GetSignContent(params)
	expected := `biz_content={"market":"US"}&charset=UTF-8&method=market_state&sign_type=RSA&tiger_id=20150001&timestamp=2024-01-01 00:00:00&version=3.0`
	if result != expected {
		t.Errorf("期望 %q，实际 %q", expected, result)
	}
}

// TestGetSignContent_SpecialCharactersInValues 测试值中包含特殊字符
func TestGetSignContent_SpecialCharactersInValues(t *testing.T) {
	params := map[string]string{
		"key1": "value with spaces",
		"key2": "value=with=equals",
		"key3": "value&with&ampersand",
	}

	result := GetSignContent(params)
	expected := "key1=value with spaces&key2=value=with=equals&key3=value&with&ampersand"
	if result != expected {
		t.Errorf("期望 %q，实际 %q", expected, result)
	}
}
