package executor

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
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
	cmd, cancel := s.buildCmd(command, dir, nil)
	if cancel != nil {
		defer cancel()
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

// RunStreamingInDirWithEnv 在指定工作目录下“流式”执行 Bash 命令：
// - stdout/stderr 会边执行边写入对应 writer（更像真实终端）
// - 同时仍会捕获 stdout/stderr 到结果中（便于上层做摘要/解析）
// extraEnv 会追加到子进程环境变量中（为空可传 nil）。
func (s *ShellExecutor) RunStreamingInDirWithEnv(command string, dir string, extraEnv []string, stdoutW, stderrW io.Writer) ExecResult {
	// 默认不保留输出内容（更省内存）；需要保留尾部日志时请用 RunStreamingTailInDirWithEnv
	return s.RunStreamingTailInDirWithEnv(command, dir, extraEnv, stdoutW, stderrW, 0)
}

// RunStreamingTailInDirWithEnv 在指定工作目录下“流式”执行 Bash 命令：
// - stdout/stderr 会边执行边写入对应 writer（更像真实终端）
// - 同时可选择性地在结果中保留 stdout/stderr 的“最后 N 行”（用于后续 AI 摘要/分析）
// extraEnv 会追加到子进程环境变量中（为空可传 nil）。
// tailLines<=0 表示不保留输出内容。
func (s *ShellExecutor) RunStreamingTailInDirWithEnv(command string, dir string, extraEnv []string, stdoutW, stderrW io.Writer, tailLines int) ExecResult {
	cmd, cancel := s.buildCmd(command, dir, extraEnv)
	if cancel != nil {
		defer cancel()
	}
	return runStreamingCmd(cmd, stdoutW, stderrW, tailLines)
}

func (s *ShellExecutor) buildCmd(command string, dir string, extraEnv []string) (cmd *exec.Cmd, cancel context.CancelFunc) {
	// Timeout<=0 表示不设超时（由用户 Ctrl+C 或外部机制中断更符合“真终端”体验）
	if s.Timeout <= 0 {
		cmd = exec.Command("bash", "-c", command)
		cancel = nil
	} else {
		var ctx context.Context
		ctx, cancel = context.WithTimeout(context.Background(), s.Timeout)
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	if len(extraEnv) > 0 {
		// 追加环境变量时必须保留父进程环境（否则 PATH 等会丢失）
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	return cmd, cancel
}

type tailLineWriter struct {
	maxLines int

	mu       sync.Mutex
	lines    []string
	partial  string
	maxChars int // 防止极端单行无限增长
}

func newTailLineWriter(maxLines int) *tailLineWriter {
	return &tailLineWriter{
		maxLines: maxLines,
		lines:    make([]string, 0, maxLines),
		maxChars: 1024 * 1024, // 1MB 单行上限
	}
}

func (t *tailLineWriter) Write(p []byte) (int, error) {
	if t.maxLines <= 0 {
		return len(p), nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	s := t.partial + string(p)
	// 限制 partial 增长
	if len(s) > t.maxChars && !strings.Contains(s, "\n") {
		s = s[len(s)-t.maxChars:]
	}

	parts := strings.Split(s, "\n")
	// 最后一段可能是不完整行，先留着
	t.partial = parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	for _, ln := range parts {
		t.lines = append(t.lines, ln)
		if len(t.lines) > t.maxLines {
			// 丢弃更早的行
			t.lines = t.lines[len(t.lines)-t.maxLines:]
		}
	}

	// partial 也要限长，防止无换行的输出爆内存
	if len(t.partial) > t.maxChars {
		t.partial = t.partial[len(t.partial)-t.maxChars:]
	}
	return len(p), nil
}

func (t *tailLineWriter) String() string {
	if t.maxLines <= 0 {
		return ""
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.lines) == 0 && t.partial == "" {
		return ""
	}
	// 这里把 partial 也当作最后一行（即便没有 \n）
	all := append([]string{}, t.lines...)
	if t.partial != "" {
		all = append(all, t.partial)
	}
	return strings.Join(all, "\n")
}

func runStreamingCmd(cmd *exec.Cmd, stdoutW, stderrW io.Writer, tailLines int) ExecResult {
	if stdoutW == nil {
		stdoutW = io.Discard
	}
	if stderrW == nil {
		stderrW = io.Discard
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return ExecResult{ExitCode: -1, Err: err}
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return ExecResult{ExitCode: -1, Err: err}
	}

	var stdoutCap, stderrCap io.Writer = io.Discard, io.Discard
	var stdoutTail, stderrTail *tailLineWriter
	if tailLines > 0 {
		stdoutTail = newTailLineWriter(tailLines)
		stderrTail = newTailLineWriter(tailLines)
		stdoutCap = stdoutTail
		stderrCap = stderrTail
	}

	stdoutMW := io.MultiWriter(stdoutW, stdoutCap)
	stderrMW := io.MultiWriter(stderrW, stderrCap)

	if err := cmd.Start(); err != nil {
		return ExecResult{ExitCode: -1, Err: err}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var copyStdoutErr, copyStderrErr error
	go func() {
		defer wg.Done()
		_, copyStdoutErr = io.Copy(stdoutMW, stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		_, copyStderrErr = io.Copy(stderrMW, stderrPipe)
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	exitCode := 0
	if exitErr, ok := waitErr.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if waitErr != nil {
		exitCode = -1
	}

	// 优先返回 Wait 的错误，其次返回拷贝错误（避免“命令没跑完”被吞掉）
	finalErr := waitErr
	if finalErr == nil {
		if copyStdoutErr != nil {
			finalErr = copyStdoutErr
		} else if copyStderrErr != nil {
			finalErr = copyStderrErr
		}
	}

	return ExecResult{
		Stdout: func() string {
			if stdoutTail != nil {
				return stdoutTail.String()
			}
			return ""
		}(),
		Stderr: func() string {
			if stderrTail != nil {
				return stderrTail.String()
			}
			return ""
		}(),
		ExitCode: exitCode,
		Err:      finalErr,
	}
}
