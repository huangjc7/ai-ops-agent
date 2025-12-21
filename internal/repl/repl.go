package repl

import (
	"ai-ops-agent/internal/executor"
	"ai-ops-agent/internal/prompt"
	"ai-ops-agent/pkg/ai"
	"ai-ops-agent/pkg/env"
	"ai-ops-agent/pkg/i18n"
	"ai-ops-agent/pkg/shell"
	"ai-ops-agent/pkg/system"
	"ai-ops-agent/pkg/text"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/sashabaranov/go-openai"
)

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// REPL 交互式命令行
type REPL struct {
	err            error
	classSvc       ai.Controller // 专门用于分类（临时对话上下文）
	TmpSvc         ai.Controller // 临时使用
	svc            ai.Controller // 主对话上下文，负责和用户持续交互
	systemInjected bool          // 初始化一次标签
	execer         *executor.ShellExecutor
	rl             *readline.Instance // readline 实例，支持行编辑和多行输入

	repairCount int // 递归计数器

	continueEnabled bool // 开关持续Ai推理模式
}

// New 构造函数
func New() *REPL {
	execer := executor.New(10 * time.Second)
	svc := ai.GetAIModel().TextGenTextModelClient
	classSvc := ai.GetAIModel().TextGenTextModelClient
	tmpSvc := ai.GetAIModel().TextGenTextModelClient

	// 创建 readline 实例，支持行编辑（退格、方向键等）和多行输入
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          colorYellow + "You: " + colorReset,
		HistoryFile:     "", // 可选：设置历史文件路径
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		// 多行输入：按 \ 结尾继续输入下一行
		EnableMask: false,
	})
	if err != nil {
		// 如果 readline 初始化失败，程序无法正常运行
		fmt.Fprintf(os.Stderr, "Failed to initialize readline: %v\n", err)
		os.Exit(1)
	}

	r := &REPL{
		svc:      svc,
		classSvc: classSvc,
		execer:   execer,
		TmpSvc:   tmpSvc,
		rl:       rl,
	}

	if env.Get("AGENT_CONTINUE_MODE", "yes") == "yes" {
		r.continueEnabled = true
	}

	return r
}

// Run 启动 REPL
func (r *REPL) Run() error {
	defer r.rl.Close()

	modeStr := i18n.T("ModeDisabled")
	if r.continueEnabled {
		modeStr = i18n.T("ModeEnabled")
	}

	// 打印欢迎信息（去掉 TUI 颜色标签）
	welcomeMsg := i18n.T("WelcomeMessage")
	welcomeMsg = stripTviewColors(welcomeMsg)
	fmt.Printf(colorCyan+welcomeMsg+colorReset+"\n", modeStr)
	fmt.Println(i18n.T("HelpTip"))
	fmt.Println()

	// 主循环
	for {
		r.rl.SetPrompt(colorYellow + "You: " + colorReset)
		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue // Ctrl+C，继续等待输入
			}
			if err == io.EOF {
				break // Ctrl+D，退出
			}
			return err
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		// 处理命令
		if strings.HasPrefix(input, "/") {
			if r.handleCommand(input) {
				continue
			}
			// /exit 返回 false，退出循环
			break
		}

		// 处理用户输入
		r.AIA(input)
		fmt.Println()
	}

	return nil
}

// readMultiLine 读取多行输入，空行结束
func (r *REPL) readMultiLine() (string, error) {
	fmt.Println(i18n.T("MultilineModeHint"))

	var lines []string

	for {
		r.rl.SetPrompt(colorYellow + " >>> " + colorReset)
		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl+C 取消多行输入
				fmt.Println(i18n.T("MultilineCanceled"))
				return "", nil
			}
			return "", err
		}

		// 空行表示结束输入
		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// stripTviewColors 去掉 tview 的颜色标签 [green] [red] [-] 等
func stripTviewColors(s string) string {
	// 简单替换，去掉 [color] 和 [-] 标签
	s = strings.ReplaceAll(s, "[blue]", "")
	s = strings.ReplaceAll(s, "[green]", "")
	s = strings.ReplaceAll(s, "[red]", "")
	s = strings.ReplaceAll(s, "[yellow]", "")
	s = strings.ReplaceAll(s, "[-]", "")
	return s
}

// handleCommand 处理斜杠命令，返回 true 继续循环，false 退出
func (r *REPL) handleCommand(input string) bool {
	switch input {
	case "/exit", "/quit", "/q":
		fmt.Println(i18n.T("Goodbye"))
		return false
	case "/clear":
		r.CleanHistory()
		return true
	case "/history":
		r.showHistory()
		return true
	case "/h", "/help":
		r.PrintHelpInfo()
		return true
	case "/m", "/multi":
		// 进入多行输入模式
		multiInput, err := r.readMultiLine()
		if err != nil && err != io.EOF {
			fmt.Printf("[error] %s\n", err.Error())
			return true
		}
		multiInput = strings.TrimSpace(multiInput)
		if multiInput != "" {
			r.AIA(multiInput)
			fmt.Println()
		}
		return true
	default:
		fmt.Println(i18n.T("UnknownCommand"))
		return true
	}
}

// CleanHistory 清理历史
func (r *REPL) CleanHistory() {
	r.svc.Close()
	fmt.Println(i18n.T("CleanHistory"))
}

// showHistory 显示历史
func (r *REPL) showHistory() {
	fmt.Println("\n" + i18n.T("HistoryTitle"))
	fmt.Println(strings.Repeat("-", 40))
	for _, msg := range r.svc.GetHistory() {
		role := i18n.T("Assistant")
		if msg.Role == "user" {
			role = i18n.T("UserRole")
		} else if msg.Role == "system" {
			role = i18n.T("SystemRole")
		}
		fmt.Printf("[%s]: %s\n", role, msg.Content)
	}
	fmt.Println(strings.Repeat("-", 40))
}

// PrintHelpInfo 打印帮助信息
func (r *REPL) PrintHelpInfo() {
	fmt.Println("\n" + i18n.T("HelpTitle"))
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("  /help, /h      - " + i18n.T("HelpCmdHelp"))
	fmt.Println("  /m, /multi     - " + i18n.T("HelpCmdMulti"))
	//fmt.Println("  /history       - " + i18n.T("HelpCmdHistory"))
	fmt.Println("  /clear         - " + i18n.T("HelpCmdClear"))
	fmt.Println("  /exit, /quit   - " + i18n.T("HelpCmdExit"))
	fmt.Println(strings.Repeat("-", 40))
}

// AIA 核心逻辑（保持不变）
func (r *REPL) AIA(input string) {
	fmt.Print(colorGreen + "AI: " + colorReset)

	// 判断用户输入
	if strings.HasPrefix(input, "cmd:") {
		// 提取真实命令
		realCmd := strings.TrimPrefix(input, "cmd:")
		result := r.execer.Run(realCmd)
		if result.ExitCode == 0 {
			r.svc.AddUserRoleSession(realCmd + "的执行结果：" + result.Stdout)
			fmt.Print(result.Stdout)
		} else {
			r.svc.AddUserRoleSession(realCmd + "的执行结果：" + result.Stderr)
			fmt.Print(result.Stderr)
		}
		return
	}

	// 判断类型变化并注入初始化 prompt
	replyAi, err := r.classSvc.AddUserRoleSession(fmt.Sprintf(prompt.GetTemplate(prompt.Class).User, input)).Send()
	if err != nil {
		fmt.Printf("\n[error] %s%s\n", i18n.T("ClassifyFailed"), err.Error())
		return
	}
	r.classSvc.Close()

	var inputClass = prompt.InputClassResult
	if r.err = json.Unmarshal([]byte(replyAi), &inputClass); r.err != nil {
		fmt.Println(r.err.Error())
		return
	}

	// 判断是否需要更新 prompt 类型
	if !r.systemInjected {
		r.svc.AddSystemRoleSession(fmt.Sprintf(prompt.GetTemplate(prompt.InitPrompt).User))
		r.systemInjected = true
	}

	// 判断是否需要切换类型 Prompt
	switch inputClass.Type {
	case strings.ToLower(prompt.Ask):
		r.Ask(input)
	case strings.ToLower(prompt.Operation):
		r.Operation(input)
	default:
		fmt.Printf("[warning]: %s\n", i18n.T("NoMatchType"))
	}
}

// Ask 问答模式（保持核心逻辑不变）
func (r *REPL) Ask(input string) {
	err := r.svc.
		AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Ask).System, system.GetSystemInfo())).
		AddUserRoleSession(input).
		SendStream(func(token string) {
			fmt.Print(token)
		})

	if err != nil {
		fmt.Printf("\n[error] %s\n", err.Error())
		return
	}

	fmt.Println()
}

// Operation 操作模式（保持核心逻辑不变）
func (r *REPL) Operation(input string) {
	// 只在最外层调用时设置 defer 重置计数器
	if r.repairCount == 0 {
		defer func() { r.repairCount = 0 }()
	}
	var cmdJsonReply string
	var err error

	r.svc.AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Operation).System, system.GetSystemInfo()))

	if input == "" {
		cmdJsonReply, err = r.svc.
			AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Operation).System, system.GetSystemInfo())).
			Send()
	} else {
		cmdJsonReply, err = r.svc.AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Operation).System, system.GetSystemInfo())).
			AddUserRoleSession(input).
			Send()
	}

	if err != nil {
		fmt.Printf("\n[error] %s\n", err.Error())
		return
	}

	for retryCount := 0; retryCount < 3; retryCount++ {
		// 找到了就退出
		if strings.Contains(cmdJsonReply, "<result>") {
			break
		}

		// 没找到，重新请求
		cmdJsonReply, err = r.svc.AddUserRoleSession(i18n.T("ReGenerateCmd")).Send()
		if err != nil {
			fmt.Printf("[error] %s\n", err.Error())
			return
		}
	}

	// 执行命令添加对话历史，方便Ai回溯
	r.svc.AddCustomRoleSession(openai.ChatMessageRoleAssistant, cmdJsonReply)

	resultDatas := text.ExtractAllResults(cmdJsonReply)

	var fmtResult string
	var commands prompt.SuggestionList

	// 防止AI抽风 多给了<result>标签对
	for _, resultData := range resultDatas {
		r.err = json.Unmarshal([]byte(resultData), &commands)
		if r.err != nil {
			fmt.Printf("[error] %s%s\n", i18n.T("ParseCmdFailed"), r.err.Error())
			return
		}

		fmt.Println()
		for step, command := range commands {
			fmt.Printf("%s:%d %s:%s %s:%s\n", i18n.T("Step"), step+1, i18n.T("Cmd"), command.Cmd, i18n.T("Desc"), command.Desc)
		}

		// 去掉 tview 颜色标签并用 ANSI 颜色
		checkMsg := stripTviewColors(i18n.T("CheckCmdList"))

		// 获取用户确认
		confirmed := r.confirmWithPrompt(colorYellow + checkMsg + colorReset)
		if !confirmed {
			cancelMsg := stripTviewColors(i18n.T("CancelMsg"))
			fmt.Println(colorGreen + cancelMsg + colorReset)
			return
		}

		confirmMsg := stripTviewColors(i18n.T("ConfirmMsg"))
		fmt.Println(colorGreen + confirmMsg + colorReset)

		for i, command := range commands {
			fmt.Printf(colorBlue+"%d) %s"+colorReset+"\n", i+1, command.Desc)

			// 检测高危命令
			if shell.IsDangerousCommandRegex(command.Cmd) {
				// 先打印警告信息
				dangerMsg := stripTviewColors(i18n.T("DangerousCmd"))
				fmt.Printf(colorRed+dangerMsg+colorReset+"\n", command.Cmd)
				// 确认提示符（不包含换行）
				dangerPrompt := colorRed + stripTviewColors(i18n.T("DangerousCmdConfirm")) + colorReset

				// 获取用户确认
				confirmed := r.confirmWithPrompt(dangerPrompt)
				if !confirmed {
					skipMsg := stripTviewColors(i18n.T("SkipCmd"))
					fmt.Print(colorYellow + skipMsg + colorReset)
					continue
				}
			}

			// shell执行
			result := r.execer.Run(command.Cmd)
			if result.ExitCode == 0 {
				fmtResult += fmt.Sprintf(i18n.T("ExecResult"), command.Cmd, result.Stdout)
			} else {
				fmtResult += fmt.Sprintf(i18n.T("ExecResult"), command.Cmd, result.Stderr)
			}
		}
	}

	// 截断以防止超出 AI 上下文限制
	const maxResultSize = 200 * 1024
	if len(fmtResult) > maxResultSize {
		fmtResult = fmtResult[:maxResultSize] + "\n\n[Truncated due to length limit]..."
	}

	// 提炼描述
	cmdExecSummary, err := r.TmpSvc.AddUserRoleSession(fmt.Sprintf(prompt.GetTemplate(prompt.Summary).User, fmtResult)).Send()
	if err != nil {
		fmt.Printf("[error] %s%s\n", i18n.T("SummaryFailed"), err.Error())
		return
	}
	r.TmpSvc.Close()

	// 清理包含命令列表的AI回复
	r.svc.RemoveOldResultMessages()

	// 重新 Send 一次，继续对话
	summaryReply, err := r.svc.AddUserRoleSession(cmdExecSummary + i18n.T("JudgeContinue")).Send()
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
	}

	if !r.continueEnabled {
		fmt.Println()
		return
	}

	if count, _ := strconv.Atoi(env.Get("CONTINUE_COUNT", "5")); r.repairCount > count-1 {
		fmt.Println(colorYellow + i18n.T("MaxRoundReached") + colorReset)

		fmt.Print(colorGreen + "AI: " + colorReset)

		r.svc.AddUserRoleSession(i18n.T("SummaryRequest")).
			SendStream(func(token string) {
				fmt.Print(token)
			})
		return
	}

	// 继续
	if strings.Contains(summaryReply, "<continue>") || strings.Contains(summaryReply, "<result>") {
		fmt.Printf("\n"+colorCyan+"[debug] %s"+colorReset+"\n\n", cmdExecSummary)
		r.repairCount++
		r.Operation(i18n.T("GenNewCmd"))
	} else {
		fmt.Print(colorGreen + "AI: " + colorReset)
		r.svc.SendStream(func(token string) {
			fmt.Print(token)
		})
	}

	fmt.Println()
}

// confirmWithPrompt 带提示的 y/n 确认
func (r *REPL) confirmWithPrompt(prompt string) bool {
	// 保存原始提示符
	originalPrompt := r.rl.Config.Prompt

	// 设置确认提示符
	r.rl.SetPrompt(prompt)

	defer r.rl.SetPrompt(originalPrompt)

	for {
		input, err := r.rl.Readline()
		if err != nil {
			return false
		}

		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Print(i18n.T("InvalidInputMsg"))
			r.rl.SetPrompt(prompt) // 重新设置提示符
		}
	}
}

// RunSingleAsk 单次问答模式（用于管道输入）
func (r *REPL) RunSingleAsk(question string) {
	r.svc.AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Ask).System, system.GetSystemInfo()))

	err := r.svc.
		AddUserRoleSession(question).
		SendStream(func(token string) {
			fmt.Print(token)
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[error] %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
