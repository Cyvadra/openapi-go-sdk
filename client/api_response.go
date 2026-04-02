package client

import (
	"encoding/json"
	"fmt"
)

// ApiResponse API 响应结构体
type ApiResponse struct {
	// Code 状态码（0=成功）
	Code int `json:"code"`
	// Message 状态描述
	Message string `json:"message"`
	// Data 业务数据（原始 JSON）
	Data json.RawMessage `json:"data"`
	// Timestamp 服务器时间戳
	Timestamp int64 `json:"timestamp"`
}

// ParseApiResponse 解析 API 响应 JSON 字符串
// 当 code 为 0 时返回 ApiResponse；当 code 不为 0 时返回 TigerError
func ParseApiResponse(body []byte) (*ApiResponse, error) {
	var resp ApiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}
	if resp.Code != 0 {
		return &resp, NewTigerError(resp.Code, resp.Message)
	}
	return &resp, nil
}
