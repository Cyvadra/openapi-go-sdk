package client

import (
	"encoding/json"
	"fmt"
)

// ApiRequest API 请求结构体
type ApiRequest struct {
	// Method API 方法名（如 "market_state"、"place_order"）
	Method string `json:"method"`
	// BizContent 业务参数 JSON 字符串
	BizContent string `json:"biz_content"`
}

// NewApiRequest 创建 API 请求，将业务参数序列化为 JSON 字符串作为 biz_content
func NewApiRequest(method string, bizParams interface{}) (*ApiRequest, error) {
	var bizContent string
	switch v := bizParams.(type) {
	case string:
		bizContent = v
	case nil:
		bizContent = "{}"
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("序列化 biz_content 失败: %w", err)
		}
		bizContent = string(data)
	}
	return &ApiRequest{
		Method:     method,
		BizContent: bizContent,
	}, nil
}
