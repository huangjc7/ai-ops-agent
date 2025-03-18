package agent

import (
	"ai-ops-agent/agent/model"
	"ai-ops-agent/pkg"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
	"regexp"
	"strings"
)

// CommandResponse 定义 JSON 结构
type CommandResponse struct {
	Step    string `json:"step"`
	Command string `json:"command"`
}

func extractAndParseJSONs(text string) (CommandResponse, error) {
	var result CommandResponse

	// 使用正则表达式匹配第一个JSON对象
	re := regexp.MustCompile(`(?s)\{.*?\}`)
	match := re.FindString(text)

	if match == "" {
		return result, fmt.Errorf("未找到有效的JSON数据")
	}

	// 尝试解析JSON
	if err := json.Unmarshal([]byte(match), &result); err != nil {
		return result, fmt.Errorf("JSON解析失败: %w", err)
	}

	return result, nil
}

// 函数：剔除 JSON 只保留总结部分
func removeJSON(input string) string {
	// 使用正则匹配 JSON 数据
	jsonRegex := regexp.MustCompile(`\{[\s\S]*?\}`)

	// 提取 JSON 部分
	jsonMatches := jsonRegex.FindAllString(input, -1)

	// 如果找到了 JSON 数据，则移除它们
	if len(jsonMatches) > 0 {
		input = jsonRegex.ReplaceAllString(input, "")
	}

	// 去除多余空格或换行
	return strings.TrimSpace(input)
}

// JsonAndExecCommandResponse 解析 JSON 并执行命令
func JsonAndExecCommandResponse(command string, sender model.Sender) {
	// 提取 JSON
	response, err := extractAndParseJSONs(command)
	if err != nil {
		fmt.Println("❌ 未找到有效的 JSON 数据 AI内容为:", command)
		return
	}
	fmt.Println("执行命令: ", response)

	// 如果 AI 说 "结束" 则进行修复阶段
	if response.Step == "结束" {
		fmt.Println("✅ AI 诊断已结束")
		fmt.Println("🧰 开始进入修复阶段")
		// 提取最后总结的异常项字段
		result := matchErrorText(response.Command)
		for index, v := range result {
			fmt.Printf("第%d次修复\n", index+1)
			model.ChatHistory = []openai.ChatCompletionMessage{} // 重置一下上下文 防止上下文溢出 只保留最后总结
			model.ChatHistory = append(model.ChatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: v,
			})
			sender.SupportSend()

			for i, msg := range model.ChatHistory {
				fmt.Printf("Message %d: Role=%s, Content=%s\n", i, msg.Role, msg.Content)
			}
		}
		os.Exit(0)
	}

	// 执行命令
	info := pkg.ExecuteCommand(response.Command)

	model.ChatHistory = append(model.ChatHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: info,
	})
}

func matchErrorText(text string) []string {
	var result []string // 用于存储收集的结果
	scanner := bufio.NewScanner(strings.NewReader(text))
	collecting := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "异常值") || strings.Contains(line, "总结") {
			collecting = true
			continue
		}
		if collecting {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" {
				collecting = false
				continue
			}
			result = append(result, trimmedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取错误：", err)
	}
	return result
}
