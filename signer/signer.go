// Package signer 提供 RSA 签名和请求参数排序拼接功能。
// 用于 OpenAPI 请求的认证签名流程。
package signer

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// LoadPrivateKey 加载 RSA 私钥，支持以下格式：
// - PKCS#1 PEM（BEGIN RSA PRIVATE KEY）
// - PKCS#8 PEM（BEGIN PRIVATE KEY）
// - 裸 Base64 编码的 DER 数据（无 PEM 头尾）
func LoadPrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	if keyStr == "" {
		return nil, fmt.Errorf("私钥不能为空")
	}

	// 尝试解析 PEM 格式
	block, _ := pem.Decode([]byte(keyStr))
	if block != nil {
		return parseDERPrivateKey(block.Bytes)
	}

	// 非 PEM 格式，尝试作为裸 Base64 解码
	derBytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析私钥：不是有效的 PEM 或 Base64 格式")
	}

	return parseDERPrivateKey(derBytes)
}

// parseDERPrivateKey 从 DER 字节解析 RSA 私钥，先尝试 PKCS#1，再尝试 PKCS#8
func parseDERPrivateKey(derBytes []byte) (*rsa.PrivateKey, error) {
	// 先尝试 PKCS#1
	if key, err := x509.ParsePKCS1PrivateKey(derBytes); err == nil {
		return key, nil
	}

	// 再尝试 PKCS#8
	keyIface, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析私钥：不是有效的 PKCS#1 或 PKCS#8 格式")
	}

	rsaKey, ok := keyIface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥不是 RSA 类型")
	}

	return rsaKey, nil
}

// SignWithRSA 使用 RSA 私钥对内容进行 SHA1WithRSA 签名，返回 Base64 编码的签名字符串。
func SignWithRSA(privateKeyStr string, content string) (string, error) {
	key, err := LoadPrivateKey(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("加载私钥失败: %w", err)
	}

	hashed := sha1.Sum([]byte(content))
	signature, err := rsa.SignPKCS1v15(nil, key, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyWithRSA 使用 RSA 公钥验证 SHA1WithRSA 签名。
// signature 为 Base64 编码的签名字符串。
func VerifyWithRSA(publicKey *rsa.PublicKey, content string, signature string) error {
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("签名 Base64 解码失败: %w", err)
	}

	hashed := sha1.Sum([]byte(content))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashed[:], sigBytes)
}
