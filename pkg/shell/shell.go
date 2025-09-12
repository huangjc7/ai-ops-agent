package shell

import (
	"regexp"
	"unicode"
)
import "strings"

// IsCommand 判断是否为命令行表达式
func IsCommand(input string) bool {
	input = strings.TrimSpace(input)

	// 如果包含中文，则认为是自然语言
	for _, r := range input {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}

	// 匹配命令结构（英文+符号组合）
	// 允许的字符集：字母数字空格和典型命令行符号
	validCmdPattern := regexp.MustCompile(`^[a-zA-Z0-9\s/\-_.|&><=~:"'\\]+$`)
	return validCmdPattern.MatchString(input)
}

var dangerousPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*rm\s+`),                                   // rm 命令
	regexp.MustCompile(`(?i)\b(reboot|shutdown|halt|poweroff|restart)\b`), // 重启类
	regexp.MustCompile(`(?i)\bsystemctl\s+(stop|disable)\b`),
	regexp.MustCompile(`(?i)^\s*dd\s+`),
	regexp.MustCompile(`(?i)\b(wipefs|mkfs|parted|drop|delete)\b`),
}

func IsDangerousCommandRegex(cmd string) bool {
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(cmd) {
			return true
		}
	}
	return false
}
