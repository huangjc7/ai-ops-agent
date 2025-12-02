# AI Operations Assistant
* ai-ops-agent 是一个基于终端的 AI 助手，专为 Linux 运维和 SRE 设计，利用大语言模型（LLMs）帮助自动化日志分析、配置检查和命令生成，使运维任务更快速、更高效。

# 使用场景
1. **自动化故障排查**:  
   只需告诉 AI 排查问题（例如，“是什么导致我的服务失败？”）。它会自动执行相关命令，如 `systemctl status` 或 `journalctl`，收集日志，并建议下一步的调试操作。

2. **Kubernetes Pod 分析**:  
   询问 AI 分析特定的 Pod 问题（例如，“pod `<pod-name>` 出现了什么问题？”）。它会自动运行 `kubectl describe pod <pod-name>` 以及其他相关命令，并提供关于潜在配置问题的建议。

3. **配置分析**:  
   描述一个配置问题（例如，“解释一下 nginx 配置”）。AI 会自动提取相关配置，进行分析，并提供详细的解释和优化建议。

4. **系统配置生成**:  
   只需用自然语言描述一个系统配置（例如，“为应用 X 设置 nginx 反向代理”）。AI 将根据描述生成必要的部署和配置命令。

# 演示
[![asciicast](https://asciinema.org/a/U53jImXIlvHUB3Gm9cqA4o5tO.svg)](https://asciinema.org/a/U53jImXIlvHUB3Gm9cqA4o5tO)

# 环境变量
| 变量名          | 描述                | 默认值   |
|------------|-------------------|-------|
| `BASE_URL` | 模型调用API地址         | `nil` |
| `API_KEY` | 调认证APIKEY         | `nil` |
| `MODEL` | 模型名称，如"ChatGPT-4o" | `nil` |
| `CONTINUE_COUNT` | 循环处理次数            | `5`   |
| `AGENT_CONTINUE_MODE` |是否启用多轮处理模式，yes开启| `no`  |
| `AI_OPS_LANG` | 语言设置 (en/zh) | `en`  |

# 使用方式
```shell
# 注意区分架构版本
$ curl -o ./ai-ops-agent_linux_amd64.tar.gz https://github.com/huangjc7/ai-ops-agent/releases/download/v2.0.10/ai-ops-agent_linux_amd64.tar.gz
$ tar xf ai-ops-agent_linux_amd64.tar.gz
$ chmod +x ./ai-ops-agent
$ export API_KEY="你的密钥"
$ export BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
$ export MODEL="qwen3-max"
$ ./ai-ops-agent
```

# 目前适配模型
| 模型        | 是否支持 |
|-----------|------|
| 阿里千问  | ✅    |
| 腾讯混元      | ✅    |
| OpenAI         | ✅    |


# TODO LIST
- 解决上下文长度问题，会话管理 -- 长期计划
- 解决全局共享单个会话历史，采用多上下文进行异步协同处理
- 允许用户上传自定义tool、mcp
- ~~国际版支持~~
- ~~支持多轮处理修复~~
- ~~在线更新更新版本能力~~
- ~~优化各类代码边界，日志收集，初始化信息没有填写，启动提示(优先)~~
