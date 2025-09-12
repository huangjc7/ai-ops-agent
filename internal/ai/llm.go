package ai

import (
	"context"
	"encoding/json"
	openai "github.com/sashabaranov/go-openai"
	"log"
)

// 最大对话历史长度，防止 Token 超限
const maxHistory = 5

type OpenClient struct {
	client      *openai.Client
	model       string
	ChatHistory []openai.ChatCompletionMessage
	RawResponse *openai.ChatCompletionResponse
}

type Controller interface {
	Send() *OpenClient                             // 请求AI
	AddUserRoleSession(content string) *OpenClient // 添加对话到历史记录中
	AddSystemRoleSession(content string) *OpenClient
	Close()                // 清空所有历史会话
	PrintResponse() string // 打印最新AI回复
	PrintHistory()         // 打印整个历史对话
	GetHistory() []openai.ChatCompletionMessage
	AddUserRoleMultiContent(contents []openai.ChatMessagePart) *OpenClient // 构造多模态内容
	SendStream(onToken func(string)) error                                 //流式输出
}

func NewAIClient(cfg *Config) Controller {
	config := openai.DefaultConfig(cfg.ApiKey)
	config.BaseURL = cfg.BaseURL
	client := openai.NewClientWithConfig(config)

	return &OpenClient{
		client:      client,
		model:       cfg.Model,
		ChatHistory: make([]openai.ChatCompletionMessage, 0),
	}
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
		log.Printf("当前 Role=assistant 消息数: %d，小于或等于 %d，无需裁剪。\n", assistantCount, maxAssistantCount)
		return history
	}

	log.Printf("当前 Role=assistant 消息数: %d，超过 %d，开始裁剪...\n", assistantCount, maxAssistantCount)

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

	log.Printf("裁剪完成，当前历史记录总数: %d，其中 Role=assistant: %d\n", len(trimmedHistory), maxAssistantCount)

	return trimmedHistory
}

// Send 提供基本对话能力

func (oc *OpenClient) AddUserRoleSession(content string) *OpenClient {
	oc.ChatHistory = append(oc.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	return oc
}
func (oc *OpenClient) AddSystemRoleSession(content string) *OpenClient {
	oc.ChatHistory = append(oc.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: content,
	})
	return oc
}

func (oc *OpenClient) AddCustomRoleSession(role string, content string) {
	oc.ChatHistory = append(oc.ChatHistory, openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
}

func (oc *OpenClient) AddUserRoleMultiContent(contents []openai.ChatMessagePart) *OpenClient {
	oc.ChatHistory = append(oc.ChatHistory, openai.ChatCompletionMessage{
		Role:         openai.ChatMessageRoleUser,
		MultiContent: contents,
	})
	return oc
}

func (oc *OpenClient) SendStream(onToken func(string)) error {
	// 裁剪历史
	// oc.ChatHistory = trimChatHistory(oc.ChatHistory)

	req := openai.ChatCompletionRequest{
		Model:    oc.model,
		Messages: oc.ChatHistory,
		Stream:   true,
	}

	stream, err := oc.client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		log.Println("创建流失败:", err)
		return err
	}
	defer stream.Close()

	var fullReply string

	for {
		resp, err := stream.Recv()
		if err != nil {
			break // io.EOF 表示结束
		}
		token := resp.Choices[0].Delta.Content
		if token != "" {
			fullReply += token
			onToken(token) // 逐 token 回调
		}
	}

	// 追加完整回复到历史
	oc.AddCustomRoleSession(openai.ChatMessageRoleAssistant, fullReply)

	return nil
}

func (oc *OpenClient) Send() *OpenClient {

	//裁剪对话历史 上下文问题比较头疼， 考虑后面版本处理
	//ChatHistory = trimChatHistory(ChatHistory)

	// 发送请求
	resp, err := oc.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    oc.model,
		Messages: oc.ChatHistory,
	})
	oc.RawResponse = &resp
	if err != nil {
		log.Println("请求 AI 失败:", err)
	}

	// 如果 AI 没有调用 Function，直接返回回答
	aiReply := resp.Choices[0].Message.Content

	// 追加 AI 回复到历史对话
	oc.AddCustomRoleSession(openai.ChatMessageRoleAssistant, aiReply)

	return oc
}

func (oc *OpenClient) Close() {
	oc.ChatHistory = oc.ChatHistory[:0]
}

func (oc *OpenClient) PrintHistory() {
	for _, msg := range oc.ChatHistory {
		log.Printf("%s===>%s\n", msg.Role, msg.Content)
	}
}

func (oc *OpenClient) GetHistory() []openai.ChatCompletionMessage { return oc.ChatHistory }

func (oc *OpenClient) PrintResponse() string {
	for i := len(oc.ChatHistory) - 1; i >= 0; i-- {
		msg := oc.ChatHistory[i]
		if msg.Role == openai.ChatMessageRoleAssistant {
			return msg.Content
		}
	}
	return ""
}

func (oc *OpenClient) PrintFullResponse() string {
	if oc.RawResponse == nil {
		return "无响应结果"
	}
	bytes, err := json.MarshalIndent(oc.RawResponse, "", "  ")
	if err != nil {
		return "序列化响应失败: " + err.Error()
	}
	return string(bytes)
}
