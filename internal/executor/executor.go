package executor

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

type ShellExecutor struct {
	Timeout time.Duration // 超时时间
}

func New(timeout time.Duration) *ShellExecutor {
	return &ShellExecutor{Timeout: timeout}
}

// Run 执行 Bash 命令，支持超时控制
func (s *ShellExecutor) Run(command string) ExecResult {
	return s.RunInDir(command, "")
}

// RunInDir 在指定工作目录下执行 Bash 命令，支持超时控制
// dir 为空则使用当前进程的工作目录
func (s *ShellExecutor) RunInDir(command string, dir string) ExecResult {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		exitCode = -1
	}

	return ExecResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Err:      err,
	}
}
