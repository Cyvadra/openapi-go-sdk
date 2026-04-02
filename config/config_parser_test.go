package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParsePropertiesFile_BasicKeyValue 测试基本键值对解析
func TestParsePropertiesFile_BasicKeyValue(t *testing.T) {
	content := "tiger_id=test123\nprivate_key=mykey\naccount=DU123456\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	assertProp(t, props, "tiger_id", "test123")
	assertProp(t, props, "private_key", "mykey")
	assertProp(t, props, "account", "DU123456")
}

// TestParsePropertiesFile_Comments 测试注释行被忽略
func TestParsePropertiesFile_Comments(t *testing.T) {
	content := "# 这是注释\ntiger_id=abc\n! 这也是注释\naccount=DU001\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if len(props) != 2 {
		t.Fatalf("期望 2 个键值对，实际 %d 个", len(props))
	}
	assertProp(t, props, "tiger_id", "abc")
	assertProp(t, props, "account", "DU001")
}

// TestParsePropertiesFile_EmptyLines 测试空行被忽略
func TestParsePropertiesFile_EmptyLines(t *testing.T) {
	content := "\n\ntiger_id=abc\n\n  \n\naccount=DU001\n\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if len(props) != 2 {
		t.Fatalf("期望 2 个键值对，实际 %d 个", len(props))
	}
}

// TestParsePropertiesFile_Continuation 测试续行（反斜杠续行）
func TestParsePropertiesFile_Continuation(t *testing.T) {
	content := "private_key=MIIEvgIBADANBg\\\nkqhkiG9w0BAQ\\\nEFAASCBKgwgg\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	expected := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwgg"
	assertProp(t, props, "private_key", expected)
}

// TestParsePropertiesFile_TrimSpaces 测试键值两端空格被去除
func TestParsePropertiesFile_TrimSpaces(t *testing.T) {
	content := "  tiger_id  =  test123  \n  account = DU001 \n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	assertProp(t, props, "tiger_id", "test123")
	assertProp(t, props, "account", "DU001")
}

// TestParsePropertiesFile_ColonSeparator 测试冒号分隔符
func TestParsePropertiesFile_ColonSeparator(t *testing.T) {
	content := "tiger_id:test123\naccount:DU001\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	assertProp(t, props, "tiger_id", "test123")
	assertProp(t, props, "account", "DU001")
}

// TestParsePropertiesFile_ValueWithEquals 测试值中包含等号
func TestParsePropertiesFile_ValueWithEquals(t *testing.T) {
	content := "private_key=abc=def=ghi\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	assertProp(t, props, "private_key", "abc=def=ghi")
}

// TestParsePropertiesFile_FileNotFound 测试文件不存在时返回错误
func TestParsePropertiesFile_FileNotFound(t *testing.T) {
	_, err := ParsePropertiesFile("/nonexistent/path/config.properties")
	if err == nil {
		t.Fatal("期望返回错误，但没有")
	}
}

// TestParsePropertiesFile_EmptyFile 测试空文件
func TestParsePropertiesFile_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if len(props) != 0 {
		t.Fatalf("期望 0 个键值对，实际 %d 个", len(props))
	}
}

// TestParsePropertiesFile_ContinuationTrimLeadingSpaces 测试续行时下一行前导空格被去除
func TestParsePropertiesFile_ContinuationTrimLeadingSpaces(t *testing.T) {
	content := "key=hello\\\n    world\n"
	path := writeTempFile(t, content)

	props, err := ParsePropertiesFile(path)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	assertProp(t, props, "key", "helloworld")
}

// --- 辅助函数 ---

// writeTempFile 写入临时文件并返回路径
func writeTempFile(t testing.TB, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.properties")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("写入临时文件失败: %v", err)
	}
	return path
}

// assertProp 断言属性值
func assertProp(t testing.TB, props map[string]string, key, expected string) {
	t.Helper()
	val, ok := props[key]
	if !ok {
		t.Fatalf("缺少键 %q", key)
	}
	if val != expected {
		t.Fatalf("键 %q: 期望 %q，实际 %q", key, expected, val)
	}
}
