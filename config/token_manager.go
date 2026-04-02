package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// 默认 Token 文件名
	defaultTokenFileName = "tiger_openapi_token.properties"
	// 默认 Token 刷新间隔
	defaultTokenRefreshInterval = 24 * time.Hour
)

// TokenManager Token 管理器
// 从 properties 文件加载 Token，支持后台定期刷新。
type TokenManager struct {
	mu       sync.RWMutex
	token    string
	filePath string
	interval time.Duration
	stopCh   chan struct{}
	// Token 刷新阈值（秒），0 表示不刷新
	refreshDuration int64
	// 可注入的刷新函数（用于测试）
	refreshFn func() (string, error)
}

// TokenManagerOption TokenManager 配置选项
type TokenManagerOption func(*TokenManager)

// WithTokenFilePath 设置 Token 文件路径
func WithTokenFilePath(path string) TokenManagerOption {
	return func(m *TokenManager) { m.filePath = path }
}

// WithTokenRefreshInterval 设置 Token 检查间隔
func WithTokenRefreshInterval(d time.Duration) TokenManagerOption {
	return func(m *TokenManager) { m.interval = d }
}

// WithRefreshDuration 设置 Token 刷新阈值（秒），当 token 生成时间超过此值时触发刷新。
// 0 表示不刷新，最小值 30 秒。
func WithRefreshDuration(seconds int64) TokenManagerOption {
	return func(m *TokenManager) {
		if seconds > 0 && seconds < 30 {
			seconds = 30
		}
		m.refreshDuration = seconds
	}
}

// NewTokenManager 创建 Token 管理器
func NewTokenManager(opts ...TokenManagerOption) *TokenManager {
	m := &TokenManager{
		filePath: defaultTokenFileName,
		interval: defaultTokenRefreshInterval,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// LoadToken 从 properties 文件加载 Token
func (m *TokenManager) LoadToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	props, err := ParsePropertiesFile(m.filePath)
	if err != nil {
		return "", fmt.Errorf("加载 Token 文件失败: %w", err)
	}

	token, ok := props["token"]
	if !ok || token == "" {
		return "", fmt.Errorf("Token 文件中未找到 token 字段")
	}

	m.token = token
	return token, nil
}

// GetToken 获取当前 Token
func (m *TokenManager) GetToken() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.token
}

// SetToken 设置 Token 并更新文件
func (m *TokenManager) SetToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.token = token
	return m.saveTokenToFile(token)
}

// saveTokenToFile 将 Token 保存到 properties 文件
func (m *TokenManager) saveTokenToFile(token string) error {
	dir := filepath.Dir(m.filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}
	content := fmt.Sprintf("token=%s\n", token)
	return os.WriteFile(m.filePath, []byte(content), 0644)
}

// ShouldTokenRefresh 判断 Token 是否需要刷新。
// 解码 base64 token，提取前 27 字符中的 gen_ts，
// 当 (当前时间秒 - gen_ts/1000) > refreshDuration 时返回 true。
func (m *TokenManager) ShouldTokenRefresh() bool {
	m.mu.RLock()
	token := m.token
	dur := m.refreshDuration
	m.mu.RUnlock()

	if token == "" || dur == 0 {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	if len(decoded) < 27 {
		return false
	}

	parts := strings.SplitN(string(decoded[:27]), ",", 2)
	if len(parts) < 2 {
		return false
	}

	genTs, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return false
	}

	return (time.Now().Unix() - genTs/1000) > dur
}

// StartAutoRefresh 启动后台定期刷新
func (m *TokenManager) StartAutoRefresh(refreshFn func() (string, error)) {
	m.mu.Lock()
	m.refreshFn = refreshFn
	m.stopCh = make(chan struct{})
	m.mu.Unlock()

	go m.refreshLoop()
}

// StopAutoRefresh 停止后台刷新
func (m *TokenManager) StopAutoRefresh() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopCh != nil {
		close(m.stopCh)
		m.stopCh = nil
	}
}

// refreshLoop 后台刷新循环
func (m *TokenManager) refreshLoop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ticker.C:
			// 先检查是否需要刷新
			if !m.ShouldTokenRefresh() {
				continue
			}
			m.mu.RLock()
			fn := m.refreshFn
			m.mu.RUnlock()
			if fn != nil {
				if newToken, err := fn(); err == nil && newToken != "" {
					m.SetToken(newToken)
				}
			}
		}
	}
}
