package client

import (
	"encoding/json"
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 6: API 响应解析与错误处理
// **Validates: Requirements 3.5, 3.6**
// 对于任意有效的 JSON 响应字符串（包含 code、message、data 字段），
// 当 code 为 0 时解析应成功返回 data；当 code 不为 0 时应返回包含对应 code 和 message 的结构化错误。
func TestProperty6_ApiResponseParseAndErrorHandling(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机 code（0 表示成功，非 0 表示错误）
		code := rapid.IntRange(-100, 5000).Draw(t, "code")
		message := rapid.StringMatching(`[a-zA-Z0-9_ ]{1,50}`).Draw(t, "message")
		timestamp := rapid.Int64Range(1000000000, 2000000000).Draw(t, "timestamp")

		// 构造 JSON 响应
		respJSON := fmt.Sprintf(
			`{"code":%d,"message":"%s","data":{"key":"value"},"timestamp":%d}`,
			code, message, timestamp,
		)

		resp, err := ParseApiResponse([]byte(respJSON))

		if code == 0 {
			// code=0 时应成功
			if err != nil {
				t.Fatalf("code=0 时不应返回错误，但得到: %v", err)
			}
			if resp == nil {
				t.Fatal("code=0 时 resp 不应为 nil")
			}
			if resp.Code != 0 {
				t.Fatalf("期望 code=0，实际为 %d", resp.Code)
			}
			if resp.Message != message {
				t.Fatalf("message 不匹配: 期望 %s，实际 %s", message, resp.Message)
			}
			// data 应该存在
			if resp.Data == nil {
				t.Fatal("data 不应为 nil")
			}
		} else {
			// code!=0 时应返回 TigerError
			if err == nil {
				t.Fatalf("code=%d 时应返回错误", code)
			}
			tigerErr, ok := err.(*TigerError)
			if !ok {
				t.Fatalf("期望 *TigerError 类型，实际为 %T", err)
			}
			if tigerErr.Code != code {
				t.Fatalf("错误码不匹配: 期望 %d，实际 %d", code, tigerErr.Code)
			}
			if tigerErr.Message != message {
				t.Fatalf("错误消息不匹配: 期望 %s，实际 %s", message, tigerErr.Message)
			}
			// resp 仍然应该被返回（包含原始响应信息）
			if resp == nil {
				t.Fatal("即使有错误，resp 也不应为 nil")
			}
		}
	})
}

// TestProperty6_InvalidJSON 测试无效 JSON 的解析
func TestProperty6_InvalidJSON(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成不是有效 JSON 的字符串
		invalidJSON := rapid.StringMatching(`[a-zA-Z]{5,20}`).Draw(t, "invalidJSON")
		// 确保不是有效 JSON
		if json.Valid([]byte(invalidJSON)) {
			return // 跳过碰巧是有效 JSON 的情况
		}
		_, err := ParseApiResponse([]byte(invalidJSON))
		if err == nil {
			t.Fatal("无效 JSON 应返回错误")
		}
	})
}

// Feature: multi-language-sdks, Property 8: 错误码分类正确性
// **Validates: Requirements 8.2, 8.3, 8.4, 8.5**
// 对于任意已知的 API 错误码，SDK 返回的错误类别应与错误码对应的预定义分类一致。
func TestProperty8_ErrorCodeClassification(t *testing.T) {
	// 定义已知错误码与期望分类的映射
	knownCodes := map[int]ErrorCategory{
		0:    CategorySuccess,
		5:    CategoryRateLimit,
		1000: CategoryCommonParam,
		1001: CategoryCommonParam,
		1009: CategoryCommonParam,
		1010: CategoryBizParam,
		1050: CategoryBizParam,
		1099: CategoryBizParam,
		1100: CategoryTradeGlobal,
		1150: CategoryTradeGlobal,
		1199: CategoryTradeGlobal,
		1200: CategoryTradePrime,
		1250: CategoryTradePrime,
		1299: CategoryTradePrime,
		1300: CategoryTradeSimulation,
		2099: CategoryTradeSimulation,
		2100: CategoryQuoteStock,
		2150: CategoryQuoteStock,
		2199: CategoryQuoteStock,
		2200: CategoryQuoteOption,
		2250: CategoryQuoteOption,
		2299: CategoryQuoteOption,
		2300: CategoryQuoteFuture,
		2350: CategoryQuoteFuture,
		2399: CategoryQuoteFuture,
		2400: CategoryToken,
		3000: CategoryToken,
		3999: CategoryToken,
		4000: CategoryPermission,
		4500: CategoryPermission,
		4999: CategoryPermission,
	}

	for code, expectedCategory := range knownCodes {
		actual := ClassifyErrorCode(code)
		if actual != expectedCategory {
			t.Errorf("错误码 %d: 期望分类 %s，实际为 %s", code, expectedCategory, actual)
		}
	}
}

// TestProperty8_ErrorCodeClassificationProperty 属性测试：随机错误码分类一致性
func TestProperty8_ErrorCodeClassificationProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		code := rapid.IntRange(0, 10000).Draw(t, "code")
		category := ClassifyErrorCode(code)

		// 验证分类结果是有效的 ErrorCategory
		validCategories := map[ErrorCategory]bool{
			CategorySuccess:         true,
			CategoryCommonParam:     true,
			CategoryBizParam:        true,
			CategoryRateLimit:       true,
			CategoryTradeGlobal:     true,
			CategoryTradePrime:      true,
			CategoryTradeSimulation: true,
			CategoryQuoteStock:      true,
			CategoryQuoteOption:     true,
			CategoryQuoteFuture:     true,
			CategoryToken:           true,
			CategoryPermission:      true,
			CategoryServer:          true,
			CategoryUnknown:         true,
		}
		if !validCategories[category] {
			t.Fatalf("错误码 %d 返回了无效的分类: %s", code, category)
		}

		// 验证特定规则
		if code == 0 && category != CategorySuccess {
			t.Fatalf("code=0 应分类为 success，实际为 %s", category)
		}
		if code == 5 && category != CategoryRateLimit {
			t.Fatalf("code=5 应分类为 rate_limit，实际为 %s", category)
		}
		if code >= 1000 && code < 1010 && category != CategoryCommonParam {
			t.Fatalf("code=%d 应分类为 common_param_error，实际为 %s", code, category)
		}
		if code >= 1010 && code < 1100 && category != CategoryBizParam {
			t.Fatalf("code=%d 应分类为 biz_param_error，实际为 %s", code, category)
		}
	})
}
