# AI Operations Assistant
* 完全自主式运维Agent，能够提供部署、排查、分析等多方面运维能力。
* 在完全使用模型规划、执行等能力的同时，保障了执行安全。

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
