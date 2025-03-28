package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// GetTopCpuProcesses 获取占用 CPU 前 5 的进程，并返回 JSON 格式
func GetTopCpuProcesses() (string, error) {
	type Process struct {
		PID  int     `json:"pid"`
		Name string  `json:"name"`
		CPU  float64 `json:"cpu"`
	}

	var processes []Process

	files, err := os.ReadDir("/proc")
	if err != nil {
		return marshalResult(false, "获取高 CPU 进程", "读取 /proc 目录失败", nil)
	}

	for _, file := range files {
		if file.IsDir() {
			pid, err := strconv.Atoi(file.Name())
			if err != nil {
				continue
			}

			statFile := fmt.Sprintf("/proc/%d/stat", pid)
			data, err := os.ReadFile(statFile)
			if err != nil {
				continue
			}

			fields := strings.Fields(string(data))
			if len(fields) < 15 {
				continue
			}

			name := strings.Trim(fields[1], "()")

			userTime, err1 := strconv.Atoi(fields[13])
			systemTime, err2 := strconv.Atoi(fields[14])
			if err1 != nil || err2 != nil {
				continue
			}

			totalTime := userTime + systemTime
			cpu := float64(totalTime) / 100.0

			processes = append(processes, Process{PID: pid, Name: name, CPU: cpu})
		}
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU
	})

	if len(processes) > 5 {
		processes = processes[:5]
	}

	return marshalResult(true, "获取高 CPU 进程", "获取成功", processes)
}

func marshalResult(success bool, actionName string, message string, data interface{}) (string, error) {
	result := map[string]interface{}{
		"action":  actionName,
		"result":  message,
		"success": success,
		"data":    data,
	}
	b, err := json.MarshalIndent(result, "", "  ")
	return string(b), err
}
