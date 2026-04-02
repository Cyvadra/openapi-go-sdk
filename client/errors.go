// Package client 提供 HTTP 客户端、重试策略和 API 请求构造功能。
package client

import "fmt"

// ErrorCategory 错误分类
type ErrorCategory string

const (
	// CategorySuccess 成功
	CategorySuccess ErrorCategory = "success"
	// CategoryCommonParam 公共参数错误
	CategoryCommonParam ErrorCategory = "common_param_error"
	// CategoryBizParam 业务参数错误
	CategoryBizParam ErrorCategory = "biz_param_error"
	// CategoryRateLimit 频率限制
	CategoryRateLimit ErrorCategory = "rate_limit"
	// CategoryTradeGlobal 环球账户交易错误
	CategoryTradeGlobal ErrorCategory = "trade_global_error"
	// CategoryTradePrime 综合账户交易错误
	CategoryTradePrime ErrorCategory = "trade_prime_error"
	// CategoryTradeSimulation 模拟账户交易错误
	CategoryTradeSimulation ErrorCategory = "trade_simulation_error"
	// CategoryQuoteStock 股票行情错误
	CategoryQuoteStock ErrorCategory = "quote_stock_error"
	// CategoryQuoteOption 期权行情错误
	CategoryQuoteOption ErrorCategory = "quote_option_error"
	// CategoryQuoteFuture 期货行情错误
	CategoryQuoteFuture ErrorCategory = "quote_future_error"
	// CategoryToken Token 错误
	CategoryToken ErrorCategory = "token_error"
	// CategoryPermission 权限错误
	CategoryPermission ErrorCategory = "permission_error"
	// CategoryServer 服务端错误
	CategoryServer ErrorCategory = "server_error"
	// CategoryUnknown 未知错误
	CategoryUnknown ErrorCategory = "unknown_error"
)

// ClassifyErrorCode 根据错误码返回对应的错误分类
func ClassifyErrorCode(code int) ErrorCategory {
	switch {
	case code == 0:
		return CategorySuccess
	case code == 5:
		return CategoryRateLimit
	case code >= 1000 && code < 1010:
		return CategoryCommonParam
	case code >= 1010 && code < 1100:
		return CategoryBizParam
	case code >= 1100 && code < 1200:
		return CategoryTradeGlobal
	case code >= 1200 && code < 1300:
		return CategoryTradePrime
	case code >= 1300 && code < 2100:
		return CategoryTradeSimulation
	case code >= 2100 && code < 2200:
		return CategoryQuoteStock
	case code >= 2200 && code < 2300:
		return CategoryQuoteOption
	case code >= 2300 && code < 2400:
		return CategoryQuoteFuture
	case code >= 2400 && code < 4000:
		return CategoryToken
	case code >= 4000 && code < 5000:
		return CategoryPermission
	default:
		return CategoryUnknown
	}
}

// TigerError 统一错误类型
type TigerError struct {
	Code     int
	Message  string
	Category ErrorCategory
}

func (e *TigerError) Error() string {
	return fmt.Sprintf("code=%d msg=%s category=%s", e.Code, e.Message, e.Category)
}

// NewTigerError 创建一个 TigerError
func NewTigerError(code int, message string) *TigerError {
	return &TigerError{
		Code:     code,
		Message:  message,
		Category: ClassifyErrorCode(code),
	}
}
