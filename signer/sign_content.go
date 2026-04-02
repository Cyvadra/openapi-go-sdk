package signer

import (
	"sort"
	"strings"
)

// GetSignContent 按参数名字母序排列所有参数，拼接为 key=value&key=value 格式。
// 用于构造 RSA 签名的待签名内容。
func GetSignContent(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	// 提取所有键名并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按排序后的键名拼接
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(params[k])
	}

	return builder.String()
}
