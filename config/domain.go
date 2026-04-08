package config

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	// 域名花园地址
	domainGardenAddress = "https://cg.play-analytics.com"
	// 动态域名查询超时（快速失败）
	domainQueryTimeout = 1 * time.Second
	// TBUS 牌照标识
	licenseTBUS = "TBUS"
	// COMMON 域名 key
	domainKeyCommon = "COMMON"
	// QUOTE 域名后缀
	domainKeyQuoteSuffix = "-QUOTE"
	// gateway 后缀
	gatewaySuffix = "/gateway"
)

// domainResponse 域名查询响应
type domainResponse struct {
	Ret   int                          `json:"ret"`
	Items []map[string]json.RawMessage `json:"items"`
}

// QueryDomains 从域名花园获取动态域名配置。
// 失败时返回空 map（静默回退）。
func QueryDomains(license string) map[string]interface{} {
	url := domainGardenAddress
	if license == licenseTBUS {
		url += "?appName=tradeup"
	}

	client := &http.Client{Timeout: domainQueryTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result domainResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}

	if len(result.Items) == 0 {
		return nil
	}

	// 提取 openapi 配置
	item := result.Items[0]
	raw, ok := item["openapi"]
	if !ok {
		return nil
	}

	var conf map[string]interface{}
	if err := json.Unmarshal(raw, &conf); err != nil {
		return nil
	}

	return conf
}

// resolveDynamicServerURL 根据动态域名配置和 license 解析服务器地址。
// 返回空字符串表示无法解析（应使用默认地址）。
func resolveDynamicServerURL(domainConf map[string]interface{}, license string) string {
	if domainConf == nil || len(domainConf) == 0 {
		return ""
	}

	var key string
	if license != "" {
		key = license
	} else {
		key = domainKeyCommon
	}

	if url, ok := domainConf[key]; ok {
		if s, ok := url.(string); ok {
			return s + gatewaySuffix
		}
	}

	// 回退到 COMMON
	if url, ok := domainConf[domainKeyCommon]; ok {
		if s, ok := url.(string); ok {
			return s + gatewaySuffix
		}
	}

	return ""
}

// resolveDynamicQuoteServerURL 根据动态域名配置和 license 解析行情服务器地址。
func resolveDynamicQuoteServerURL(domainConf map[string]interface{}, license string) string {
	if domainConf == nil || len(domainConf) == 0 {
		return ""
	}

	if license != "" {
		if url, ok := domainConf[license+domainKeyQuoteSuffix]; ok {
			if s, ok := url.(string); ok {
				return s + gatewaySuffix
			}
		}
	}

	return resolveDynamicServerURL(domainConf, license)
}
