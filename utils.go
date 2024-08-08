package sqlca

import (
	"regexp"
	"strings"
)

// convertCamelToSnake converts a CamelCase string to snake_case
func convertCamelToSnake(s string) string {
	// 使用正则表达式匹配大写字母前的所有字符
	re := regexp.MustCompile("([A-Z])")
	// 将匹配到的大写字母替换为下划线加小写字母
	snake := re.ReplaceAllString(s, "_$1")
	// 将字符串转换为小写并去掉开头的下划线
	return strings.ToLower(strings.TrimPrefix(snake, "_"))
}

