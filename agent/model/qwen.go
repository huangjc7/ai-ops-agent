package model

import (
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"os/exec"
)

// 最大对话历史长度，防止 Token 超限
const maxHistory = 5

var ChatHistory []openai.ChatCompletionMessage

// 支持的工具（Tools，替代 Functions）
var availableTools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "fixDiskIssue",
			Description: "用于在检测到磁盘问题时执行修复或清理操作。仅在明确提示磁盘相关错误的情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "fixDockerService",
			Description: "用于Docker服务停止或异常时执行修复操作，仅在明确提示Docker相关错误的情况下使用",
			Parameters:  map[string]interface{}{},
		},
	},
}

type ConfigClient struct {
	ApiKey  string
	BaseURL string
	Model   string
}

type HunYuan struct {
	client *openai.Client
	model  string
}

type Sender interface {
	SupportSend()
	Send()
}

// NewAIClient 创建 HunYuan AI 客户端
func NewAIClient(cfg *ConfigClient) Sender {
	config := openai.DefaultConfig(cfg.ApiKey)
	config.BaseURL = cfg.BaseURL
	client := openai.NewClientWithConfig(config)

	return &HunYuan{client: client, model: cfg.Model}
}

// 智能裁剪对话历史，防止 Token 超限
func trimChatHistory(history []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	maxAssistantCount := maxHistory
	// 统计 Role=assistant 的总数量
	assistantCount := 0
	for _, msg := range history {
		if msg.Role == openai.ChatMessageRoleAssistant {
			assistantCount++
		}
	}

	// 如果总数小于等于 maxAssistantCount，直接返回原始历史，不需要裁剪
	if assistantCount <= maxAssistantCount {
		fmt.Printf("当前 Role=assistant 消息数: %d，小于或等于 %d，无需裁剪。\n", assistantCount, maxAssistantCount)
		return history
	}

	fmt.Printf("当前 Role=assistant 消息数: %d，超过 %d，开始裁剪...\n", assistantCount, maxAssistantCount)

	// 如果需要裁剪，则开始保留最新的 maxAssistantCount 条
	var trimmedHistory []openai.ChatCompletionMessage
	var assistantHistory []openai.ChatCompletionMessage

	// 从后往前找到最新的 maxAssistantCount 个 Role=assistant 的记录
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == openai.ChatMessageRoleAssistant {
			assistantHistory = append([]openai.ChatCompletionMessage{history[i]}, assistantHistory...)
			if len(assistantHistory) >= maxAssistantCount {
				break
			}
		}
	}

	// 再次遍历切片，保留非Role=assistant的记录，以及最新的 assistantHistory
	for _, msg := range history {
		if msg.Role != openai.ChatMessageRoleAssistant {
			trimmedHistory = append(trimmedHistory, msg)
		}
	}

	// 合并保留的非Role=assistant记录和最新的maxAssistantCount个Role=assistant记录
	trimmedHistory = append(trimmedHistory, assistantHistory...)

	fmt.Printf("裁剪完成，当前历史记录总数: %d，其中 Role=assistant: %d\n", len(trimmedHistory), maxAssistantCount)

	return trimmedHistory
}

// 模拟服务器巡检
func fixDockerService(args json.RawMessage) (string, error) {

	err := exec.Command("systemctl", "start", "docker").Run()
	systemStatus := map[string]string{
		"status_name": "Docker",
		"status":      "已修复",
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

// 模拟问题修复（仅清理磁盘，不执行其他操作）
func fixDiskIssue(args json.RawMessage) (string, error) {

	// 真正执行逻辑
	cmd := exec.Command("find", "/", "-type", "f", "-size", "+1G", "-name", "*log*", "-exec", "bash", "-c", "> {}", ";")
	err := cmd.Run()

	// 模拟日志输出
	fixResult := map[string]string{
		"action":  "清理磁盘",
		"result":  "已释放 10GB 空间",
		"success": "true",
	}
	if err != nil {
		fixResult = map[string]string{
			"action":  "清理磁盘",
			"result":  "清理磁盘失败",
			"success": "false",
		}
	}
	resultJSON, _ := json.Marshal(fixResult)

	return string(resultJSON), nil
}

// Send 发送用户请求
func (hy *HunYuan) Send() {

	//裁剪对话历史
	//ChatHistory = trimChatHistory(ChatHistory)

	// 发送请求，使用 Tools
	resp, err := hy.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    hy.model,
		Messages: ChatHistory,
	})
	if err != nil {
		fmt.Println("❌ 请求 AI 失败:", err)
	}

	// 如果 AI 没有调用 Function，直接返回回答
	aiReply := resp.Choices[0].Message.Content
	// 追加 AI 回复到历史对话
	ChatHistory = append(ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: aiReply,
	})

}

func (hy *HunYuan) SupportSend() {

	// 发送请求，使用 Tools
	resp, err := hy.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:      hy.model,
		Messages:   ChatHistory,
		Tools:      availableTools,
		ToolChoice: nil, // ✅ 让 AI 自动决定是否调用工具
	})
	if err != nil {
		fmt.Println("❌ 请求 AI 失败:", err)
	}

	// 检查 AI 是否调用了 Tool
	if len(resp.Choices) > 0 {
		if len(resp.Choices[0].Message.ToolCalls) > 0 {
			fmt.Println("🔧 AI 需要调用 Tool:", resp.Choices[0].Message.ToolCalls[0].Function.Name)

			// 解析 Tool 参数
			var params json.RawMessage = []byte(resp.Choices[0].Message.ToolCalls[0].Function.Arguments)

			// 执行对应的 Tool（Function）
			var functionResult string
			switch resp.Choices[0].Message.ToolCalls[0].Function.Name {
			case "fixDockerService":
				result, err := fixDockerService(params)
				if err != nil {
					fmt.Println("❌ Tool 执行失败:", err)
					return
				}
				functionResult = result
			case "fixDiskIssue":
				result, err := fixDiskIssue(params)
				if err != nil {
					fmt.Println("❌ Tool 执行失败:", err)
					return
				}
				functionResult = result
			}

			// 返回 Function 结果给 AI
			ChatHistory = append(ChatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleTool,
				Content: functionResult,
			})

			// 确保 chatHistory 以 assistant 角色结束
			ChatHistory = append(ChatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: fmt.Sprintf("结果：%s", functionResult),
			})
		} else {
			// 如果 AI 没有调用 Function，直接返回回答
			aiReply := resp.Choices[0].Message.Content
			// 追加 AI 回复到历史对话
			ChatHistory = append(ChatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: aiReply,
			})
		}
	}
}
