
<div align="center">

 <a href="https://github.com/huangjc7/ai-ops-agent">Website</a> •
 <a href="https://github.com/huangjc7/ai-ops-agent/releases">Downloads</a> •
 <a href="./README-zh.md">简体中文</a>

</div>

# AI Operations Assistant
* ai-ops-agent is a terminal-based AI assistant for Linux ops and SREs, powered by large language models (LLMs). It helps automate log analysis, configuration inspections, and command generation, making ops tasks faster and more efficient.

# Use Cases
1. **Automated Troubleshooting**:  
   Simply ask the AI to troubleshoot an issue (e.g., "What is causing my service to fail?"). It will automatically run relevant commands, such as `systemctl status` or `journalctl`, gather logs, and suggest the next steps for debugging.

2. **Kubernetes Pod Analysis**:  
   Ask the AI to analyze a specific pod issue (e.g., "What is wrong with pod `<pod-name>`?"). It will automatically run `kubectl describe pod <pod-name>` and other relevant commands, and provide recommendations on potential misconfigurations.

3. **Configuration Analysis**:  
   Describe a configuration issue (e.g., "Explain the nginx configuration"). The AI will automatically fetch relevant configurations, analyze them, and provide detailed explanations and optimization suggestions.

4. **System Configuration Generation**:  
   Simply describe a system configuration in natural language (e.g., "Set up a reverse proxy with nginx for app X"). The AI will generate the necessary commands for deployment and configuration based on the description.

# Demo
[![asciicast](https://asciinema.org/a/R8mG62leelpF5GNJcJc6l9hog.svg)](https://asciinema.org/a/R8mG62leelpF5GNJcJc6l9hog)

# Environment Variables
| Variable Name | Description | Default Value |
|------------|-------------------|---------------|
| `BASE_URL` | Model API URL | `nil`         |
| `API_KEY` | Authentication API KEY | `nil`         |
| `MODEL` | Model name, e.g., "ChatGPT-4o" | `nil`         |
| `CONTINUE_COUNT` | Max loop count for processing | `5`           |
| `AGENT_CONTINUE_MODE` | Enable multi-turn processing mode (yes to enable) | `no`          |
| `AI_OPS_LANG` | Language setting (en/zh) | `en`          |

# Usage
```shell
# Note: Choose the correct architecture version
$ curl -o ./ai-ops-agent_linux_amd64.tar.gz https://github.com/huangjc7/ai-ops-agent/releases/download/v2.0.11/ai-ops-agent_linux_amd64.tar.gz
$ tar xf ai-ops-agent_linux_amd64.tar.gz
$ chmod +x ./ai-ops-agent
$ export API_KEY="your_api_key"
$ export BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
$ export MODEL="qwen3-max"
$ ./ai-ops-agent
```

# Supported Models
| Model | Supported |
|-----------|------|
| Qwen (Aliyun) | ✅ |
| Hunyuan (Tencent) | ✅ |
| OpenAI | ✅ |

# Features
- **Automated Log Analysis**: AI automatically analyzes logs (e.g., `journalctl`, `kubectl`) and provides insights, suggesting the next steps for troubleshooting.
- **Configuration Insights**: AI reviews system configurations (e.g., systemd, Kubernetes) and explains their impact or potential issues.
- **Command Suggestions**: Based on your logs or system state, AI generates the appropriate commands for the next troubleshooting steps.
- **Human-in-the-loop**: AI assists with decision-making, but you stay in control — review and approve commands before execution.
- **Lightweight**: No agents or additional software needed on target machines — the tool runs directly from your control machine.

# TODO LIST
- Solve context length issues, session management -- Long-term plan
- Solve global shared single session history, use multiple contexts for asynchronous collaborative processing
- Allow users to upload custom tools, MCP
- Internationalization support (Completed)
- ~~Support multi-turn processing repair~~
- ~~Online update version capability~~
- ~~Optimize various code boundaries, log collection, initialization info check, startup prompts (Priority)~~

