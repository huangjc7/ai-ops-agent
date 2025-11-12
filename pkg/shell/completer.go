package shell

import (
	"os"
	"path/filepath"
	"strings"
)

// Completer 提供命令和路径的自动补全功能
type Completer struct {
	commands []string // 常用命令列表
}

// NewCompleter 创建新的补全器
func NewCompleter() *Completer {
	return &Completer{
		commands: []string{
			// 基础命令
			"ls", "cd", "pwd", "cat", "grep", "find", "ps", "top", "htop",
			"kill", "killall", "jobs", "fg", "bg",
			// 文件操作
			"cp", "mv", "rm", "mkdir", "rmdir", "touch", "chmod", "chown",
			"tar", "zip", "unzip", "gzip", "gunzip",
			// 文本处理
			"head", "tail", "less", "more", "sed", "awk", "sort", "uniq",
			"wc", "cut", "paste", "diff", "patch",
			// 网络
			"curl", "wget", "ping", "netstat", "ss", "ifconfig", "ip",
			// 系统服务
			"systemctl", "service", "journalctl",
			// 容器
			"docker", "docker-compose", "kubectl",
			// 版本控制
			"git", "svn",
			// 编辑器
			"vim", "vi", "nano", "emacs",
			// 其他常用
			"sudo", "su", "whoami", "id", "uname", "df", "du", "free",
			"history", "alias", "export", "env", "echo", "printf",
		},
	}
}

// Complete 根据输入返回补全建议
// 返回补全建议列表和当前需要补全的部分（最后一个词）
func (c *Completer) Complete(line string) ([]string, string, int) {
	line = strings.TrimRight(line, " ")
	if line == "" {
		return c.commands, "", 0
	}

	// 找到最后一个空格的位置，确定需要补全的部分
	lastSpaceIdx := strings.LastIndex(line, " ")
	var toComplete string
	var startPos int

	if lastSpaceIdx == -1 {
		// 没有空格，补全第一个词（命令）
		toComplete = line
		startPos = 0
	} else {
		// 有空格，补全最后一个词（可能是路径或参数）
		toComplete = line[lastSpaceIdx+1:]
		startPos = lastSpaceIdx + 1
	}

	// 判断是否是路径补全
	if c.isPath(toComplete) {
		matches := c.completePath(toComplete)
		return matches, toComplete, startPos
	}

	// 命令补全（第一个词）
	if lastSpaceIdx == -1 {
		matches := c.completeCommand(toComplete)
		return matches, toComplete, startPos
	}

	// 参数补全（可能是命令名，也可能是路径）
	// 先尝试作为路径补全
	if strings.Contains(toComplete, "/") || strings.HasPrefix(toComplete, ".") {
		matches := c.completePath(toComplete)
		if len(matches) > 0 {
			return matches, toComplete, startPos
		}
	}

	// 否则尝试作为命令补全
	matches := c.completeCommand(toComplete)
	return matches, toComplete, startPos
}

// isPath 判断是否是路径
func (c *Completer) isPath(s string) bool {
	return strings.HasPrefix(s, "/") ||
		strings.HasPrefix(s, "./") ||
		strings.HasPrefix(s, "../") ||
		strings.HasPrefix(s, "~") ||
		strings.Contains(s, "/")
}

// completeCommand 补全命令
func (c *Completer) completeCommand(prefix string) []string {
	var matches []string
	prefixLower := strings.ToLower(prefix)
	for _, cmd := range c.commands {
		if strings.HasPrefix(strings.ToLower(cmd), prefixLower) {
			matches = append(matches, cmd)
		}
	}
	return matches
}

// completePath 补全文件路径
func (c *Completer) completePath(path string) []string {
	var matches []string
	var searchDir string
	var base string
	originalPath := path
	needsTilde := strings.HasPrefix(originalPath, "~")

	// 处理 ~ 符号
	if needsTilde {
		home, err := os.UserHomeDir()
		if err != nil {
			return matches
		}
		path = strings.Replace(path, "~", home, 1)
	}

	// 分离目录和文件名
	if strings.HasSuffix(path, "/") {
		searchDir = path
		base = ""
	} else {
		searchDir = filepath.Dir(path)
		base = filepath.Base(path)
	}

	// 处理空目录或当前目录
	if searchDir == "" || searchDir == "." {
		searchDir = "."
	}

	// 确保路径是绝对路径（用于读取目录）
	absSearchDir, err := filepath.Abs(searchDir)
	if err != nil {
		return matches
	}

	// 读取目录内容
	entries, err := os.ReadDir(absSearchDir)
	if err != nil {
		return matches
	}

	// 获取原始路径的目录部分（用于构建返回路径）
	origDir := filepath.Dir(originalPath)
	if origDir == "." || origDir == "" {
		origDir = ""
	} else if !strings.HasSuffix(origDir, "/") && origDir != "" {
		origDir = origDir + "/"
	}

	for _, entry := range entries {
		name := entry.Name()
		// 如果 base 为空或者是前缀匹配
		if base == "" || strings.HasPrefix(name, base) {
			var resultPath string

			// 构建返回路径
			if needsTilde {
				// ~ 开头的路径
				home, _ := os.UserHomeDir()
				fullPath := filepath.Join(absSearchDir, name)
				if strings.HasPrefix(fullPath, home) {
					resultPath = "~" + strings.TrimPrefix(fullPath, home)
				} else {
					// 如果不在 home 目录下，使用相对路径
					if origDir != "" {
						resultPath = origDir + name
					} else {
						resultPath = name
					}
				}
			} else if filepath.IsAbs(originalPath) {
				// 绝对路径
				resultPath = filepath.Join(absSearchDir, name)
			} else {
				// 相对路径
				if origDir != "" {
					resultPath = origDir + name
				} else {
					resultPath = name
				}
			}

			// 统一路径分隔符
			resultPath = filepath.ToSlash(resultPath)

			// 如果是目录，添加斜杠
			if entry.IsDir() {
				if !strings.HasSuffix(resultPath, "/") {
					resultPath += "/"
				}
			}

			matches = append(matches, resultPath)
		}
	}

	return matches
}

// GetCommonPrefix 获取字符串列表的最长公共前缀
func GetCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		prefix = commonPrefix(prefix, strs[i])
		if prefix == "" {
			break
		}
	}
	return prefix
}

func commonPrefix(a, b string) string {
	i := 0
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i < minLen && a[i] == b[i] {
		i++
	}
	return a[:i]
}
