// Package config 提供老虎证券 OpenAPI 的配置管理功能。
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParsePropertiesFile 解析 Java properties 格式的配置文件。
// 支持 key=value 和 key:value 格式，支持 \ 续行，忽略 # 和 ! 注释行及空行。
func ParsePropertiesFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件 %s: %w", path, err)
	}
	defer file.Close()

	props := make(map[string]string)
	scanner := bufio.NewScanner(file)

	var currentLine string
	var continuation bool

	for scanner.Scan() {
		line := scanner.Text()

		if continuation {
			// 续行：去除前导空格后拼接
			line = strings.TrimLeft(line, " \t")
			if strings.HasSuffix(line, "\\") {
				currentLine += line[:len(line)-1]
				continue
			}
			currentLine += line
			continuation = false
		} else {
			// 去除前后空格
			trimmed := strings.TrimSpace(line)

			// 跳过空行
			if trimmed == "" {
				continue
			}

			// 跳过注释行
			if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "!") {
				continue
			}

			// 检查续行
			if strings.HasSuffix(trimmed, "\\") {
				currentLine = trimmed[:len(trimmed)-1]
				continuation = true
				continue
			}

			currentLine = trimmed
		}

		// 解析键值对
		key, value := parseKeyValue(currentLine)
		if key != "" {
			props[key] = value
		}
		currentLine = ""
	}

	// 处理最后一行是续行但文件结束的情况
	if continuation && currentLine != "" {
		key, value := parseKeyValue(currentLine)
		if key != "" {
			props[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	return props, nil
}

// parseKeyValue 解析单行键值对，支持 = 和 : 分隔符。
// 值中可以包含 = 或 :，只按第一个分隔符拆分。
func parseKeyValue(line string) (string, string) {
	// 找到第一个 = 或 : 的位置
	eqIdx := strings.Index(line, "=")
	colonIdx := strings.Index(line, ":")

	sepIdx := -1
	if eqIdx >= 0 && colonIdx >= 0 {
		if eqIdx < colonIdx {
			sepIdx = eqIdx
		} else {
			sepIdx = colonIdx
		}
	} else if eqIdx >= 0 {
		sepIdx = eqIdx
	} else if colonIdx >= 0 {
		sepIdx = colonIdx
	}

	if sepIdx < 0 {
		return "", ""
	}

	key := strings.TrimSpace(line[:sepIdx])
	value := strings.TrimSpace(line[sepIdx+1:])
	return key, value
}
