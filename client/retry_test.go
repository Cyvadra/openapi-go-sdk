package client

import (
	"testing"
	"time"
)

// TestRetryPolicy_DefaultValues 测试默认重试策略参数
func TestRetryPolicy_DefaultValues(t *testing.T) {
	p := DefaultRetryPolicy()
	if p.MaxRetries != 5 {
		t.Errorf("期望 MaxRetries=5，实际为 %d", p.MaxRetries)
	}
	if p.MaxRetryTime != 60*time.Second {
		t.Errorf("期望 MaxRetryTime=60s，实际为 %v", p.MaxRetryTime)
	}
	if p.BaseDelay != 1*time.Second {
		t.Errorf("期望 BaseDelay=1s，实际为 %v", p.BaseDelay)
	}
	if p.MaxDelay != 16*time.Second {
		t.Errorf("期望 MaxDelay=16s，实际为 %v", p.MaxDelay)
	}
}

// TestRetryPolicy_ExponentialBackoff 测试指数退避时间计算
func TestRetryPolicy_ExponentialBackoff(t *testing.T) {
	p := DefaultRetryPolicy()
	// 期望退避序列：1s → 2s → 4s → 8s → 16s
	expected := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
	}
	for i, exp := range expected {
		actual := p.CalculateBackoff(i)
		if actual != exp {
			t.Errorf("重试 %d: 期望退避 %v，实际为 %v", i, exp, actual)
		}
	}
}

// TestRetryPolicy_BackoffCap 测试退避时间上限
func TestRetryPolicy_BackoffCap(t *testing.T) {
	p := DefaultRetryPolicy()
	// 超过最大退避时间应被截断
	backoff := p.CalculateBackoff(10) // 2^10 = 1024s，应被截断为 16s
	if backoff != 16*time.Second {
		t.Errorf("期望退避被截断为 16s，实际为 %v", backoff)
	}
}

// TestRetryPolicy_ShouldRetry_TradeOperations 测试交易操作跳过重试
func TestRetryPolicy_ShouldRetry_TradeOperations(t *testing.T) {
	p := DefaultRetryPolicy()

	// 交易操作不应重试
	tradeOps := []string{"place_order", "modify_order", "cancel_order"}
	for _, op := range tradeOps {
		if p.ShouldRetry(op) {
			t.Errorf("交易操作 %s 不应重试", op)
		}
	}

	// 非交易操作应重试
	nonTradeOps := []string{"market_state", "brief", "kline", "positions", "assets"}
	for _, op := range nonTradeOps {
		if !p.ShouldRetry(op) {
			t.Errorf("非交易操作 %s 应允许重试", op)
		}
	}
}

// TestIsTradeOperation 测试交易操作判断
func TestIsTradeOperation(t *testing.T) {
	if !IsTradeOperation("place_order") {
		t.Error("place_order 应为交易操作")
	}
	if !IsTradeOperation("modify_order") {
		t.Error("modify_order 应为交易操作")
	}
	if !IsTradeOperation("cancel_order") {
		t.Error("cancel_order 应为交易操作")
	}
	if IsTradeOperation("market_state") {
		t.Error("market_state 不应为交易操作")
	}
}
