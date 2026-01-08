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
	"os/exec"
	"path/filepath"
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

	cwd  string // 当前工作目录（用于模拟终端）
	home string // 用户 home（用于 ~ 展开）
}

// New 构造函数
func New() *REPL {
	execer := executor.New(10 * time.Second)
	svc := ai.GetAIModel().TextGenTextModelClient
	classSvc := ai.GetAIModel().TextGenTextModelClient
	tmpSvc := ai.GetAIModel().TextGenTextModelClient

	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = "/"
	}
	home, _ := os.UserHomeDir()

	r := &REPL{
		svc:      svc,
		classSvc: classSvc,
		execer:   execer,
		TmpSvc:   tmpSvc,
		cwd:      cwd,
		home:     home,
	}

	historyFile := ""
	if home != "" {
		historyFile = filepath.Join(home, ".ai-ops-agent_history")
	}

	// 创建 readline 实例，支持行编辑（退格、方向键等）和多行输入
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          r.terminalPrompt(),
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		// 多行输入：按 \ 结尾继续输入下一行
		EnableMask:   false,
		AutoComplete: newTerminalCompleter(r),
	})
	if err != nil {
		// 如果 readline 初始化失败，程序无法正常运行
		fmt.Fprintf(os.Stderr, "Failed to initialize readline: %v\n", err)
		os.Exit(1)
	}

	r.rl = rl

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
	//fmt.Println(i18n.T("HelpTip"))
	fmt.Println()

	// 主循环
	for {
		r.rl.SetPrompt(r.terminalPrompt())
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

		// 处理斜杠命令（兼容旧用法）
		if strings.HasPrefix(input, "/") {
			if r.handleCommand(input) {
				continue
			}
			// /exit 返回 false，退出循环
			break
		}

		// 处理终端输入（命令执行 or AI 交互）
		if cont := r.handleTerminalInput(input); !cont {
			break
		}
		fmt.Fprintln(r.rl.Stdout())
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
			if cont := r.handleTerminalInput(multiInput); !cont {
				return false
			}
			fmt.Fprintln(r.rl.Stdout())
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

// handleTerminalInput 用于判断当前输入是“直接执行的命令”还是“交给 AI 处理的对话/任务”。
// 返回 false 表示需要退出 REPL。
func (r *REPL) handleTerminalInput(input string) bool {
	// 快路径：明显像“要立即执行的命令”，就不走 AI 分类，直接执行。
	if r.looksLikeImmediateShell(input) {
		return r.execShellLine(input)
	}

	// 否则交给 AI 进行分类：shell / ask / operation。
	replyAi, err := r.classSvc.AddUserRoleSession(fmt.Sprintf(prompt.GetTemplate(prompt.Class).User, input)).Send()
	if err != nil {
		fmt.Fprintf(r.rl.Stderr(), "\n[error] %s%s\n", i18n.T("ClassifyFailed"), err.Error())
		return true
	}
	r.classSvc.Close()

	var inputClass = prompt.InputClassResult
	if r.err = json.Unmarshal([]byte(replyAi), &inputClass); r.err != nil {
		fmt.Fprintln(r.rl.Stderr(), r.err.Error())
		return true
	}

	// 只注入一次系统初始化 prompt（用于 AI 对话上下文）。
	if !r.systemInjected {
		r.svc.AddSystemRoleSession(fmt.Sprintf(prompt.GetTemplate(prompt.InitPrompt).User))
		r.systemInjected = true
	}

	switch strings.ToLower(inputClass.Type) {
	case "shell":
		return r.execShellLine(input)
	case strings.ToLower(prompt.Ask):
		r.Ask(input)
		return true
	case strings.ToLower(prompt.Operation):
		r.Operation(input)
		return true
	default:
		// 默认回退到 ask，避免误执行。
		r.Ask(input)
		return true
	}
}

func (r *REPL) terminalPrompt() string {
	dir := r.cwd
	if r.home != "" && strings.HasPrefix(dir, r.home) {
		dir = "~" + strings.TrimPrefix(dir, r.home)
	}
	return colorYellow + dir + " $ " + colorReset
}

func (r *REPL) looksLikeImmediateShell(input string) bool {
	s := strings.TrimSpace(input)
	if s == "" {
		return false
	}
	// 注释行
	if strings.HasPrefix(s, "#") {
		return true
	}

	// “问句”更可能是咨询/解释，应走 AI（避免误执行）。
	if strings.Contains(s, "?") || strings.Contains(s, "？") {
		// 例外：用户明确输入了 --help，本质还是命令。
		if strings.Contains(s, "--help") {
			return true
		}
		return false
	}

	// 中英文“在问命令用法/参数”的提示词：应走 AI（避免把提问当成要执行的命令）。
	askHints := []string{"参数", "用法", "怎么", "如何", "为什么", "是什么", "解释", "含义", "什么意思", "option", "options", "usage", "what", "how", "why"}
	// man xxx 本身就是命令（交互类，后续会走直通执行）。
	if strings.HasPrefix(strings.ToLower(s), "man ") {
		return true
	}
	for _, h := range askHints {
		if strings.Contains(strings.ToLower(s), h) {
			// allow explicit help flag
			if strings.Contains(s, "--help") || strings.Contains(s, "-h") {
				return true
			}
			return false
		}
	}

	fields := strings.Fields(s)
	if len(fields) == 0 {
		return false
	}
	first := fields[0]
	switch first {
	case "cd", "pwd", "clear", "exit", "quit":
		return true
	}
	// 如果第一个 token 能在 PATH 中找到，大概率就是要执行的命令。
	if _, err := exec.LookPath(first); err == nil {
		return true
	}
	return false
}

const shellSentinel = "__AIOPS_AGENT__"

func (r *REPL) execShellLine(input string) bool {
	s := strings.TrimSpace(input)
	if s == "" {
		return true
	}
	if strings.HasPrefix(s, "#") {
		return true
	}

	// 交互/全屏命令（vim/less/top/ssh/...）需要真实 TTY。
	// 这里走“直通模式”：不捕获 stdout/stderr，把终端控制权交给子进程。
	if r.isInteractiveCommandLine(s) {
		if err := r.execInteractive(s); err != nil {
			fmt.Fprintln(r.rl.Stderr(), err.Error())
		}
		return true
	}

	// 基础内建命令：提升“像真终端”的体验（这些应影响 REPL 自己的 cwd/状态）。
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return true
	}
	switch fields[0] {
	case "exit", "quit":
		fmt.Fprintln(r.rl.Stdout(), i18n.T("Goodbye"))
		return false
	case "clear":
		// ANSI 清屏
		fmt.Fprint(r.rl.Stdout(), "\033[H\033[2J")
		return true
	case "pwd":
		fmt.Fprintln(r.rl.Stdout(), r.cwd)
		return true
	case "cd":
		target := ""
		if len(fields) >= 2 {
			target = strings.Join(fields[1:], " ")
		}
		if target == "" {
			target = r.home
		}
		if target == "" {
			target = "/"
		}
		if err := r.changeDir(target); err != nil {
			fmt.Fprintln(r.rl.Stderr(), err.Error())
		}
		return true
	}

	// 用一层包装命令捕获执行后的 cwd，并尽量保留原始 exit code。
	wrapped := fmt.Sprintf(`%s; __ec=$?; printf '\n%sEXIT=%%d\n%sCWD=%%s\n' "$__ec" "$(pwd)"; exit $__ec`, s, shellSentinel, shellSentinel)
	result := r.execer.RunInDir(wrapped, r.cwd)

	stdout, newCwd := stripShellSentinels(result.Stdout)
	if newCwd != "" {
		r.cwd = newCwd
	}

	if stdout != "" {
		fmt.Fprint(r.rl.Stdout(), stdout)
		// 贴近真实终端行为：若输出末尾没有换行，补一个换行。
		if !strings.HasSuffix(stdout, "\n") {
			fmt.Fprintln(r.rl.Stdout())
		}
	}
	if result.Stderr != "" {
		fmt.Fprint(r.rl.Stderr(), result.Stderr)
		if !strings.HasSuffix(result.Stderr, "\n") {
			fmt.Fprintln(r.rl.Stderr())
		}
	}
	return true
}

func (r *REPL) isInteractiveCommandLine(cmd string) bool {
	s := strings.TrimSpace(cmd)
	if s == "" {
		return false
	}

	// 如果管道接到 pager（less/more），通常需要交互能力，按交互命令处理。
	lower := strings.ToLower(s)
	if strings.Contains(lower, "| less") || strings.Contains(lower, "|more") || strings.Contains(lower, "| more") {
		return true
	}

	// 容器/远程工具常见的 TTY 直连参数。
	if strings.Contains(lower, " -it ") || strings.Contains(lower, " --tty ") || strings.Contains(lower, " --interactive ") {
		return true
	}

	fields := strings.Fields(s)
	if len(fields) == 0 {
		return false
	}
	first := strings.ToLower(fields[0])

	switch first {
	case "vi", "vim", "nvim", "nano", "less", "more", "top", "htop", "watch", "man",
		"ssh", "sftp", "scp", "ftp", "telnet",
		"bash", "sh", "zsh", "fish",
		"python", "python3", "node", "irb",
		"mysql", "psql", "sqlite3", "redis-cli":
		return true
	default:
		return false
	}
}

func (r *REPL) execInteractive(cmdline string) error {
	// 先离开 readline 的提示行，再把控制权交给子进程，避免画面错位。
	r.rl.Clean()
	fmt.Fprintln(r.rl.Stdout())

	// 尽力恢复终端到“非 raw 模式”，让交互程序能正常读写 TTY。
	// 注意：这里一般不在 Readline() 阻塞中，readline 的 Terminal goroutine 通常是空闲的。
	_ = r.rl.Terminal.ExitRawMode()
	defer func() { _ = r.rl.Terminal.EnterRawMode() }()

	cmd := exec.Command("bash", "-c", cmdline)
	cmd.Dir = r.cwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w", strings.Fields(cmdline)[0], err)
	}

	// 全屏程序退出后，刷新一次，保证下一次 prompt 绘制干净。
	r.rl.Refresh()
	return nil
}

func stripShellSentinels(stdout string) (clean string, cwd string) {
	lines := strings.Split(stdout, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		if strings.HasPrefix(ln, shellSentinel+"CWD=") {
			cwd = strings.TrimPrefix(ln, shellSentinel+"CWD=")
			continue
		}
		if strings.HasPrefix(ln, shellSentinel+"EXIT=") {
			continue
		}
		out = append(out, ln)
	}
	clean = strings.Join(out, "\n")
	// 去掉一个由于 Split 行尾空元素导致的尾随换行，避免输出多一行空行。
	clean = strings.TrimSuffix(clean, "\n")
	return clean, strings.TrimSpace(cwd)
}

func (r *REPL) changeDir(target string) error {
	t := strings.TrimSpace(target)
	if t == "" {
		return nil
	}
	if strings.HasPrefix(t, "~") && r.home != "" {
		t = filepath.Join(r.home, strings.TrimPrefix(t, "~"))
	}
	if !filepath.IsAbs(t) {
		t = filepath.Join(r.cwd, t)
	}
	t = filepath.Clean(t)

	fi, err := os.Stat(t)
	if err != nil {
		return fmt.Errorf("cd: %s: %v", target, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("cd: %s: not a directory", target)
	}
	r.cwd = t
	return nil
}

// Ask 问答模式（保持核心逻辑不变）
func (r *REPL) Ask(input string) {
	err := r.svc.
		AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Ask).System, system.GetSystemInfo())).
		AddUserRoleSession(input).
		SendStream(func(token string) {
			fmt.Fprint(r.rl.Stdout(), token)
		})

	if err != nil {
		fmt.Fprintf(r.rl.Stderr(), "\n[error] %s\n", err.Error())
		return
	}

	fmt.Fprintln(r.rl.Stdout())
}

// Operation 操作模式（保持核心逻辑不变）
func (r *REPL) Operation(input string) {
	// 只在最外层调用时设置 defer 重置计数器
	if r.repairCount == 0 {
		defer func() { r.repairCount = 0 }()
	}
	var cmdJsonReply string
	var err error

	//r.svc.AddSystemRoleSessionOne(fmt.Sprintf(prompt.GetTemplate(prompt.Operation).System, system.GetSystemInfo()))

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
		fmt.Fprintf(r.rl.Stderr(), "\n[error] %s\n", err.Error())
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
			fmt.Fprintf(r.rl.Stderr(), "[error] %s\n", err.Error())
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
			fmt.Fprintf(r.rl.Stderr(), "[error] %s%s\n", i18n.T("ParseCmdFailed"), r.err.Error())
			return
		}

		fmt.Fprintln(r.rl.Stdout())
		for step, command := range commands {
			fmt.Fprintf(r.rl.Stdout(), "%s:%d %s:%s %s:%s\n", i18n.T("Step"), step+1, i18n.T("Cmd"), command.Cmd, i18n.T("Desc"), command.Desc)
		}

		// 去掉 tview 颜色标签并用 ANSI 颜色
		checkMsg := stripTviewColors(i18n.T("CheckCmdList"))

		// 获取用户确认
		confirmed := r.confirmWithPrompt(colorYellow + checkMsg + colorReset)
		if !confirmed {
			cancelMsg := stripTviewColors(i18n.T("CancelMsg"))
			fmt.Fprintln(r.rl.Stdout(), colorGreen+cancelMsg+colorReset)
			return
		}

		confirmMsg := stripTviewColors(i18n.T("ConfirmMsg"))
		fmt.Fprintln(r.rl.Stdout(), colorGreen+confirmMsg+colorReset)

		for i, command := range commands {
			fmt.Fprintf(r.rl.Stdout(), colorBlue+"%d) %s"+colorReset+"\n", i+1, command.Desc)

			// 检测高危命令
			if shell.IsDangerousCommandRegex(command.Cmd) {
				// 先打印警告信息
				dangerMsg := stripTviewColors(i18n.T("DangerousCmd"))
				fmt.Fprintf(r.rl.Stdout(), colorRed+dangerMsg+colorReset+"\n", command.Cmd)
				// 确认提示符（不包含换行）
				dangerPrompt := colorRed + stripTviewColors(i18n.T("DangerousCmdConfirm")) + colorReset

				// 获取用户确认
				confirmed := r.confirmWithPrompt(dangerPrompt)
				if !confirmed {
					skipMsg := stripTviewColors(i18n.T("SkipCmd"))
					fmt.Fprint(r.rl.Stdout(), colorYellow+skipMsg+colorReset)
					continue
				}
			}

			// shell执行
			result := r.execer.RunInDir(command.Cmd, r.cwd)
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
		fmt.Fprintf(r.rl.Stderr(), "[error] %s%s\n", i18n.T("SummaryFailed"), err.Error())
		return
	}
	r.TmpSvc.Close()

	// 清理包含命令列表的AI回复
	r.svc.RemoveOldResultMessages()

	// 重新 Send 一次，继续对话
	summaryReply, err := r.svc.AddUserRoleSession(cmdExecSummary + i18n.T("JudgeContinue")).Send()
	if err != nil {
		fmt.Fprintf(r.rl.Stderr(), "[error] %s\n", err.Error())
		return
	}

	if !r.continueEnabled {
		fmt.Fprintln(r.rl.Stdout())
		return
	}

	if count, _ := strconv.Atoi(env.Get("CONTINUE_COUNT", "20")); r.repairCount > count-1 {
		fmt.Fprintln(r.rl.Stdout(), colorYellow+i18n.T("MaxRoundReached")+colorReset)

		r.svc.AddUserRoleSession(i18n.T("SummaryRequest")).
			SendStream(func(token string) {
				fmt.Fprint(r.rl.Stdout(), token)
			})
		return
	}

	// 继续
	if strings.Contains(summaryReply, "<continue>") || strings.Contains(summaryReply, "<result>") {
		fmt.Fprintf(r.rl.Stdout(), "\n"+colorCyan+"[CmdSummary] %s"+colorReset+"\n\n", cmdExecSummary)
		r.repairCount++
		r.Operation(i18n.T("GenNewCmd"))
	} else {
		r.svc.SendStream(func(token string) {
			fmt.Fprint(r.rl.Stdout(), token)
		})
	}

	fmt.Fprintln(r.rl.Stdout())
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
			fmt.Fprint(os.Stdout, token)
		})

	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[error] %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
