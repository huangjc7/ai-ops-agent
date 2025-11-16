package ui

import (
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
		SetTitle(" 帮助信息（按 Esc 返回）")

	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", "/h", "帮助信息")))
	helpView.Write([]byte(fmt.Sprintf("%-15s %s\n", "/clear", "清除本次对话AI记忆")))

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
