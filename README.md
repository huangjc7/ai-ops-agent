
<div align="center">

 <a href="https://github.com/huangjc7/ai-ops-agent">Website</a> •
 <a href="https://github.com/huangjc7/ai-ops-agent/releases">Downloads</a> •
 <a href="./README-zh.md">简体中文</a>

</div>

# AI Operations Assistant

* A fully autonomous operations Agent capable of providing deployment, troubleshooting, analysis, and other operations capabilities.
* Ensures execution safety while utilizing model planning and execution capabilities.

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

# TODO LIST
- Solve context length issues, session management -- Long-term plan
- Solve global shared single session history, use multiple contexts for asynchronous collaborative processing
- Allow users to upload custom tools, MCP
- Internationalization support (Completed)
- ~~Support multi-turn processing repair~~
- ~~Online update version capability~~
- ~~Optimize various code boundaries, log collection, initialization info check, startup prompts (Priority)~~

