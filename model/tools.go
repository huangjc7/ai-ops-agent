package model

import "github.com/sashabaranov/go-openai"

// AvailableTools 支持的工具（Tools，替代 Functions）
var availableTools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "FixDiskIssue",
			Description: "用于在检测到磁盘问题时执行修复或清理操作。仅在明确提示磁盘相关错误的情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "FixDockerService",
			Description: "用于Docker服务停止或异常时执行修复操作，仅在明确提示Docker相关错误的情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "GetTopCpuProcesses",
			Description: "用于主机CPU负载高执行的修复操作，仅在明确CPU负载高时情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "GetTopMemoryProcesses",
			Description: "用于主机内存使用率高执行的修复操作，仅在明确内存使用率高时情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
}
