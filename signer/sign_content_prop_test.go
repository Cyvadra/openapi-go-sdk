package signer

import (
	"sort"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Feature: multi-language-sdks, Property 5: 请求参数按字母序排列
// **Validates: Requirements 3.3**
//
// 对于任意参数名-值的映射（map），GetSignContent 函数输出的字符串中，
// 参数应严格按参数名的字母序排列，格式为 key1=value1&key2=value2&...
func TestProperty5_SignContentAlphabeticalOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// 生成随机参数映射（1-10 个键值对）
		// 键名：仅包含字母和数字，避免特殊字符干扰
		// 值：任意非空字符串（不含 & 和换行符，以便解析验证）
		numParams := rapid.IntRange(1, 10).Draw(t, "numParams")
		params := make(map[string]string)
		for i := 0; i < numParams; i++ {
			key := rapid.StringMatching(`[a-z][a-z0-9_]{0,15}`).Draw(t, "key")
			value := rapid.StringMatching(`[a-zA-Z0-9_.\-/: ]{1,30}`).Draw(t, "value")
			params[key] = value
		}

		// 跳过空映射（去重后可能为空）
		if len(params) == 0 {
			return
		}

		result := GetSignContent(params)

		// 验证 1：结果非空
		if result == "" {
			t.Fatal("非空参数映射的结果不应为空")
		}

		// 验证 2：按 & 分割后，每个部分都是 key=value 格式
		parts := strings.Split(result, "&")
		if len(parts) != len(params) {
			t.Fatalf("期望 %d 个参数，实际 %d 个", len(params), len(parts))
		}

		// 验证 3：提取所有键名，检查是否严格按字母序排列
		keys := make([]string, 0, len(parts))
		for _, part := range parts {
			idx := strings.Index(part, "=")
			if idx < 0 {
				t.Fatalf("参数格式错误，缺少 '=': %q", part)
			}
			key := part[:idx]
			value := part[idx+1:]
			keys = append(keys, key)

			// 验证 4：每个键值对的值与原始映射一致
			expectedValue, ok := params[key]
			if !ok {
				t.Fatalf("结果中包含未知的键: %q", key)
			}
			if value != expectedValue {
				t.Fatalf("键 %q 的值不匹配：期望 %q，实际 %q", key, expectedValue, value)
			}
		}

		// 验证 5：键名严格按字母序排列
		if !sort.StringsAreSorted(keys) {
			t.Fatalf("参数未按字母序排列: %v", keys)
		}
	})
}
