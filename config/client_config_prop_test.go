package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pgregory.net/rapid"
)

func clearTigerEnv(t testing.TB) {
	t.Helper()
	for _, key := range []string{
		"TIGEROPEN_TIGER_ID",
		"TIGEROPEN_PRIVATE_KEY",
		"TIGEROPEN_ACCOUNT",
	} {
		key := key
		previous, ok := os.LookupEnv(key)
		os.Unsetenv(key)
		t.Cleanup(func() {
			if ok {
				os.Setenv(key, previous)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

// Feature: multi-language-sdks, Property 2: ClientConfig 字段设置 round-trip
// **Validates: Requirements 2.1, 2.6**
//
// 对于任意有效的配置参数组合，通过代码设置到 ClientConfig 后，
// 读取各字段的值应与设置的值完全一致。
func TestClientConfigFieldsRoundTrip(t *testing.T) {
	clearTigerEnv(t)
	rapid.Check(t, func(t *rapid.T) {
		tigerID := rapid.StringMatching(`[a-zA-Z0-9]{1,20}`).Draw(t, "tigerID")
		privateKey := rapid.StringMatching(`[a-zA-Z0-9]{1,100}`).Draw(t, "privateKey")
		account := rapid.StringMatching(`[A-Z]{2}[0-9]{4,10}`).Draw(t, "account")
		language := rapid.SampledFrom([]string{"zh_CN", "zh_TW", "en_US"}).Draw(t, "language")
		timezone := rapid.StringMatching(`[A-Za-z/_]{3,30}`).Draw(t, "timezone")
		timeoutSec := rapid.IntRange(1, 120).Draw(t, "timeoutSec")
		sandbox := rapid.Bool().Draw(t, "sandbox")

		cfg, err := NewClientConfig(
			WithTigerID(tigerID),
			WithPrivateKey(privateKey),
			WithAccount(account),
			WithLanguage(language),
			WithTimezone(timezone),
			WithTimeout(time.Duration(timeoutSec)*time.Second),
			WithSandboxDebug(sandbox),
		)
		if err != nil {
			t.Fatalf("创建配置失败: %v", err)
		}

		// 验证 round-trip：读取值应与设置值一致
		if cfg.TigerID != tigerID {
			t.Fatalf("TigerID: 期望 %q，实际 %q", tigerID, cfg.TigerID)
		}
		if cfg.PrivateKey != privateKey {
			t.Fatalf("PrivateKey: 期望 %q，实际 %q", privateKey, cfg.PrivateKey)
		}
		if cfg.Account != account {
			t.Fatalf("Account: 期望 %q，实际 %q", account, cfg.Account)
		}
		if cfg.Language != language {
			t.Fatalf("Language: 期望 %q，实际 %q", language, cfg.Language)
		}
		if cfg.Timezone != timezone {
			t.Fatalf("Timezone: 期望 %q，实际 %q", timezone, cfg.Timezone)
		}
		if cfg.Timeout != time.Duration(timeoutSec)*time.Second {
			t.Fatalf("Timeout: 期望 %v，实际 %v", time.Duration(timeoutSec)*time.Second, cfg.Timeout)
		}
		if cfg.SandboxDebug != sandbox {
			t.Fatalf("SandboxDebug: 期望 %v，实际 %v", sandbox, cfg.SandboxDebug)
		}
	})
}

// Feature: multi-language-sdks, Property 3: 环境变量优先级高于配置文件
// **Validates: Requirements 2.4**
//
// 对于任意配置字段（tiger_id、private_key、account），当环境变量和配置文件同时提供该字段的值时，
// ClientConfig 最终使用的值应等于环境变量中的值。
func TestEnvOverridesConfigFile(t *testing.T) {
	clearTigerEnv(t)
	rapid.Check(t, func(t *rapid.T) {
		// 生成配置文件中的值
		fileTigerID := rapid.StringMatching(`file_[a-zA-Z0-9]{1,15}`).Draw(t, "fileTigerID")
		filePrivateKey := rapid.StringMatching(`file_[a-zA-Z0-9]{1,50}`).Draw(t, "filePrivateKey")
		fileAccount := rapid.StringMatching(`file_[A-Z]{2}[0-9]{4,8}`).Draw(t, "fileAccount")

		// 生成环境变量中的值
		envTigerID := rapid.StringMatching(`env_[a-zA-Z0-9]{1,15}`).Draw(t, "envTigerID")
		envPrivateKey := rapid.StringMatching(`env_[a-zA-Z0-9]{1,50}`).Draw(t, "envPrivateKey")
		envAccount := rapid.StringMatching(`env_[A-Z]{2}[0-9]{4,8}`).Draw(t, "envAccount")

		// 写入配置文件
		content := fmt.Sprintf("tiger_id=%s\nprivate_key=%s\naccount=%s\n",
			fileTigerID, filePrivateKey, fileAccount)
		dir := os.TempDir()
		path := filepath.Join(dir, fmt.Sprintf("env_test_%d.properties",
			rapid.IntRange(0, 999999).Draw(t, "fileId")))
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		defer os.Remove(path)

		// 设置环境变量
		os.Setenv("TIGEROPEN_TIGER_ID", envTigerID)
		os.Setenv("TIGEROPEN_PRIVATE_KEY", envPrivateKey)
		os.Setenv("TIGEROPEN_ACCOUNT", envAccount)
		defer func() {
			os.Unsetenv("TIGEROPEN_TIGER_ID")
			os.Unsetenv("TIGEROPEN_PRIVATE_KEY")
			os.Unsetenv("TIGEROPEN_ACCOUNT")
		}()

		cfg, err := NewClientConfig(
			WithPropertiesFile(path),
		)
		if err != nil {
			t.Fatalf("创建配置失败: %v", err)
		}

		// 环境变量应覆盖配置文件
		if cfg.TigerID != envTigerID {
			t.Fatalf("TigerID: 期望环境变量值 %q，实际 %q", envTigerID, cfg.TigerID)
		}
		if cfg.PrivateKey != envPrivateKey {
			t.Fatalf("PrivateKey: 期望环境变量值 %q，实际 %q", envPrivateKey, cfg.PrivateKey)
		}
		if cfg.Account != envAccount {
			t.Fatalf("Account: 期望环境变量值 %q，实际 %q", envAccount, cfg.Account)
		}
	})
}
