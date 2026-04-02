package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 4: RSA 签名-验签 round-trip
// **Validates: Requirements 3.2**
//
// 对于任意非空字符串内容和有效的 RSA 密钥对，使用私钥对内容进行 SHA1WithRSA 签名后，
// 使用对应公钥验签应成功。
func TestProperty4_RSASignVerifyRoundTrip(t *testing.T) {
	// 预先生成密钥对（避免每次迭代都生成，提高测试速度）
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	// PKCS#1 格式 PEM
	privDER1 := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM1 := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER1,
	}))

	// PKCS#8 格式 PEM
	privDER8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("编码 PKCS#8 私钥失败: %v", err)
	}
	privPEM8 := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privDER8,
	}))

	pubKey := &privateKey.PublicKey

	rapid.Check(t, func(t *rapid.T) {
		// 生成任意非空字符串作为签名内容
		content := rapid.StringMatching(`.+`).Draw(t, "content")

		// 随机选择 PKCS#1 或 PKCS#8 格式
		usePKCS8 := rapid.Bool().Draw(t, "usePKCS8")
		var privPEM string
		if usePKCS8 {
			privPEM = privPEM8
		} else {
			privPEM = privPEM1
		}

		// 使用 SignWithRSA 签名
		signature, err := SignWithRSA(privPEM, content)
		if err != nil {
			t.Fatalf("签名失败: %v", err)
		}

		// 验证签名非空
		if signature == "" {
			t.Fatal("签名结果不应为空")
		}

		// 验证签名是有效的 Base64
		sigBytes, err := base64.StdEncoding.DecodeString(signature)
		if err != nil {
			t.Fatalf("签名结果不是有效的 Base64: %v", err)
		}

		// 使用公钥验签（SHA1WithRSA + PKCS1v15）
		hashed := sha1.Sum([]byte(content))
		err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA1, hashed[:], sigBytes)
		if err != nil {
			t.Fatalf("验签失败: %v", err)
		}
	})
}
