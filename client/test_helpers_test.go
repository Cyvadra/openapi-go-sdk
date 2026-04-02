package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/tigerfintech/openapi-go-sdk/config"
)

// mustGenerateKeyPEM 生成测试用 PKCS#1 PEM 格式私钥（不依赖 testing.T）
func mustGenerateKeyPEM() string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("生成 RSA 密钥对失败: " + err.Error())
	}
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})
	return string(privPEM)
}

// 缓存一个测试用私钥，避免每次测试都生成
var cachedTestPrivateKey = mustGenerateKeyPEM()

// newTestConfigWithURL 创建测试用 ClientConfig（不依赖 testing.T）
func newTestConfigWithURL(serverURL string) *config.ClientConfig {
	return &config.ClientConfig{
		TigerID:    "test_tiger_id",
		PrivateKey: cachedTestPrivateKey,
		Account:    "test_account",
		Language:   "zh_CN",
		Timeout:    5 * time.Second,
		ServerURL:  serverURL,
	}
}

// newTestConfig 创建测试用 ClientConfig（依赖 testing.T 版本）
func newTestConfig(t *testing.T, serverURL string) *config.ClientConfig {
	t.Helper()
	return newTestConfigWithURL(serverURL)
}
