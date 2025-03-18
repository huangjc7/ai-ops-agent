package pkg

import (
	"bytes"
	"os/exec"
)

// ExecuteCommand 执行 shell 命令，实时打印，并返回输出内容
func ExecuteCommand(cmd string) string {
	//fmt.Println("🚀 执行命令:", cmd)

	// 创建命令
	out := exec.Command("bash", "-c", cmd)

	// 用于存储执行结果
	var outputBuffer bytes.Buffer

	out.Stdout = &outputBuffer
	out.Stderr = &outputBuffer

	// 启动并等待完成
	out.Run()

	// 返回完整输出
	//fmt.Println("✅ 命令执行完成，结果:\n", result)
	return outputBuffer.String()
}
