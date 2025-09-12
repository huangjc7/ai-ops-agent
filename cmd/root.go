package cmd

import (
	"ai-ops-agent/internal/ui"

	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ai",
	Short: "Start the Linux AI TUI assistant",
	Long: `该命令启动主服务，并支持通过环境变量配置以下参数：

环境变量说明：
  • API_KEY   - API 授权密钥（默认：空字符串）
  • BASE_URL  - 模型服务接口地址（默认：https://dashscope.aliyuncs.com/compatible-mode/v1）
  • MODEL     - 模型名称（默认：qwen3-max）

示例：
  API_KEY=yourkey BASE_URL=https://api.openai.com/v1 MODEL=chatGPT-4o ./ai-ops-agent run
`,
	Run: func(cmd *cobra.Command, args []string) {
		chat := ui.NewChatUI()
		if err := chat.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
