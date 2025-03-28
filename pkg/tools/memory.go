package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// GetTopMemoryProcesses 获取占用内存前 5 的进程，并返回 JSON 格式
func GetTopMemoryProcesses() (string, error) {
	type Process struct {
		PID   int     `json:"pid"`
		Name  string  `json:"name"`
		MemMB float64 `json:"mem_mb"`
	}

	var processes []Process

	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return marshalResult(false, "获取高内存进程", "读取 /proc 目录失败", nil)
	}

	for _, entry := range procEntries {
		if !entry.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		statusPath := fmt.Sprintf("/proc/%d/status", pid)
		data, err := ioutil.ReadFile(statusPath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		var name string
		var rssKB float64
		for _, line := range lines {
			if strings.HasPrefix(line, "Name:") {
				name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			}
			if strings.HasPrefix(line, "VmRSS:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					rssKB, _ = strconv.ParseFloat(fields[1], 64)
				}
			}
		}

		processes = append(processes, Process{
			PID:   pid,
			Name:  name,
			MemMB: rssKB / 1024.0, // 转 MB
		})
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].MemMB > processes[j].MemMB
	})

	if len(processes) > 5 {
		processes = processes[:5]
	}

	return marshalResult(true, "获取高内存进程", "获取成功", processes)
}
