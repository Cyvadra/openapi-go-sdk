package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// makeTestToken 构造一个 base64 编码的测试 token，前 27 字符包含 "gen_ts,expire_ts"
func makeTestToken(genTsMs int64, expireTsMs int64) string {
	// 前 27 字符格式: "gen_ts_ms,expire_ts_m" 补齐到 27 字符
	header := fmt.Sprintf("%013d,%013d", genTsMs, expireTsMs)
	// header 恰好 27 字符: 13 + 1 + 13
	payload := header + "some_extra_payload_data"
	return base64.StdEncoding.EncodeToString([]byte(payload))
}

func TestTokenManager_LoadToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")
	os.WriteFile(path, []byte("token=test_token_123\n"), 0644)

	m := NewTokenManager(WithTokenFilePath(path))
	token, err := m.LoadToken()
	if err != nil {
		t.Fatalf("加载 Token 失败: %v", err)
	}
	if token != "test_token_123" {
		t.Errorf("Token 应为 test_token_123，实际为 %s", token)
	}
	if m.GetToken() != "test_token_123" {
		t.Error("GetToken 返回值不一致")
	}
}

func TestTokenManager_LoadToken_FileNotFound(t *testing.T) {
	m := NewTokenManager(WithTokenFilePath("/nonexistent/path"))
	_, err := m.LoadToken()
	if err == nil {
		t.Fatal("文件不存在时应返回错误")
	}
}

func TestTokenManager_LoadToken_NoTokenField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")
	os.WriteFile(path, []byte("other_key=value\n"), 0644)

	m := NewTokenManager(WithTokenFilePath(path))
	_, err := m.LoadToken()
	if err == nil {
		t.Fatal("无 token 字段时应返回错误")
	}
}

func TestTokenManager_SetToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")

	m := NewTokenManager(WithTokenFilePath(path))
	err := m.SetToken("new_token_456")
	if err != nil {
		t.Fatalf("设置 Token 失败: %v", err)
	}
	if m.GetToken() != "new_token_456" {
		t.Error("内存中 Token 未更新")
	}

	// 验证文件已更新
	content, _ := os.ReadFile(path)
	if string(content) != "token=new_token_456\n" {
		t.Errorf("文件内容不正确: %s", string(content))
	}
}

func TestTokenManager_SetToken_ThenLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")

	m := NewTokenManager(WithTokenFilePath(path))
	m.SetToken("round_trip_token")

	m2 := NewTokenManager(WithTokenFilePath(path))
	token, err := m2.LoadToken()
	if err != nil {
		t.Fatalf("重新加载 Token 失败: %v", err)
	}
	if token != "round_trip_token" {
		t.Errorf("Token 应为 round_trip_token，实际为 %s", token)
	}
}

func TestTokenManager_AutoRefresh(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")

	// 构造一个 gen_ts 在 100 秒前的 token，refreshDuration 设为 30 秒
	oldGenTs := (time.Now().Unix() - 100) * 1000
	oldToken := makeTestToken(oldGenTs, oldGenTs+3600000)

	m := NewTokenManager(
		WithTokenFilePath(path),
		WithTokenRefreshInterval(50*time.Millisecond),
		WithRefreshDuration(30),
	)
	m.SetToken(oldToken)

	callCount := 0
	m.StartAutoRefresh(func() (string, error) {
		callCount++
		return "refreshed_token", nil
	})

	time.Sleep(200 * time.Millisecond)
	m.StopAutoRefresh()

	if callCount == 0 {
		t.Error("刷新函数应至少被调用一次")
	}
	if m.GetToken() != "refreshed_token" {
		t.Errorf("Token 应为 refreshed_token，实际为 %s", m.GetToken())
	}
}

func TestTokenManager_AutoRefresh_SkipsWhenNotNeeded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token.properties")

	// 构造一个刚刚生成的 token，refreshDuration 设为 3600 秒（远未到期）
	freshGenTs := time.Now().Unix() * 1000
	freshToken := makeTestToken(freshGenTs, freshGenTs+3600000)

	m := NewTokenManager(
		WithTokenFilePath(path),
		WithTokenRefreshInterval(50*time.Millisecond),
		WithRefreshDuration(3600),
	)
	m.SetToken(freshToken)

	callCount := 0
	m.StartAutoRefresh(func() (string, error) {
		callCount++
		return "should_not_be_set", nil
	})

	time.Sleep(200 * time.Millisecond)
	m.StopAutoRefresh()

	if callCount != 0 {
		t.Errorf("Token 未过期时不应调用刷新函数，但被调用了 %d 次", callCount)
	}
}

func TestTokenManager_ShouldTokenRefresh_EmptyToken(t *testing.T) {
	m := NewTokenManager(WithRefreshDuration(30))
	// 空 token 不需要刷新
	if m.ShouldTokenRefresh() {
		t.Error("空 token 不应需要刷新")
	}
}

func TestTokenManager_ShouldTokenRefresh_ZeroDuration(t *testing.T) {
	m := NewTokenManager()
	m.SetToken(makeTestToken(1000000, 2000000))
	// refreshDuration 为 0 时不刷新
	if m.ShouldTokenRefresh() {
		t.Error("refreshDuration 为 0 时不应需要刷新")
	}
}

func TestTokenManager_ShouldTokenRefresh_Expired(t *testing.T) {
	m := NewTokenManager(WithRefreshDuration(30))
	// gen_ts 在 100 秒前
	oldGenTs := (time.Now().Unix() - 100) * 1000
	m.SetToken(makeTestToken(oldGenTs, oldGenTs+3600000))
	if !m.ShouldTokenRefresh() {
		t.Error("Token 已过期应需要刷新")
	}
}

func TestTokenManager_ShouldTokenRefresh_NotExpired(t *testing.T) {
	m := NewTokenManager(WithRefreshDuration(3600))
	// gen_ts 刚刚生成
	freshGenTs := time.Now().Unix() * 1000
	m.SetToken(makeTestToken(freshGenTs, freshGenTs+7200000))
	if m.ShouldTokenRefresh() {
		t.Error("Token 未过期不应需要刷新")
	}
}

func TestTokenManager_StopAutoRefresh_NoOp(t *testing.T) {
	m := NewTokenManager()
	// 未启动时停止不应 panic
	m.StopAutoRefresh()
}
