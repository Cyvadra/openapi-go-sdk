package client

import (
	"math"
	"testing"
	"time"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 9: 指数退避时间计算
// **Validates: Requirements 11.3**
// 对于任意重试次数 n（0 ≤ n < 最大重试次数），计算的退避等待时间应等于 min(2^n, 最大退避时间) 秒。
func TestProperty9_ExponentialBackoffCalculation(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		p := DefaultRetryPolicy()
		retryCount := rapid.IntRange(0, p.MaxRetries-1).Draw(t, "retryCount")

		actual := p.CalculateBackoff(retryCount)

		// 期望值：min(2^n * baseDelay, maxDelay)
		expectedSeconds := math.Pow(2, float64(retryCount))
		expectedDuration := time.Duration(expectedSeconds) * p.BaseDelay
		if expectedDuration > p.MaxDelay {
			expectedDuration = p.MaxDelay
		}

		if actual != expectedDuration {
			t.Fatalf("重试 %d: 期望退避 %v，实际为 %v", retryCount, expectedDuration, actual)
		}

		// 验证退避时间不超过最大退避时间
		if actual > p.MaxDelay {
			t.Fatalf("退避时间 %v 超过最大退避时间 %v", actual, p.MaxDelay)
		}

		// 验证退避时间不小于基础退避时间
		if actual < p.BaseDelay {
			t.Fatalf("退避时间 %v 小于基础退避时间 %v", actual, p.BaseDelay)
		}
	})
}

// Feature: multi-language-sdks, Property 10: 交易操作跳过重试
// **Validates: Requirements 11.4**
// 对于任意服务类型标识，当且仅当该标识属于 {place_order, modify_order, cancel_order} 时，
// 重试策略应返回"不重试"。
func TestProperty10_TradeOperationsSkipRetry(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		p := DefaultRetryPolicy()

		// 从包含交易和非交易操作的集合中随机选择
		allMethods := []string{
			"place_order", "modify_order", "cancel_order",
			"market_state", "brief", "kline", "timeline",
			"trade_tick", "quote_depth", "positions", "assets",
			"orders", "active_orders", "contract", "contracts",
			"option_expiration", "option_chain", "future_exchange",
			"financial_daily", "capital_flow", "market_scanner",
		}
		idx := rapid.IntRange(0, len(allMethods)-1).Draw(t, "methodIdx")
		method := allMethods[idx]

		shouldRetry := p.ShouldRetry(method)
		isTradeOp := IsTradeOperation(method)

		// 交易操作不应重试，非交易操作应重试
		if isTradeOp && shouldRetry {
			t.Fatalf("交易操作 %s 不应重试", method)
		}
		if !isTradeOp && !shouldRetry {
			t.Fatalf("非交易操作 %s 应允许重试", method)
		}

		// 验证 IsTradeOperation 与 ShouldRetry 互斥
		if isTradeOp == shouldRetry {
			t.Fatalf("方法 %s: IsTradeOperation=%v 与 ShouldRetry=%v 应互斥",
				method, isTradeOp, shouldRetry)
		}
	})
}
