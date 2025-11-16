package ui

import (
	"ai-ops-agent/internal/executor"
	"ai-ops-agent/internal/prompt"
	"ai-ops-agent/pkg/ai"
	"ai-ops-agent/pkg/shell"
	"ai-ops-agent/pkg/system"
	"ai-ops-agent/pkg/text"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatUI struct {
	err            error
	app            *tview.Application
	classSvc       ai.Controller // 专门用于分类（临时对话上下文）
	TmpSvc         ai.Controller // 临时使用
	svc            ai.Controller // 主对话上下文，负责和用户持续交互
	chatView       *tview.TextView
	input          *tview.InputField
	systemInjected bool            // 初始化一次标签
	rootLayout     tview.Primitive // 用于恢复主视图
	execer         *executor.ShellExecutor

	baseURL string
	model   string
	apiKey  string
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 打印欢迎信息
func (ui *ChatUI) printWelcomeSlowly(text string) {
	go func() {
		for _, c := range text {
			ui.chatView.Write([]byte(string(c)))
			ui.app.Draw()
			time.Sleep(30 * time.Millisecond)
		}
		ui.chatView.Write([]byte("\n"))
	}()
}

// NewChatUI 构造函数
func NewChatUI() *ChatUI {

	execer := executor.New(10 * time.Second)
	svc := ai.GetAIModel().TextGenTextModelClient
	classSvc := ai.GetAIModel().TextGenTextModelClient
	tmpSvc := ai.GetAIModel().TextGenTextModelClient

	ui := &ChatUI{
		app:      tview.NewApplication(),
		svc:      svc,
		classSvc: classSvc,
		execer:   execer,
		TmpSvc:   tmpSvc,
	}

	var scrollOffset int
	var autoScroll = true // 控制滚动和历史上翻

	// 创建显示视图
	ui.chatView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			if autoScroll {
				ui.chatView.ScrollToEnd()
			}
			ui.app.Draw() // 刷新整个app
		})

	// 新建输入视图
	ui.input = tview.NewInputField().SetLabel("You: ").SetFieldWidth(0).SetDoneFunc(ui.handleInput)

	// 设置标题
	ui.chatView.SetBorder(true).SetTitle(" Linux AI Assistant ")

	// 设置显示视图的键盘捕捉
	ui.chatView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := ui.chatView.GetScrollOffset()
		switch event.Key() {
		case tcell.KeyUp:
			autoScroll = false
			if row > 0 {
				scrollOffset = row - 1
			}
			ui.chatView.ScrollTo(scrollOffset, 0)
		case tcell.KeyDown:
			autoScroll = false
			scrollOffset = row + 1
			ui.chatView.ScrollTo(scrollOffset, 0)
		case tcell.KeyPgUp:
			autoScroll = false
			scrollOffset = max(0, row-5)
			ui.chatView.ScrollTo(scrollOffset, 0)
		case tcell.KeyPgDn:
			autoScroll = false
			scrollOffset = row + 5
			ui.chatView.ScrollTo(scrollOffset, 0)
		}
		return event
	})

	// 创建一个垂直容器
	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.chatView, 0, 1, true). // 添加显示视图
		AddItem(ui.input, 1, 0, true)     //添加输入视图

	ui.app.SetRoot(layout, true) // 设置容器全屏显示
	ui.app.SetFocus(ui.input)    // 设置聚焦在输入视图

	// 切换聚焦 需要在app里面切换聚焦
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			current := ui.app.GetFocus()

			// 检查当前焦点是否在主视图的组件上
			// 如果焦点不在 input 或 chatView 上，说明在帮助视图或历史视图中
			if current != ui.input && current != ui.chatView {
				// 不在主视图，忽略 Tab 键
				return event
			}

			// 只在主视图中处理 Tab 键切换
			switch current.(type) {
			case *tview.InputField:
				ui.app.SetFocus(ui.chatView)
				autoScroll = true         // 恢复 autoScroll
				ui.chatView.ScrollToEnd() // 立即滚到底部
			case *tview.TextView:
				ui.app.SetFocus(ui.input)
			default:
				ui.app.SetFocus(ui.chatView)
			}
			return nil
		}
		return event
	})
	ui.rootLayout = layout
	return ui
}

func (ui *ChatUI) CleanHistory() {
	ui.svc.Close()
	ui.chatView.Write([]byte("清理会话完毕"))
}

func (ui *ChatUI) showHistory() {
	historyView := tview.NewTextView()
	historyView.SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(true).
		SetBorder(true).
		SetTitle(" 历史对话（按 Esc 返回）")

	var builder strings.Builder
	for _, msg := range ui.svc.GetHistory() {
		role := "[green]助手[-]"
		if msg.Role == "user" {
			role = "[blue]用户[-]"
		} else if msg.Role == "system" {
			role = "[red]系统[-]"
		}
		builder.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}

	historyView.SetText(builder.String())

	historyView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// 不阻断换行或未来扩展
		case tcell.KeyRune:
			if event.Rune() == '/' {
				// 临时进入命令模式，判断是否是 /exit
				ui.input.SetText("/")
				ui.app.SetRoot(ui.rootLayout, true)
				ui.app.SetFocus(ui.input)
				return nil
			}
		case tcell.KeyEsc:
			ui.app.SetRoot(ui.rootLayout, true)
			ui.app.SetFocus(ui.input)
			return nil
		}
		return event
	})

	ui.app.SetRoot(historyView, true).SetFocus(historyView)
}

func (ui *ChatUI) handleInput(key tcell.Key) {
	if key != tcell.KeyEnter {
		return
	}

	userInput := ui.input.GetText()
	if userInput == "" {
		return
	} else if userInput == "/history" {
		ui.showHistory()
		return
	} else if userInput == "/clear" {
		ui.CleanHistory()
		return
	} else if userInput == "/h" || userInput == "/help" {
		ui.PrintHelpInfo()
		return
	}

	ui.chatView.Write([]byte("[yellow]You:[-] " + userInput + "\n"))
	ui.input.SetDisabled(true)
	ui.input.SetText("")

	go ui.AIA(userInput)

}

func (ui *ChatUI) AIA(input string) {

	defer func() {
		ui.input.SetDisabled(false) // 关闭输入框锁定
		ui.app.SetFocus(ui.input)   // 聚焦输入框
		ui.app.Draw()               // 刷新
	}()

	ui.chatView.Write([]byte("[green]AI:[-] "))

	// 判断用户输入
	if shell.IsCommand(input) {
		result := ui.execer.Run(input)
		if result.ExitCode == 0 {
			ui.svc.AddUserRoleSession(input + "的执行结果：" + result.Stdout)
			ui.chatView.Write([]byte(result.Stdout))
		} else {
			ui.svc.AddUserRoleSession(input + "的执行结果：" + result.Stderr)
			ui.chatView.Write([]byte(result.Stdout))
		}
		return
	}

	// 判断类型变化并注入初始化 prompt
	if ui.err = ui.classSvc.AddUserRoleSession(fmt.Sprintf(prompt.Templates[prompt.Class].User, input)).Send(); ui.err != nil {
		ui.chatView.Write([]byte("[red]分类请求失败: " + ui.err.Error() + "[-]\n"))
		return
	}
	replyAi := ui.classSvc.PrintResponse()
	ui.classSvc.Close()

	var inputClass = prompt.InputClassResult
	if ui.err = json.Unmarshal([]byte(replyAi), &inputClass); ui.err != nil {
		ui.chatView.Write([]byte(ui.err.Error()))
		return
	}

	// 判断是否需要更新 prompt 类型
	if !ui.systemInjected {
		ui.svc.AddSystemRoleSession(fmt.Sprintf(prompt.Templates[prompt.InitPrompt].User))
		ui.systemInjected = true
	}

	// 判断是否需要切换类型 Prompt
	switch inputClass.Type {
	case strings.ToLower(prompt.Ask):
		ui.Ask(input)
	case strings.ToLower(prompt.Operation):
		ui.Operation(input)
	default:
		ui.chatView.Write([]byte("[debug]: 没有匹配到类型"))
	}

}

func (ui *ChatUI) Ask(input string) {
	err := ui.svc.
		AddSystemRoleSessionOne(fmt.Sprintf(prompt.Templates[prompt.Ask].System, system.GetSystemInfo())).
		AddUserRoleSession(input).
		SendStream(func(token string) {
			ui.chatView.Write([]byte(token))
		})

	if err != nil {
		ui.chatView.Write([]byte("[red]\n错误: " + ui.err.Error() + "[-]\n"))
		return
	}

	ui.chatView.Write([]byte("\n"))
}

func (ui *ChatUI) Operation(input string) {
	if ui.err = ui.svc.AddSystemRoleSessionOne(fmt.Sprintf(prompt.Templates[prompt.Operation].System, system.GetSystemInfo())).
		AddUserRoleSession(input).Send(); ui.err != nil {
		ui.chatView.Write([]byte("[red]\n错误: " + ui.err.Error() + "[-]\n"))
		return
	}

	// 获取最新的完整回复
	reply := ui.svc.PrintResponse()

	// 如果回复包含 <result> 标签，尝试解析并执行命令
	var commands prompt.SuggestionList
	if strings.Contains(reply, "<result>") {
		resultDatas := text.ExtractAllResults(reply)

		var fmtResult string

		// 防止AI抽风 给多了<result>对
		for _, resultData := range resultDatas {
			ui.err = json.Unmarshal([]byte(resultData), &commands)
			if ui.err != nil {
				ui.chatView.Write([]byte("[red]解析命令失败：" + ui.err.Error() + "[-]\n"))
				return
			}

			for i, command := range commands {

				//desc := strings.TrimSpace(command.Desc) // 去掉首尾空格
				ui.chatView.Write([]byte(fmt.Sprintf("[blue]%d) %s[-]\n", i+1, command.Desc)))

				// 检测高位命令
				if shell.IsDangerousCommandRegex(command.Cmd) {
					ui.chatView.Write([]byte(fmt.Sprintf("[red][警告] 检测到高风险命令: %s\n是否确认执行？(y/n): [-]", command.Cmd)))
					ui.app.QueueUpdateDraw(func() {
						ui.input.SetDisabled(false)
						ui.app.SetFocus(ui.input)
					})

					// 获取用户输入
					confirmed := ui.waitForUserConfirmation()
					if !confirmed {
						ui.chatView.Write([]byte("[yellow]已跳过该命令执行[-]\n"))
						continue
					}

					ui.app.QueueUpdateDraw(func() {
						ui.input.SetDisabled(true)
					})
				}

				// shell执行
				result := ui.execer.Run(command.Cmd)
				if result.ExitCode == 0 {
					fmtResult += fmt.Sprintf("%s 的执行结果: \n%s\n\n", command.Cmd, result.Stdout)
				} else {
					fmtResult += fmt.Sprintf("%s 的执行结果: \n%s\n\n", command.Cmd, result.Stderr)
				}
			}
		}

		if ui.err = ui.TmpSvc.AddUserRoleSession(fmt.Sprintf(prompt.Templates[prompt.Summary].User, fmtResult)).Send(); ui.err != nil {
			ui.chatView.Write([]byte("[red]总结请求失败: " + ui.err.Error() + "[-]\n"))
			return
		}
		cmdExecSummary := ui.TmpSvc.PrintResponse()
		ui.TmpSvc.Close() // 清理动作
		ui.svc.AddUserRoleSession(fmt.Sprintf(prompt.Templates[prompt.FollowupPrompt].User, cmdExecSummary))

		// 重新 Send 一次，继续对话
		ui.svc.SendStream(func(token string) {
			ui.chatView.Write([]byte(token))
		})

		ui.chatView.Write([]byte("\n"))
	}
}

// Run 启动 UI 应用
func (ui *ChatUI) Run() error {

	go func() {
		ui.printWelcomeSlowly(
			"[blue]欢迎使用 Linux AI 助手！输入问题并按 Enter 开始对话[-]\n" +
				"[blue]输入 /h 并按 Enter 可以进入帮助信息[-]")
		ui.app.Draw()
	}()

	return ui.app.Run()
}

func (ui *ChatUI) waitForUserConfirmation() bool {
	inputChan := make(chan bool)

	// 恢复原始输入逻辑
	defer func() {
		ui.app.QueueUpdateDraw(func() {
			ui.input.SetLabel("You: ")
			ui.input.SetText("")
			ui.input.SetDoneFunc(ui.handleInput)
		})
	}()

	ui.app.QueueUpdateDraw(func() {
		ui.input.SetLabel("确认 (y/n): ")
		ui.input.SetText("")
		ui.input.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				resp := strings.TrimSpace(strings.ToLower(ui.input.GetText()))
				ui.input.SetText("")
				if resp == "y" {
					ui.chatView.Write([]byte("[green]用户已确认执行命令[-]\n"))
					inputChan <- true
				} else {
					ui.chatView.Write([]byte("[yellow]用户已取消该命令[-]\n"))
					inputChan <- false
				}
			}
		})
	})

	select {
	case result := <-inputChan:
		return result
	case <-time.After(30 * time.Second):
		ui.chatView.Write([]byte("\n[red]超时未确认，已跳过该命令[-]\n"))
		return false
	}
}
