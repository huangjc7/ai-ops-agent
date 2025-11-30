package ui

import (
	"ai-ops-agent/pkg/i18n"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const defaultTimeout = 300 * time.Second

// ConfirmationPrompt 管理输入框的确认交互，复用标签、回调和超时逻辑
type ConfirmationPrompt struct {
	app    *tview.Application
	input  *tview.InputField
	output *tview.TextView

	defaultLabel    string
	defaultDoneFunc func(tcell.Key)
	timeout         time.Duration
}

// ConfirmOptions 控制一次确认操作的提示文案和超时
type ConfirmOptions struct {
	PromptLabel    string
	ConfirmMessage string
	CancelMessage  string
	TimeoutMessage string
	Timeout        time.Duration
}

// NewConfirmationPrompt 创建一个确认提示管理器
func NewConfirmationPrompt(app *tview.Application, input *tview.InputField, output *tview.TextView) *ConfirmationPrompt {
	return &ConfirmationPrompt{
		app:     app,
		input:   input,
		output:  output,
		timeout: defaultTimeout,
	}
}

// SetDefaultState 注册默认标签和 DoneFunc，便于确认结束后恢复
func (cp *ConfirmationPrompt) SetDefaultState(label string, doneFunc func(tcell.Key)) {
	cp.defaultLabel = label
	cp.defaultDoneFunc = doneFunc
}

// SetTimeout 设置全局默认超时
func (cp *ConfirmationPrompt) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		cp.timeout = timeout
	}
}

// Confirm 触发一次确认输入，返回用户是否确认
func (cp *ConfirmationPrompt) Confirm(opts ConfirmOptions) bool {
	label := opts.PromptLabel
	if label == "" {
		label = i18n.T("DefaultConfirmLabel")
	}

	confirmMsg := opts.ConfirmMessage
	if confirmMsg == "" {
		confirmMsg = i18n.T("ConfirmMsg")
	}

	cancelMsg := opts.CancelMessage
	if cancelMsg == "" {
		cancelMsg = i18n.T("CancelMsg")
	}

	timeoutMsg := opts.TimeoutMessage
	if timeoutMsg == "" {
		timeoutMsg = i18n.T("TimeoutMsg")
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = cp.timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
	}

	resultCh := make(chan bool, 1)

	cp.app.QueueUpdateDraw(func() {
		cp.input.SetLabel(label)
		cp.input.SetText("")
		cp.input.SetDoneFunc(func(key tcell.Key) {
			if key != tcell.KeyEnter {
				return
			}

			text := strings.TrimSpace(strings.ToLower(cp.input.GetText()))
			cp.input.SetText("")

			switch text {
			case "y", "yes":
				cp.output.Write([]byte(confirmMsg))
				resultCh <- true
			case "n", "no":
				cp.output.Write([]byte(cancelMsg))
				resultCh <- false
			default:
				cp.output.Write([]byte(i18n.T("InvalidInputMsg")))
			}
		})
	})

	defer cp.restoreDefaultState()

	select {
	case result := <-resultCh:
		return result
	case <-time.After(timeout):
		cp.output.Write([]byte(timeoutMsg))
		return false
	}
}

func (cp *ConfirmationPrompt) restoreDefaultState() {
	cp.app.QueueUpdateDraw(func() {
		if cp.defaultLabel != "" {
			cp.input.SetLabel(cp.defaultLabel)
		}
		if cp.defaultDoneFunc != nil {
			cp.input.SetDoneFunc(cp.defaultDoneFunc)
		}
		cp.input.SetText("")
	})
}
