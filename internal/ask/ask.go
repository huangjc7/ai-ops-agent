package ask

import (
	"ai-ops-agent/internal/prompt"
	"ai-ops-agent/pkg/ai"
	"ai-ops-agent/pkg/system"
	"fmt"
	"os"
)

// Run 非交互 ask 模式 后续解耦AIA方法
func Run(question string) {
	svc := ai.GetAIModel().TextGenTextModelClient

	// 添加系统 prompt + 用户问题，流式输出
	err := svc.
		AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Ask).System, system.GetSystemInfo())).
		AddUserRoleSession(question).
		SendStream(func(token string) {
			fmt.Print(token)
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Println() // 结尾换行
}
