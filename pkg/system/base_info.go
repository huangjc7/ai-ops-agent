package system

import (
	"ai-ops-agent/pkg/i18n"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// GetSystemInfo 收集系统基本信息并返回拼接后的字符串
func GetSystemInfo() string {
	var b bytes.Buffer

	b.WriteString(i18n.T("SysInfoTitle"))

	// 操作系统和架构
	b.WriteString(fmt.Sprintf(i18n.T("SysOS"), runtime.GOOS))
	b.WriteString(fmt.Sprintf(i18n.T("SysArch"), runtime.GOARCH))

	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		b.WriteString(fmt.Sprintf(i18n.T("SysHostname"), hostname))
	}

	// 尝试读取 /etc/os-release 获取 Linux 发行版信息（仅限 Linux）
	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			osRelease := parseOSRelease(string(data))
			if prettyName, ok := osRelease["PRETTY_NAME"]; ok {
				b.WriteString(fmt.Sprintf(i18n.T("SysDistro"), prettyName))
			}
		}
	}

	// CPU 核心数
	b.WriteString(fmt.Sprintf(i18n.T("SysCPU"), runtime.NumCPU()))
	// Go 版本
	b.WriteString(fmt.Sprintf(i18n.T("SysGoVer"), runtime.Version()))
	u, _ := user.Current()
	b.WriteString(fmt.Sprintf(i18n.T("SysUser"), u.Username))

	// 当前目录
	wd, _ := os.Getwd()
	b.WriteString(fmt.Sprintf(i18n.T("SysPwd"), wd))
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
