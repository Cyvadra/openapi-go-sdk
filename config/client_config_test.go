package config

import (
	"os"
	"testing"
	"time"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"TIGEROPEN_TIGER_ID",
		"TIGEROPEN_PRIVATE_KEY",
		"TIGEROPEN_ACCOUNT",
		"TIGEROPEN_TOKEN",
		"TIGEROPEN_TOKEN_FILE",
	} {
		t.Setenv(key, "")
	}
}

// TestNewClientConfig_WithOptions 测试通过 Option 模式设置字段
func TestNewClientConfig_WithOptions(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := NewClientConfig(
		WithTigerID("test_tiger_id"),
		WithPrivateKey("test_private_key"),
		WithAccount("DU123456"),
		WithLanguage("en_US"),
		WithTimezone("America/New_York"),
		WithTimeout(30*time.Second),
		WithSandboxDebug(true),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.TigerID != "test_tiger_id" {
		t.Fatalf("TigerID: 期望 %q，实际 %q", "test_tiger_id", cfg.TigerID)
	}
	if cfg.PrivateKey != "test_private_key" {
		t.Fatalf("PrivateKey: 期望 %q，实际 %q", "test_private_key", cfg.PrivateKey)
	}
	if cfg.Account != "DU123456" {
		t.Fatalf("Account: 期望 %q，实际 %q", "DU123456", cfg.Account)
	}
	if cfg.Language != "en_US" {
		t.Fatalf("Language: 期望 %q，实际 %q", "en_US", cfg.Language)
	}
	if cfg.Timezone != "America/New_York" {
		t.Fatalf("Timezone: 期望 %q，实际 %q", "America/New_York", cfg.Timezone)
	}
	if cfg.Timeout != 30*time.Second {
		t.Fatalf("Timeout: 期望 %v，实际 %v", 30*time.Second, cfg.Timeout)
	}
	if !cfg.SandboxDebug {
		t.Fatal("SandboxDebug: 期望 true，实际 false")
	}
}

// TestNewClientConfig_Defaults 测试默认值
func TestNewClientConfig_Defaults(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := NewClientConfig(
		WithTigerID("tid"),
		WithPrivateKey("pk"),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.Language != "zh_CN" {
		t.Fatalf("默认 Language: 期望 %q，实际 %q", "zh_CN", cfg.Language)
	}
	if cfg.Timeout != 15*time.Second {
		t.Fatalf("默认 Timeout: 期望 %v，实际 %v", 15*time.Second, cfg.Timeout)
	}
	if cfg.SandboxDebug {
		t.Fatal("默认 SandboxDebug: 期望 false，实际 true")
	}
	if cfg.ServerURL != "https://openapi.tigerfintech.com/gateway" {
		t.Fatalf("默认 ServerURL: 期望生产环境 URL，实际 %q", cfg.ServerURL)
	}
}

// TestNewClientConfig_SandboxURL 测试沙箱模式 URL
func TestNewClientConfig_SandboxURL(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := NewClientConfig(
		WithTigerID("tid"),
		WithPrivateKey("pk"),
		WithSandboxDebug(true),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.ServerURL != "https://openapi-sandbox.tigerfintech.com/gateway" {
		t.Fatalf("沙箱 ServerURL: 期望沙箱 URL，实际 %q", cfg.ServerURL)
	}
}

// TestNewClientConfig_MissingTigerID 测试缺少 tiger_id 时返回错误
func TestNewClientConfig_MissingTigerID(t *testing.T) {
	clearConfigEnv(t)
	_, err := NewClientConfig(
		WithPrivateKey("pk"),
	)
	if err == nil {
		t.Fatal("期望返回错误（缺少 tiger_id），但没有")
	}
}

// TestNewClientConfig_MissingPrivateKey 测试缺少 private_key 时返回错误
func TestNewClientConfig_MissingPrivateKey(t *testing.T) {
	clearConfigEnv(t)
	_, err := NewClientConfig(
		WithTigerID("tid"),
	)
	if err == nil {
		t.Fatal("期望返回错误（缺少 private_key），但没有")
	}
}

// TestNewClientConfig_EnvOverride 测试环境变量覆盖代码设置
func TestNewClientConfig_EnvOverride(t *testing.T) {
	clearConfigEnv(t)
	// 设置环境变量
	os.Setenv("TIGEROPEN_TIGER_ID", "env_tiger_id")
	os.Setenv("TIGEROPEN_PRIVATE_KEY", "env_private_key")
	os.Setenv("TIGEROPEN_ACCOUNT", "env_account")
	defer func() {
		os.Unsetenv("TIGEROPEN_TIGER_ID")
		os.Unsetenv("TIGEROPEN_PRIVATE_KEY")
		os.Unsetenv("TIGEROPEN_ACCOUNT")
	}()

	cfg, err := NewClientConfig(
		WithTigerID("code_tiger_id"),
		WithPrivateKey("code_private_key"),
		WithAccount("code_account"),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	// 环境变量应覆盖代码设置
	if cfg.TigerID != "env_tiger_id" {
		t.Fatalf("TigerID: 期望环境变量值 %q，实际 %q", "env_tiger_id", cfg.TigerID)
	}
	if cfg.PrivateKey != "env_private_key" {
		t.Fatalf("PrivateKey: 期望环境变量值 %q，实际 %q", "env_private_key", cfg.PrivateKey)
	}
	if cfg.Account != "env_account" {
		t.Fatalf("Account: 期望环境变量值 %q，实际 %q", "env_account", cfg.Account)
	}
}

// TestNewClientConfig_EnvOverridePartial 测试环境变量部分覆盖
func TestNewClientConfig_EnvOverridePartial(t *testing.T) {
	clearConfigEnv(t)
	os.Setenv("TIGEROPEN_TIGER_ID", "env_tid")
	defer os.Unsetenv("TIGEROPEN_TIGER_ID")

	cfg, err := NewClientConfig(
		WithTigerID("code_tid"),
		WithPrivateKey("code_pk"),
		WithAccount("code_acc"),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	// 只有 tiger_id 被环境变量覆盖
	if cfg.TigerID != "env_tid" {
		t.Fatalf("TigerID: 期望 %q，实际 %q", "env_tid", cfg.TigerID)
	}
	if cfg.PrivateKey != "code_pk" {
		t.Fatalf("PrivateKey: 期望 %q，实际 %q", "code_pk", cfg.PrivateKey)
	}
	if cfg.Account != "code_acc" {
		t.Fatalf("Account: 期望 %q，实际 %q", "code_acc", cfg.Account)
	}
}

// TestNewClientConfig_FromPropertiesFile 测试从配置文件加载
func TestNewClientConfig_FromPropertiesFile(t *testing.T) {
	clearConfigEnv(t)
	content := "tiger_id=file_tid\nprivate_key=file_pk\naccount=file_acc\nlicense=TBNZ\n"
	path := writeTempFile(t, content)

	cfg, err := NewClientConfig(
		WithPropertiesFile(path),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.TigerID != "file_tid" {
		t.Fatalf("TigerID: 期望 %q，实际 %q", "file_tid", cfg.TigerID)
	}
	if cfg.PrivateKey != "file_pk" {
		t.Fatalf("PrivateKey: 期望 %q，实际 %q", "file_pk", cfg.PrivateKey)
	}
	if cfg.Account != "file_acc" {
		t.Fatalf("Account: 期望 %q，实际 %q", "file_acc", cfg.Account)
	}
	if cfg.License != "TBNZ" {
		t.Fatalf("License: 期望 %q，实际 %q", "TBNZ", cfg.License)
	}
}

// TestNewClientConfig_EnvOverridesFile 测试环境变量优先级高于配置文件
func TestNewClientConfig_EnvOverridesFile(t *testing.T) {
	clearConfigEnv(t)
	content := "tiger_id=file_tid\nprivate_key=file_pk\naccount=file_acc\n"
	path := writeTempFile(t, content)

	os.Setenv("TIGEROPEN_TIGER_ID", "env_tid")
	defer os.Unsetenv("TIGEROPEN_TIGER_ID")

	cfg, err := NewClientConfig(
		WithPropertiesFile(path),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	// 环境变量覆盖配置文件
	if cfg.TigerID != "env_tid" {
		t.Fatalf("TigerID: 期望环境变量值 %q，实际 %q", "env_tid", cfg.TigerID)
	}
	// 配置文件值保留
	if cfg.PrivateKey != "file_pk" {
		t.Fatalf("PrivateKey: 期望配置文件值 %q，实际 %q", "file_pk", cfg.PrivateKey)
	}
}

// TestNewClientConfig_WithToken 测试 Token 相关字段
func TestNewClientConfig_WithToken(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := NewClientConfig(
		WithTigerID("tid"),
		WithPrivateKey("pk"),
		WithToken("my_token"),
		WithTokenRefreshDuration(3600*time.Second),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.Token != "my_token" {
		t.Fatalf("Token: 期望 %q，实际 %q", "my_token", cfg.Token)
	}
	if cfg.TokenRefreshDuration != 3600*time.Second {
		t.Fatalf("TokenRefreshDuration: 期望 %v，实际 %v", 3600*time.Second, cfg.TokenRefreshDuration)
	}
}

// TestNewClientConfig_WithLicense 测试 License 字段
func TestNewClientConfig_WithLicense(t *testing.T) {
	clearConfigEnv(t)
	cfg, err := NewClientConfig(
		WithTigerID("tid"),
		WithPrivateKey("pk"),
		WithLicense("TBHK"),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.License != "TBHK" {
		t.Fatalf("License: 期望 %q，实际 %q", "TBHK", cfg.License)
	}
}

// TestNewClientConfig_PKCS1FromFile 测试从配置文件加载 PKCS#1 私钥
func TestNewClientConfig_PKCS1FromFile(t *testing.T) {
	clearConfigEnv(t)
	content := "tiger_id=tid\nprivate_key_pk1=pk1_key_content\naccount=acc\n"
	path := writeTempFile(t, content)

	cfg, err := NewClientConfig(
		WithPropertiesFile(path),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.PrivateKey != "pk1_key_content" {
		t.Fatalf("PrivateKey: 期望 %q，实际 %q", "pk1_key_content", cfg.PrivateKey)
	}
}

// TestNewClientConfig_PKCS8FromFile 测试从配置文件加载 PKCS#8 私钥
func TestNewClientConfig_PKCS8FromFile(t *testing.T) {
	clearConfigEnv(t)
	content := "tiger_id=tid\nprivate_key_pk8=pk8_key_content\naccount=acc\n"
	path := writeTempFile(t, content)

	cfg, err := NewClientConfig(
		WithPropertiesFile(path),
	)
	if err != nil {
		t.Fatalf("创建配置失败: %v", err)
	}

	if cfg.PrivateKey != "pk8_key_content" {
		t.Fatalf("PrivateKey: 期望 %q，实际 %q", "pk8_key_content", cfg.PrivateKey)
	}
}
