package cmd

import (
	"ai-ops-agent/internal/ui"
	"ai-ops-agent/internal/version"
	"ai-ops-agent/pkg/precheck"
	"time"

	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	showVersion bool
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
	PreRunE: func(cmd *cobra.Command, args []string) error {

		ok, err := precheck.CheckVersion()
		if err == nil {
			if !ok {
				fmt.Println("有新版本，下载新版本得到更好体验。4秒后进入工具。")
				time.Sleep(time.Second * 4)
			}
		}

		if err := precheck.CheckVar(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n",
				version.Version, version.Commit, version.BuildDate)
			os.Exit(0)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		chat := ui.NewChatUI()
		if err := chat.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version and exit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
