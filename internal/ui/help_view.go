package ui

import (
	"ai-ops-agent/pkg/i18n"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *ChatUI) PrintHelpInfo() {
	helpView := tview.NewTextView()
	helpView.SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(true).
		SetBorder(true).
		SetTitle(i18n.T("HelpTitle"))

	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", "/h", i18n.T("CmdHelp"))))
	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", "/clear", i18n.T("CmdClear"))))
	helpView.Write([]byte(fmt.Sprintf("\n\n")))
	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", i18n.T("HeaderShortcuts"), "")))
	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", "tab", i18n.T("ShortcutTab"))))

	helpView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			ui.app.SetRoot(ui.rootLayout, true) // 返回主视图
			ui.app.SetFocus(ui.input)           // 设置聚焦在input框
			return nil
		}
		return event
	})

	ui.app.SetRoot(helpView, true).SetFocus(helpView)

}
