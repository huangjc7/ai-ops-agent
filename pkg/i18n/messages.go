package i18n

var messages = map[string]map[string]string{
	"zh": {
		"CmdLong": `该命令启动主服务，并支持通过环境变量配置以下参数：

环境变量说明：
  • API_KEY   - API 授权密钥（默认：空字符串）
  • BASE_URL  - 模型服务接口地址（默认：https://dashscope.aliyuncs.com/compatible-mode/v1）
  • MODEL     - 模型名称（默认：qwen3-max）

示例：
  API_KEY=yourkey BASE_URL=https://api.openai.com/v1 MODEL=chatGPT-4o ./ai-ops-agent
`,
		"NewVersion":          "有新版本，下载新版本得到更好体验。下载地址：https://github.com/huangjc7/ai-ops-agent/releases\n4秒后进入工具",
		"WelcomeMessage":      "[blue]欢迎使用 Linux AI 助手！我可以协助你处理各类运维相关任务。\n输入问题并按 Enter 即可开始对话\n例如，你可以尝试输入：帮我分析系统日志 或 部署一个 Nginx 容器。\n\n输入 /h 并按 Enter 可查看帮助信息\n多轮处理模式： %s[-]",
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
		"DangerousCmd":        "[red][警告] 检测到高风险命令: %s\n是否确认执行？(y/n): [-]",
		"SkipCmd":             "[yellow]已跳过该命令执行[-]\n",
		"SummaryFailed":       "总结请求失败: ",
		"MaxRoundReached":     "[debug]处理轮次达到最大",
		"SummaryRequest":      "请你基于上述对话历史做一个总结，总结一下主体目标 + 目前操作进度和建议",
		"GenNewCmd":           "请生产新的命令组来继续解决上述出现所有的问题",
		"CleanHistory":        "清理会话完毕",
		"ReGenerateCmd":       "请重新生成一组命令并且使用<result>标签对包裹",
		"JudgeContinue":       "\n请你判断一下是否需要继续处理用户的问题，需要则输出<continue>即可，不需要就直接进行总结了",
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
	},
	"en": {
		"CmdLong": `This command starts the main service and supports configuration via environment variables:

Environment Variables:
  • API_KEY   - API Key (Default: empty string)
  • BASE_URL  - Model Service URL (Default: https://dashscope.aliyuncs.com/compatible-mode/v1)
  • MODEL     - Model Name (Default: qwen3-max)

Example:
  API_KEY=yourkey BASE_URL=https://api.openai.com/v1 MODEL=chatGPT-4o ./ai-ops-agent
`,
		"NewVersion":          "New version available. Download for a better experience: https://github.com/huangjc7/ai-ops-agent/releases\nEntering tool in 4 seconds...",
		"WelcomeMessage":      "[blue]Welcome to Linux AI Assistant! I can help you with various operations tasks.\nType your question and press Enter to start chatting.\nFor example: 'Analyze system logs' or 'Deploy an Nginx container'.\n\nType /h and press Enter for help.\nMulti-turn mode: %s[-]",
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
		"DangerousCmd":        "[red][Warning] High-risk command detected: %s\nConfirm execution? (y/n): [-]",
		"SkipCmd":             "[yellow]Command execution skipped[-]\n",
		"SummaryFailed":       "Summary request failed: ",
		"MaxRoundReached":     "[debug]Maximum processing rounds reached",
		"SummaryRequest":      "Please summarize the conversation history, including the main goal, current progress, and suggestions.",
		"GenNewCmd":           "Please generate a new set of commands to continue solving the issues mentioned above.",
		"CleanHistory":        "Session cleared",
		"ReGenerateCmd":       "Please regenerate the command set and wrap it with <result> tags.",
		"JudgeContinue":       "\nPlease determine if further processing is needed. If so, output <continue>. Otherwise, provide a summary.",
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
	},
}
