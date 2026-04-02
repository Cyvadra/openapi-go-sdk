package client

import (
	"math"
	"time"
)

// 交易操作方法名集合，这些操作不应自动重试
var tradeOperations = map[string]bool{
	"place_order":  true,
	"modify_order": true,
	"cancel_order": true,
}

// RetryPolicy 重试策略
type RetryPolicy struct {
	// MaxRetries 最大重试次数，默认 5
	MaxRetries int
	// MaxRetryTime 最大重试总时间，默认 60 秒
	MaxRetryTime time.Duration
	// BaseDelay 基础退避时间，默认 1 秒
	BaseDelay time.Duration
	// MaxDelay 最大单次退避时间，默认 16 秒
	MaxDelay time.Duration
}

// DefaultRetryPolicy 返回默认重试策略
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:   5,
		MaxRetryTime: 60 * time.Second,
		BaseDelay:    1 * time.Second,
		MaxDelay:     16 * time.Second,
	}
}

// ShouldRetry 判断指定的 API 方法是否应该重试
// 交易操作（place_order、modify_order、cancel_order）跳过重试
func (p *RetryPolicy) ShouldRetry(apiMethod string) bool {
	return !tradeOperations[apiMethod]
}

// IsTradeOperation 判断是否为交易操作
func IsTradeOperation(apiMethod string) bool {
	return tradeOperations[apiMethod]
}

// CalculateBackoff 计算第 n 次重试的退避等待时间（从 0 开始计数）
// 退避公式：min(2^n * baseDelay, maxDelay)
func (p *RetryPolicy) CalculateBackoff(retryCount int) time.Duration {
	if retryCount < 0 {
		return p.BaseDelay
	}
	delay := time.Duration(math.Pow(2, float64(retryCount))) * p.BaseDelay
	if delay > p.MaxDelay {
		delay = p.MaxDelay
	}
	return delay
}
