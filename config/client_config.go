package config

import (
	"fmt"
	"os"
	"time"
)

const (
	// 默认值
	defaultLanguage  = "zh_CN"
	defaultTimeout   = 15 * time.Second
	defaultServerURL = "https://openapi.tigerfintech.com/gateway"
	sandboxServerURL = "https://openapi-sandbox.tigerfintech.com/gateway"

	// 环境变量名
	envTigerID    = "TIGEROPEN_TIGER_ID"
	envPrivateKey = "TIGEROPEN_PRIVATE_KEY"
	envAccount    = "TIGEROPEN_ACCOUNT"
)

// ClientConfig 客户端配置，包含认证信息和运行参数。
type ClientConfig struct {
	TigerID              string        `json:"tiger_id"`
	PrivateKey           string        `json:"private_key"`
	Account              string        `json:"account"`
	License              string        `json:"license"`
	Language             string        `json:"language"`
	Timezone             string        `json:"timezone"`
	Timeout              time.Duration `json:"-"`
	SandboxDebug         bool          `json:"-"`
	Token                string        `json:"-"`
	TokenRefreshDuration time.Duration `json:"-"`
	ServerURL            string        `json:"-"`
	EnableDynamicDomain  bool          `json:"-"`
}

// Option 配置选项函数类型
type Option func(*ClientConfig)

// WithTigerID 设置开发者 ID
func WithTigerID(id string) Option {
	return func(c *ClientConfig) { c.TigerID = id }
}

// WithPrivateKey 设置 RSA 私钥
func WithPrivateKey(key string) Option {
	return func(c *ClientConfig) { c.PrivateKey = key }
}

// WithAccount 设置交易账户
func WithAccount(account string) Option {
	return func(c *ClientConfig) { c.Account = account }
}

// WithLicense 设置牌照类型
func WithLicense(license string) Option {
	return func(c *ClientConfig) { c.License = license }
}

// WithLanguage 设置语言（zh_CN/zh_TW/en_US）
func WithLanguage(lang string) Option {
	return func(c *ClientConfig) { c.Language = lang }
}

// WithTimezone 设置时区
func WithTimezone(tz string) Option {
	return func(c *ClientConfig) { c.Timezone = tz }
}

// WithTimeout 设置请求超时时间
func WithTimeout(d time.Duration) Option {
	return func(c *ClientConfig) { c.Timeout = d }
}

// WithSandboxDebug 设置是否使用沙箱环境
func WithSandboxDebug(sandbox bool) Option {
	return func(c *ClientConfig) { c.SandboxDebug = sandbox }
}

// WithToken 设置 TBHK 牌照 Token
func WithToken(token string) Option {
	return func(c *ClientConfig) { c.Token = token }
}

// WithTokenRefreshDuration 设置 Token 刷新间隔
func WithTokenRefreshDuration(d time.Duration) Option {
	return func(c *ClientConfig) { c.TokenRefreshDuration = d }
}

// WithEnableDynamicDomain 设置是否启用动态域名获取（默认启用）
func WithEnableDynamicDomain(enable bool) Option {
	return func(c *ClientConfig) { c.EnableDynamicDomain = enable }
}

// WithPropertiesFile 从 properties 配置文件加载配置
func WithPropertiesFile(path string) Option {
	return func(c *ClientConfig) {
		props, err := ParsePropertiesFile(path)
		if err != nil {
			// 文件加载失败时静默跳过，后续校验会捕获必填字段缺失
			return
		}
		applyProperties(c, props)
	}
}

// applyProperties 将 properties 键值对应用到配置对象
func applyProperties(c *ClientConfig, props map[string]string) {
	if v, ok := props["tiger_id"]; ok && c.TigerID == "" {
		c.TigerID = v
	}
	// 私钥优先级：private_key > private_key_pk8 > private_key_pk1
	if c.PrivateKey == "" {
		if v, ok := props["private_key"]; ok {
			c.PrivateKey = v
		} else if v, ok := props["private_key_pk8"]; ok {
			c.PrivateKey = v
		} else if v, ok := props["private_key_pk1"]; ok {
			c.PrivateKey = v
		}
	}
	if v, ok := props["account"]; ok && c.Account == "" {
		c.Account = v
	}
	if v, ok := props["license"]; ok && c.License == "" {
		c.License = v
	}
	if v, ok := props["language"]; ok && c.Language == "" {
		c.Language = v
	}
	if v, ok := props["timezone"]; ok && c.Timezone == "" {
		c.Timezone = v
	}
}

// NewClientConfig 创建客户端配置。
// 优先级：环境变量 > Option 设置（含配置文件） > 默认值。
// 必填字段 tiger_id 和 private_key 为空时返回错误。
func NewClientConfig(opts ...Option) (*ClientConfig, error) {
	cfg := &ClientConfig{
		EnableDynamicDomain: true, // 默认启用动态域名
	}

	// 应用 Option（包括代码设置和配置文件加载）
	for _, opt := range opts {
		opt(cfg)
	}

	// 环境变量覆盖（最高优先级）
	if v := os.Getenv(envTigerID); v != "" {
		cfg.TigerID = v
	}
	if v := os.Getenv(envPrivateKey); v != "" {
		cfg.PrivateKey = v
	}
	if v := os.Getenv(envAccount); v != "" {
		cfg.Account = v
	}

	// 设置默认值
	if cfg.Language == "" {
		cfg.Language = defaultLanguage
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = defaultTimeout
	}

	// 确定服务器地址：sandbox > 动态域名 > 默认
	if cfg.SandboxDebug {
		cfg.ServerURL = sandboxServerURL
	} else if cfg.ServerURL == "" {
		// 尝试动态域名获取
		if cfg.EnableDynamicDomain {
			domainConf := QueryDomains(cfg.License)
			if dynamicURL := resolveDynamicServerURL(domainConf, cfg.License); dynamicURL != "" {
				cfg.ServerURL = dynamicURL
			}
		}
		// 动态域名获取失败或未启用，使用默认地址
		if cfg.ServerURL == "" {
			cfg.ServerURL = defaultServerURL
		}
	}

	// 校验必填字段
	if cfg.TigerID == "" {
		return nil, fmt.Errorf("tiger_id 不能为空，请通过 WithTigerID 或环境变量 %s 设置", envTigerID)
	}
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("private_key 不能为空，请通过 WithPrivateKey 或环境变量 %s 设置", envPrivateKey)
	}

	return cfg, nil
}
