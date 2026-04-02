package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 1: Properties 配置文件解析 round-trip
// **Validates: Requirements 2.8, 10.7**
//
// 对于任意有效的键值对集合（键和值均为非空字符串，不含特殊字符），
// 将其序列化为 Java properties 格式后再解析，得到的键值对集合应与原始集合等价。
func TestPropertiesRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机键值对
		n := rapid.IntRange(1, 10).Draw(t, "numPairs")
		expected := make(map[string]string, n)

		for i := 0; i < n; i++ {
			// 键：字母数字下划线，不含特殊字符
			key := rapid.StringMatching(`[a-zA-Z_][a-zA-Z0-9_]{0,19}`).Draw(t, fmt.Sprintf("key_%d", i))
			// 值：不含换行、反斜杠、#、! 的可打印字符串，至少包含一个非空格字符
			// Java properties 格式会 trim 值两端空格，所以生成器不生成纯空格值
			val := rapid.StringMatching(`[a-zA-Z0-9.,;@$%^&*()\[\]{}<>/?|~+\-]{1,50}`).Draw(t, fmt.Sprintf("val_%d", i))
			expected[key] = val
		}

		// 序列化为 properties 格式
		var sb strings.Builder
		for k, v := range expected {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}

		// 写入临时文件
		dir := os.TempDir()
		path := filepath.Join(dir, fmt.Sprintf("prop_test_%d.properties", rapid.IntRange(0, 999999).Draw(t, "fileId")))
		if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
			t.Fatalf("写入临时文件失败: %v", err)
		}
		defer os.Remove(path)

		// 解析
		parsed, err := ParsePropertiesFile(path)
		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		// 验证 round-trip
		if len(parsed) != len(expected) {
			t.Fatalf("键值对数量不匹配: 期望 %d，实际 %d", len(expected), len(parsed))
		}
		for k, v := range expected {
			if parsed[k] != v {
				t.Fatalf("键 %q: 期望 %q，实际 %q", k, v, parsed[k])
			}
		}
	})
}
