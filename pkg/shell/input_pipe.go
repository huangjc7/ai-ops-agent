package shell

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// 从管道读取输入
func GetInputFromPipe() string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		var builder strings.Builder
		for {
			line, err := reader.ReadString('\n')
			builder.WriteString(line)
			if err == io.EOF {
				break
			}
		}
		return strings.TrimSpace(builder.String())
	}
	return ""
}
