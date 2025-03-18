package cmd

import (
	"ai-ops-agent/agent"
	"ai-ops-agent/agent/model"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:     "hugo",
	Short:   "🚀这是一个大模型AI运维助手",
	Version: "v1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 启动AI分析")

		//创建
		bar := progressbar.NewOptions(-1, // 设置 -1 代表未知总步数
			progressbar.OptionSetDescription("AI 任务执行中..."),
			progressbar.OptionSetWidth(30),                   // 限制进度条宽度
			progressbar.OptionShowCount(),                    // 显示计数
			progressbar.OptionThrottle(100*time.Millisecond), // 避免频繁刷新
			progressbar.OptionSpinnerType(14),                // 动态进度条
			progressbar.OptionClearOnFinish(),
		)

		// 初始化出 AI client
		client := model.NewAIClient(&model.ConfigClient{
			ApiKey:  os.Getenv("AI_API_KEY"),
			BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1/",
			Model:   "qwen-plus",
		})

		// 初始化提示工程加入历史上下文
		model.ChatHistory = append(model.ChatHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: agent.Message,
		})

		maxStep := 30 // 最大AI请求次数
		for i := 1; i <= maxStep; i++ {
			// 更新进度条（未知总步数）
			bar.Add(1)
			// 开始询问AI
			client.Send()
			// 解析AI回答并执行命令
			agent.JsonAndExecCommandResponse(model.ChatHistory[len(model.ChatHistory)-1].Content, client)
		}
		fmt.Println("\n⚠️ AI 运维分析超过最大步骤")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
