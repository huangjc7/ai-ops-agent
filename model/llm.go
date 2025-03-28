package model

import (
	"ai-ops-agent/pkg/tools"
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

// 最大对话历史长度，防止 Token 超限
const maxHistory = 5

var ChatHistory []openai.ChatCompletionMessage

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

// ChatHistory 智能裁剪对话历史，防止 Token 超限
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

// Send 提供基本对话能力
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

// SupportSend 提供修复能力
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
			//var params json.RawMessage = []byte(resp.Choices[0].Message.ToolCalls[0].Function.Arguments)

			// 执行对应的 Tool（Function）
			var functionResult string
			switch resp.Choices[0].Message.ToolCalls[0].Function.Name {
			case "FixDockerService":
				result, err := tools.FixDockerService()
				if err != nil {
					fmt.Println("❌ Tool 执行失败:", err)
					return
				}
				functionResult = result
			case "FixDiskIssue":
				result, err := tools.FixDiskIssue()
				if err != nil {
					fmt.Println("❌ Tool 执行失败:", err)
					return
				}
				functionResult = result
			case "GetTopMemoryProcesses":
				result, err := tools.GetTopMemoryProcesses()
				if err != nil {
					fmt.Println("❌ Tool 执行失败:", err)
					return
				}
				functionResult = result
			case "GetTopCpuProcesses":
				result, err := tools.GetTopCpuProcesses()
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
