# AI运维助手
* 自动巡检并处理异常问题
* 最小可行验证方式，仅提供思路，代码部份边界约束条件较不完善
# 实现方式
* prompt + function calling
* 后续可通过RAG增强巡检部份功能，支持特定业务巡检
# 目前适配模型
| 模型          | 是否支持 |
|-------------|------|
| 阿里千文plus    | ✅    |
| 阿里千文max     | ✅    |
| 阿里千文turbo   | ✅    |
| 腾讯混元        | ✅    |
| ChatGPT-4o  | ✅    |
| ChatGPT-4.5 | ✅    |
# 目前支持巡检项目
```text
    - CPU负载
    - 网络联通性
    - 磁盘空间
    - 内存使用
    - SSH服务状态
    - 网络监听端口
    - Docker进程
    - Docker服务运行状态
```
# 使用效果如下
[![asciicast](https://asciinema.org/a/Ht6kxUK5r8WeeHWNYRmsEezrr.svg)](https://asciinema.org/a/Ht6kxUK5r8WeeHWNYRmsEezrr)
# 后续演进方向
Multi-agent，多智能体协作
![示例图片](img/1742663323334.jpg)