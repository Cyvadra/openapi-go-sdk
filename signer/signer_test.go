package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
)

// generateTestKeyPair 生成测试用 RSA 密钥对，返回 PKCS#1 格式私钥 PEM 和公钥
func generateTestKeyPair(t *testing.T) (string, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}
	// 编码为 PKCS#1 PEM 格式
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})
	return string(privPEM), &privateKey.PublicKey
}

// generateTestKeyPairPKCS8 生成测试用 RSA 密钥对，返回 PKCS#8 格式私钥 PEM 和公钥
func generateTestKeyPairPKCS8(t *testing.T) (string, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}
	// 编码为 PKCS#8 PEM 格式
	privDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("编码 PKCS#8 私钥失败: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privDER,
	})
	return string(privPEM), &privateKey.PublicKey
}

// generateRawBase64Key 生成裸 Base64 编码的私钥（无 PEM 头尾），用于测试兼容性
func generateRawBase64Key(t *testing.T) (string, *rsa.PublicKey) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	rawBase64 := base64.StdEncoding.EncodeToString(privDER)
	return rawBase64, &privateKey.PublicKey
}

// TestLoadPrivateKey_PKCS1 测试加载 PKCS#1 格式私钥
func TestLoadPrivateKey_PKCS1(t *testing.T) {
	privPEM, _ := generateTestKeyPair(t)
	key, err := LoadPrivateKey(privPEM)
	if err != nil {
		t.Fatalf("加载 PKCS#1 私钥失败: %v", err)
	}
	if key == nil {
		t.Fatal("加载的私钥不应为 nil")
	}
	// 验证密钥有效性
	if err := key.Validate(); err != nil {
		t.Fatalf("私钥验证失败: %v", err)
	}
}

// TestLoadPrivateKey_PKCS8 测试加载 PKCS#8 格式私钥
func TestLoadPrivateKey_PKCS8(t *testing.T) {
	privPEM, _ := generateTestKeyPairPKCS8(t)
	key, err := LoadPrivateKey(privPEM)
	if err != nil {
		t.Fatalf("加载 PKCS#8 私钥失败: %v", err)
	}
	if key == nil {
		t.Fatal("加载的私钥不应为 nil")
	}
	if err := key.Validate(); err != nil {
		t.Fatalf("私钥验证失败: %v", err)
	}
}

// TestLoadPrivateKey_RawBase64 测试加载裸 Base64 编码的私钥（无 PEM 头尾）
func TestLoadPrivateKey_RawBase64(t *testing.T) {
	rawBase64, _ := generateRawBase64Key(t)
	key, err := LoadPrivateKey(rawBase64)
	if err != nil {
		t.Fatalf("加载裸 Base64 私钥失败: %v", err)
	}
	if key == nil {
		t.Fatal("加载的私钥不应为 nil")
	}
	if err := key.Validate(); err != nil {
		t.Fatalf("私钥验证失败: %v", err)
	}
}

// TestLoadPrivateKey_Invalid 测试加载无效私钥应返回错误
func TestLoadPrivateKey_Invalid(t *testing.T) {
	_, err := LoadPrivateKey("invalid-key-data")
	if err == nil {
		t.Fatal("加载无效私钥应返回错误")
	}
}

// TestLoadPrivateKey_Empty 测试加载空私钥应返回错误
func TestLoadPrivateKey_Empty(t *testing.T) {
	_, err := LoadPrivateKey("")
	if err == nil {
		t.Fatal("加载空私钥应返回错误")
	}
}

// TestSignWithRSA_PKCS1 测试使用 PKCS#1 私钥签名并验签
func TestSignWithRSA_PKCS1(t *testing.T) {
	privPEM, pubKey := generateTestKeyPair(t)
	content := "tiger_id=test123&timestamp=1234567890"

	signature, err := SignWithRSA(privPEM, content)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}
	if signature == "" {
		t.Fatal("签名结果不应为空")
	}

	// 验证签名是有效的 Base64
	_, err = base64.StdEncoding.DecodeString(signature)
	if err != nil {
		t.Fatalf("签名结果不是有效的 Base64: %v", err)
	}

	// 使用公钥验签
	if err := VerifyWithRSA(pubKey, content, signature); err != nil {
		t.Fatalf("验签失败: %v", err)
	}
}

// TestSignWithRSA_PKCS8 测试使用 PKCS#8 私钥签名并验签
func TestSignWithRSA_PKCS8(t *testing.T) {
	privPEM, pubKey := generateTestKeyPairPKCS8(t)
	content := "biz_content={}&method=market_state"

	signature, err := SignWithRSA(privPEM, content)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}
	if signature == "" {
		t.Fatal("签名结果不应为空")
	}

	if err := VerifyWithRSA(pubKey, content, signature); err != nil {
		t.Fatalf("验签失败: %v", err)
	}
}

// TestSignWithRSA_DifferentContentDifferentSignature 测试不同内容产生不同签名
func TestSignWithRSA_DifferentContentDifferentSignature(t *testing.T) {
	privPEM, _ := generateTestKeyPair(t)

	sig1, err := SignWithRSA(privPEM, "content1")
	if err != nil {
		t.Fatalf("签名 content1 失败: %v", err)
	}

	sig2, err := SignWithRSA(privPEM, "content2")
	if err != nil {
		t.Fatalf("签名 content2 失败: %v", err)
	}

	if sig1 == sig2 {
		t.Fatal("不同内容的签名不应相同")
	}
}

// TestSignWithRSA_InvalidKey 测试使用无效私钥签名应返回错误
func TestSignWithRSA_InvalidKey(t *testing.T) {
	_, err := SignWithRSA("invalid-key", "test content")
	if err == nil {
		t.Fatal("使用无效私钥签名应返回错误")
	}
}
