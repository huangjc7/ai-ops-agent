package tools

import (
	"encoding/json"
	"os/exec"
)

// FixDiskIssue 模拟问题修复（仅清理磁盘，不执行其他操作）
func FixDiskIssue() (string, error) {

	// 真正执行逻辑
	cmd := exec.Command("find", "/", "-type", "f", "-size", "+1G", "-name", "*log*", "-exec", "bash", "-c", "> {}", ";")
	err := cmd.Run()

	// 模拟日志输出
	fixResult := map[string]string{
		"action":  "清理磁盘",
		"result":  "已释放空间,具体请使用df查看",
		"success": "true",
	}

	if err != nil {
		fixResult = map[string]string{
			"action":  "清理磁盘",
			"result":  "清理磁盘失败" + err.Error(),
			"success": "false",
		}
	}

	resultJSON, _ := json.Marshal(fixResult)

	return string(resultJSON), nil
}
