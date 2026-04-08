package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tigerfintech/openapi-go-sdk/config"
	"github.com/tigerfintech/openapi-go-sdk/signer"
)

const (
	// UserAgentPrefix User-Agent 前缀
	UserAgentPrefix = "openapi-go-sdk-"
	// SDKVersion SDK 版本号
	SDKVersion = "0.1.0"
	// DefaultCharset 默认字符集
	DefaultCharset = "UTF-8"
	// DefaultSignType 默认签名类型
	DefaultSignType = "RSA"
	// DefaultVersion 默认 API 版本
	DefaultVersion = "2.0"
)

// HttpClient 封装 HTTP 请求、签名、重试、超时
type HttpClient struct {
	config      *config.ClientConfig
	httpClient  *http.Client
	retryPolicy *RetryPolicy
}

// NewHttpClient 创建 HttpClient 实例
func NewHttpClient(cfg *config.ClientConfig) *HttpClient {
	return &HttpClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		retryPolicy: DefaultRetryPolicy(),
	}
}

// NewQuoteHttpClient 创建使用行情网关地址的 HttpClient。
func NewQuoteHttpClient(cfg *config.ClientConfig) *HttpClient {
	cloned := *cfg
	if cloned.QuoteServerURL != "" {
		cloned.ServerURL = cloned.QuoteServerURL
	}
	return NewHttpClient(&cloned)
}

// buildCommonParams 构造公共请求参数
func (c *HttpClient) buildCommonParams(apiMethod string, bizContent string, version string) map[string]string {
	if version == "" {
		version = DefaultVersion
	}
	params := map[string]string{
		"tiger_id":    c.config.TigerID,
		"method":      apiMethod,
		"charset":     DefaultCharset,
		"sign_type":   DefaultSignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     version,
		"biz_content": bizContent,
	}
	if c.config.Language != "" {
		params["language"] = c.config.Language
	}
	if c.config.DeviceID != "" {
		params["device_id"] = c.config.DeviceID
	}
	return params
}

// signParams 对参数进行签名
func (c *HttpClient) signParams(params map[string]string) (string, error) {
	content := signer.GetSignContent(params)
	sign, err := signer.SignWithRSA(c.config.PrivateKey, content)
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}
	return sign, nil
}

// Execute 执行结构化 API 请求，返回解析后的 ApiResponse
func (c *HttpClient) Execute(request *ApiRequest) (*ApiResponse, error) {
	params := c.buildCommonParams(request.Method, request.BizContent, request.Version)
	sign, err := c.signParams(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign

	var lastErr error
	maxAttempts := 1
	if c.retryPolicy.ShouldRetry(request.Method) {
		maxAttempts = c.retryPolicy.MaxRetries + 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := c.retryPolicy.CalculateBackoff(attempt - 1)
			time.Sleep(backoff)
		}

		resp, err := c.doHTTPPost(params)
		if err != nil {
			lastErr = err
			if !c.retryPolicy.ShouldRetry(request.Method) {
				return nil, err
			}
			continue
		}

		apiResp, parseErr := ParseApiResponse(resp)
		if parseErr != nil {
			return apiResp, parseErr
		}
		return apiResp, nil
	}

	return nil, lastErr
}

// ExecuteRaw 通用 API 调用方法
// apiMethod: API 方法名（如 "market_state"、"place_order"）
// requestJSON: 原始 biz_content JSON 字符串
// 返回原始 response JSON 字符串，不做任何解析
func (c *HttpClient) ExecuteRaw(apiMethod string, requestJSON string) (string, error) {
	// 参数校验
	if apiMethod == "" {
		return "", &TigerError{Code: -1, Message: "api_method 不能为空", Category: CategoryUnknown}
	}
	if !json.Valid([]byte(requestJSON)) {
		return "", &TigerError{Code: -1, Message: "request_json 不是有效的 JSON", Category: CategoryUnknown}
	}

	params := c.buildCommonParams(apiMethod, requestJSON, "")
	sign, err := c.signParams(params)
	if err != nil {
		return "", err
	}
	params["sign"] = sign

	var lastErr error
	maxAttempts := 1
	if c.retryPolicy.ShouldRetry(apiMethod) {
		maxAttempts = c.retryPolicy.MaxRetries + 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			backoff := c.retryPolicy.CalculateBackoff(attempt - 1)
			time.Sleep(backoff)
		}

		body, err := c.doHTTPPost(params)
		if err != nil {
			lastErr = err
			if !c.retryPolicy.ShouldRetry(apiMethod) {
				return "", err
			}
			continue
		}

		return string(body), nil
	}

	return "", lastErr
}

// doHTTPPost 发送 HTTP POST 请求
func (c *HttpClient) doHTTPPost(params map[string]string) ([]byte, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.ServerURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("创建 HTTP 请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", UserAgentPrefix+SDKVersion)
	if c.config.Token != "" {
		req.Header.Set("Authorization", c.config.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}
