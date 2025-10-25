# AI运维助手
* 对话式进行运维操作，可以部署、排查、分析等能力，... 比较考验模型能力
* 最小可行验证方式，仅提供思路，代码部份边界约束条件较不完善，可尝试玩一玩；后续有时间或者issues留言完善
# 使用方式
```shell
# 注意区分架构版本
$ curl -o ./ai-ops-agent_linux_amd64.tar.gz https://github.com/huangjc7/ai-ops-agent/releases/download/v2.0.0/ai-ops-agent_linux_amd64.tar.gz
$ tar xf ai-ops-agent_linux_amd64.tar.gz 
$ chmod +x ./ai-ops-agent
$ export API_KEY="你的密钥"
$ export BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
$ export MODEL="qwen3-max"
$ ./ai-ops-agent
```

# 使用效果演示
[![asciicast](https://asciinema.org/a/xVjkj1DYvElhxT2fTmQVehGSM.svg)](https://asciinema.org/a/xVjkj1DYvElhxT2fTmQVehGSM)

# 实现方式
* 纯Prompt实现，依赖模型自身规划能力

# 目前适配模型
| 模型        | 是否支持 |
|-----------|------|
| 阿里千问  | ✅    |
| 腾讯混元      | ✅    |
| OpenAI         | ✅    |


# 后续代办
- 解决上下文长度问题，会话管理
- 解决全局共享单个会话历史，采用多上下文进行异步协同处理
- 允许用户上传自定义tool、mcp
- 在线更新更新版本能力
- 将Ask、Agent模式手工切换，ai判断的不好用