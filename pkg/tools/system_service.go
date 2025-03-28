package tools

import (
	"encoding/json"
	"os/exec"
)

// FixDockerService模拟服务器巡检
func FixDockerService() (string, error) {

	err := exec.Command("systemctl", "start", "docker").Run()
	systemStatus := map[string]string{
		"status_name": "Docker",
		"status":      "已启动",
	}
	if err != nil {
		systemStatus = map[string]string{
			"status_name": "Docker",
			"status":      "启动失败",
		}
	}
	// 转换为 JSON
	statusJSON, _ := json.Marshal(systemStatus)
	return string(statusJSON), nil
}
