package system

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// GetSystemInfo 收集系统基本信息并返回拼接后的字符串
func GetSystemInfo() string {
	var b bytes.Buffer

	b.WriteString("### 系统信息 ###\n")

	// 操作系统和架构
	b.WriteString(fmt.Sprintf("操作系统: %s\n", runtime.GOOS))
	b.WriteString(fmt.Sprintf("架构: %s\n", runtime.GOARCH))

	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		b.WriteString(fmt.Sprintf("主机名: %s\n", hostname))
	}

	// 尝试读取 /etc/os-release 获取 Linux 发行版信息（仅限 Linux）
	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			osRelease := parseOSRelease(string(data))
			if prettyName, ok := osRelease["PRETTY_NAME"]; ok {
				b.WriteString(fmt.Sprintf("发行版: %s\n", prettyName))
			}
		}
	}

	// CPU 核心数
	b.WriteString(fmt.Sprintf("CPU 核心数: %d\n", runtime.NumCPU()))

	// Go 版本
	b.WriteString(fmt.Sprintf("Go版本: %s\n", runtime.Version()))

	return b.String()
}

// parseOSRelease 解析 /etc/os-release 的内容为 map
func parseOSRelease(content string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], `"`)
			result[key] = value
		}
	}
	return result
}
