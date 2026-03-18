package i18n

var messages = map[string]map[string]string{
	"zh": {
		"CmdLong": `该命令启动主服务，并支持通过环境变量配置以下参数：

环境变量说明：
  • API_KEY             - API 授权密钥（默认：空字符串）
  • BASE_URL            - 模型服务接口地址（默认：https://dashscope.aliyuncs.com/compatible-mode/v1）
  • MODEL               - 模型名称（默认：qwen3-max）
  • AGENT_CONTINUE_MODE - 启用多轮处理模式（默认：启用）
  • AGENT_CONFIRM_MODE  - 命令执行二次确认开关（默认：关闭）
  • CONTINUE_COUNT      - 模型推理最大循环轮次（默认：5次）
  • AI_OPS_LANG         - 界面语言(默认：英语)

示例：
  API_KEY=yourkey BASE_URL=https://api.openai.com/v1 MODEL=chatGPT-4o ./ai-ops-agent
`,
		"NewVersion":          "有新版本，下载新版本得到更好体验。下载地址：https://github.com/huangjc7/ai-ops-agent/releases\n4秒后进入工具",
		"WelcomeMessage":      "[blue]欢迎使用 Linux AI 助手！我可以协助你处理各类运维相关任务。\n输入问题并按 Enter 即可开始对话\n例如，你可以尝试输入：帮我分析系统日志 或 部署一个 Nginx 容器。\n\n输入 /h 并按 Enter 可查看帮助信息，输入 /exit 退出\n多轮处理模式： %s  命令确认模式： %s[-]",
		"HistoryTitle":        " 历史对话（按 Esc 返回）",
		"Assistant":           "助手",
		"UserRole":            "用户",
		"SystemRole":          "系统",
		"ClassifyFailed":      "分类请求失败: ",
		"NoMatchType":         "没有匹配到类型",
		"ParseCmdFailed":      "解析命令失败：",
		"Step":                "步骤",
		"Cmd":                 "命令",
		"Desc":                "描述",
		"CheckCmdList":        "[yellow][提示] 执行命令清单检查,是否确认上述执行？(y/n): [-]",
		"DangerousCmd":        "[警告] 检测到高风险命令: %s",
		"DangerousCmdConfirm": "是否确认执行？(y/n): ",
		"SkipCmd":             "[yellow]已跳过该命令执行[-]\n",
		"CmdExecuted":         "✔ 已执行",
		"CmdCancelled":        "✘ 已取消",
		"CmdSkipped":          "- 跳过（前序命令已取消）",
		"ExecResultCancelled": "命令: %s [用户取消]\n\n",
		"ExecResultSkipped":   "命令: %s [未执行 - 前序命令已取消]\n\n",
		"SummaryFailed":       "总结请求失败: ",
		"MaxRoundReached":     "[debug]处理轮次达到最大",
		"SummaryRequest": "请你基于上述对话历史做一个“任务总结/记忆”，要求如下：\n" +
			"1) 必须包含：主体目标、当前状态（是否解决/卡点）、已执行的关键动作（按时间列出3-10条，包含关键命令与关键输出结论）、关键证据（错误信息/指标/日志要点）、未解决原因假设（如仍未解决）、下一步建议（最多5条，按优先级）。\n" +
			"2) 必须简洁，不要复述大量原始输出；只提炼关键点。\n" +
			"3) 输出纯文本，不要使用 markdown。\n" +
			"4) 这份总结将作为后续对话的唯一上下文记忆来源，请确保信息足够让后续继续排查。",
		"GenNewCmd":           "请生产新的命令组来继续解决上述出现所有的问题",
		"CleanHistory":        "清理会话完毕",
		"ReGenerateCmd":       "请重新生成一组命令并且使用<result>标签对包裹",
		"JudgeContinue":       "\n上述为用户反馈阶段（Observation），请你判断一下是否需要继续处理用户的问题，需要则输出<continue>即可，不需要的话就直接进行总结",
		"ModeEnabled":         "启用",
		"ModeDisabled":        "关闭",
		"ConfirmCheck":        "执行命令清单检查,是否确认上述执行？(y/n): ", // Re-check usage
		"ExecResult":          "%s 的执行结果: \n%s\n\n",
		"ErrBaseUrlMissing":   "没有配置BASE_URL环境变量",
		"ErrModelMissing":     "没有配置MODEL环境变量",
		"ErrApiKeyMissing":    "没有配置API_KEY环境变量",
		"HelpTitle":           " 帮助信息（按 Esc 返回）",
		"CmdHelp":             "帮助信息",
		"CmdClear":            "清除本次对话AI记忆",
		"HeaderShortcuts":     "快捷键",
		"ShortcutTab":         "切换聚焦框",
		"DefaultConfirmLabel": "确认 (y/n): ",
		"ConfirmMsg":          "[green]用户已确认执行命令[-]\n",
		"CancelMsg":           "[green]用户已取消该命令[-]\n",
		"TimeoutMsg":          "\n[red]超时未确认，已跳过该命令[-]\n",
		"InvalidInputMsg":     "[yellow]请输入 y 或 n[-]\n",
		"SysInfoTitle":        "### 系统信息 ###\n",
		"SysOS":               "操作系统: %s\n",
		"SysArch":             "架构: %s\n",
		"SysHostname":         "主机名: %s\n",
		"SysDistro":           "发行版: %s\n",
		"SysCPU":              "CPU 核心数: %d\n",
		"SysGoVer":            "Go版本: %s\n",
		"SysUser":             "当前登录用户 %s \n",
		"SysPwd":              "当前目录位置 %s \n",
		"HelpCmdMulti":        "进入多行输入模式",
		"MultilineModeHint":   "[多行模式] 每行按 Enter 换行，输入空行（再按一次 Enter）发送，Ctrl+C 取消",
		"MultilineCanceled":   "已取消多行输入",
		"Goodbye":             "再见！",
		"UnknownCommand":      "未知命令，输入 /help 查看帮助",
		"HelpCmdHelp":         "显示帮助信息",
		"HelpCmdHistory":      "显示对话历史",
		"HelpCmdClear":        "清除对话历史",
		"HelpCmdExit":         "退出程序",
		"HelpCmdCapture":      "执行命令并将输出采集到对话历史（例：?? tail -100 /var/log/messages）",
		"CapturedDirect":      "\n[captured] 已直接写入对话历史，可直接继续提问分析。",
		"CapturedSummary":     "\n[captured] 已生成摘要写入对话历史，可直接继续提问分析。",
		"CapturedDirectCtx":   "我执行了命令：%s\n\n%s",
		"CapturedSummaryCtx":  "我执行了命令：%s\n（仅保留输出最后500行并生成摘要）\n摘要：%s",
	},
	"en": {
		"CmdLong": `This command starts the main service and supports configuration via environment variables:

Environment Variables:
  • API_KEY             - API Key (Default: empty string)
  • BASE_URL            - Model Service URL (Default: https://dashscope.aliyuncs.com/compatible-mode/v1)
  • MODEL               - Model Name (Default: qwen3-max)
  • AGENT_CONTINUE_MODE - Enable multi-turn processing mode（Default：enable）
  • AGENT_CONFIRM_MODE  - Command execution confirmation toggle（Default：disable）
  • CONTINUE_COUNT      - Max loop count for processing（Default：five times）
  • AI_OPS_LANG         - Language setting(Default：English)

Example:
  API_KEY=yourkey BASE_URL=https://api.openai.com/v1 MODEL=chatGPT-4o ./ai-ops-agent
`,
		"NewVersion":          "New version available. Download for a better experience: https://github.com/huangjc7/ai-ops-agent/releases\nEntering tool in 4 seconds...",
		"WelcomeMessage":      "[blue]Welcome to Linux AI Assistant! I can help you with various operations tasks.\nType your question and press Enter to start chatting.\nFor example: 'Analyze system logs' or 'Deploy an Nginx container'.\n\nType /h and press Enter for help.\nMulti-turn mode: %s  Command confirm mode: %s[-]",
		"HistoryTitle":        " History (Press Esc to return)",
		"Assistant":           "Assistant",
		"UserRole":            "User",
		"SystemRole":          "System",
		"ClassifyFailed":      "Classification failed: ",
		"NoMatchType":         "No matching type found",
		"ParseCmdFailed":      "Failed to parse command: ",
		"Step":                "Step",
		"Cmd":                 "Command",
		"Desc":                "Description",
		"CheckCmdList":        "[yellow][Prompt] Confirm execution of the command list? (y/n): [-]",
		"DangerousCmd":        "[Warning] High-risk command detected: %s",
		"DangerousCmdConfirm": "Confirm execution? (y/n): ",
		"SkipCmd":             "[yellow]Command execution skipped[-]\n",
		"CmdExecuted":         "✔ Executed",
		"CmdCancelled":        "✘ Cancelled",
		"CmdSkipped":          "- Skipped (previous command was cancelled)",
		"ExecResultCancelled": "Command: %s [Cancelled by user]\n\n",
		"ExecResultSkipped":   "Command: %s [Not executed - previous command was cancelled]\n\n",
		"SummaryFailed":       "Summary request failed: ",
		"MaxRoundReached":     "[debug]Maximum processing rounds reached",
		"SummaryRequest": "Please create a concise 'task summary/memory' based on the conversation history:\n" +
			"1) Must include: main goal, current status (solved/stuck), key actions taken (3-10 items with key commands and key conclusions), key evidence (errors/metrics/log highlights), hypotheses if still unresolved, next steps (max 5, prioritized).\n" +
			"2) Keep it concise; do NOT paste large raw outputs.\n" +
			"3) Output plain text only (no markdown).\n" +
			"4) This summary will be the only memory for subsequent turns; make it sufficient to continue troubleshooting.",
		"GenNewCmd":           "Please generate a new set of commands to continue solving the issues mentioned above.",
		"CleanHistory":        "Session cleared",
		"ReGenerateCmd":       "Please regenerate the command set and wrap it with <result> tags.",
		"JudgeContinue":       "\nThe above is the user feedback phase (Observation). Please determine if further processing is needed. If so, output <continue>. Otherwise, provide a summary.",
		"ModeEnabled":         "Enabled",
		"ModeDisabled":        "Disabled",
		"ExecResult":          "Result of %s: \n%s\n\n",
		"ErrBaseUrlMissing":   "BASE_URL environment variable not set",
		"ErrModelMissing":     "MODEL environment variable not set",
		"ErrApiKeyMissing":    "API_KEY environment variable not set",
		"HelpTitle":           " Help (Press Esc to return)",
		"CmdHelp":             "Help information",
		"CmdClear":            "Clear current session AI memory",
		"HeaderShortcuts":     "Shortcuts",
		"ShortcutTab":         "Switch focus",
		"DefaultConfirmLabel": "Confirm (y/n): ",
		"ConfirmMsg":          "[green]User confirmed command execution[-]\n",
		"CancelMsg":           "[green]User canceled the command[-]\n",
		"TimeoutMsg":          "\n[red]Confirmation timed out, command skipped[-]\n",
		"InvalidInputMsg":     "[yellow]Please enter y or n[-]\n",
		"SysInfoTitle":        "### System Info ###\n",
		"SysOS":               "OS: %s\n",
		"SysArch":             "Arch: %s\n",
		"SysHostname":         "Hostname: %s\n",
		"SysDistro":           "Distro: %s\n",
		"SysCPU":              "CPU Cores: %d\n",
		"SysGoVer":            "Go Version: %s\n",
		"SysUser":             "Current User: %s \n",
		"SysPwd":              "Current Directory: %s \n",
		"HelpTip":             "Type /help for help, /exit to quit",
		"HelpCmdMulti":        "Enter multiline input mode",
		"MultilineModeHint":   "[Multiline Mode] Press Enter to add new line, empty line (press Enter again) to send, Ctrl+C to cancel",
		"MultilineCanceled":   "Multiline input canceled",
		"Goodbye":             "Goodbye!",
		"UnknownCommand":      "Unknown command, type /help for help",
		"HelpCmdHelp":         "Show help information",
		"HelpCmdHistory":      "Show conversation history",
		"HelpCmdClear":        "Clear conversation history",
		"HelpCmdExit":         "Exit the program",
		"HelpCmdCapture":      "Run a command and capture output into chat history (e.g. ?? tail -100 /var/log/messages)",
		"CapturedDirect":      "\n[captured] Written directly into chat history. You can continue asking questions.",
		"CapturedSummary":     "\n[captured] Summary written into chat history. You can continue asking questions.",
		"CapturedDirectCtx":   "I ran command: %s\n\n%s",
		"CapturedSummaryCtx":  "I ran command: %s\n(kept last 500 lines and generated summary)\nSummary: %s",
	},
}
