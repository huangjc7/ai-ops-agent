package repl

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chzyer/readline"
)

// terminalCompleter 提供基础的“类 shell”TAB 补全能力：
// - 第一段（命令名）：内建命令 + PATH 里的可执行文件
// - 后续段（参数）：文件/目录补全（相对当前 cwd）
//
// 这里刻意保持解析简单（按空格切分），避免在项目里重复实现完整的 shell 语法解析。
type terminalCompleter struct {
	r *REPL

	once sync.Once
	mu   sync.RWMutex

	commands []string
}

func newTerminalCompleter(r *REPL) readline.AutoCompleter {
	tc := &terminalCompleter{r: r}
	return readline.SegmentAutoComplete(tc)
}

func (c *terminalCompleter) DoSegment(segs [][]rune, _ int) [][]rune {
	c.once.Do(c.buildCommandCache)

	if len(segs) == 0 {
		return nil
	}
	last := string(segs[len(segs)-1])

	// If completing the first token (command), prefer command completion.
	if len(segs) == 1 {
		if strings.Contains(last, "/") || strings.HasPrefix(last, ".") || strings.HasPrefix(last, "~") {
			return c.completePath(last)
		}
		return c.completeCommand(last)
	}

	// Otherwise, complete paths for arguments.
	return c.completePath(last)
}

func (c *terminalCompleter) buildCommandCache() {
	// REPL 直接支持的内建命令。
	builtins := []string{"cd", "pwd", "clear", "exit", "quit"}

	set := map[string]struct{}{}
	for _, b := range builtins {
		set[b] = struct{}{}
	}

	pathEnv := os.Getenv("PATH")
	for _, dir := range strings.Split(pathEnv, string(os.PathListSeparator)) {
		if dir == "" {
			continue
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, ent := range entries {
			name := ent.Name()
			if _, ok := set[name]; ok {
				continue
			}
			info, err := ent.Info()
			if err != nil {
				continue
			}
			// Any executable bit set.
			if info.Mode().IsRegular() && (info.Mode().Perm()&0o111) != 0 {
				set[name] = struct{}{}
			}
		}
	}

	cmds := make([]string, 0, len(set))
	for name := range set {
		cmds = append(cmds, name)
	}
	sort.Strings(cmds)

	c.mu.Lock()
	c.commands = cmds
	c.mu.Unlock()
}

func (c *terminalCompleter) completeCommand(prefix string) [][]rune {
	c.mu.RLock()
	cmds := c.commands
	c.mu.RUnlock()

	out := make([][]rune, 0, 16)
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, prefix) {
			out = append(out, []rune(cmd))
		}
	}
	return out
}

func (c *terminalCompleter) completePath(prefix string) [][]rune {
	// 文件系统操作时需要展开 ~，但返回给 readline 的候选仍尽量保持用户输入的形式（~ / 相对路径）。
	typed := prefix
	typedFS := prefix
	typedOutBase := ""

	if strings.HasPrefix(prefix, "~") {
		home, _ := os.UserHomeDir()
		rest := strings.TrimPrefix(prefix, "~")
		typedFS = filepath.Join(home, rest)
		typedOutBase = "~"
	}

	// 确定要读取的目录（base）以及用户当前正在输入的最后一段前缀（partial）。
	var baseFS, baseOut, partial string
	if strings.HasSuffix(typedFS, string(filepath.Separator)) {
		baseFS = typedFS
		partial = ""
		baseOut = typed
	} else {
		baseFS = filepath.Dir(typedFS)
		partial = filepath.Base(typedFS)
		baseOut = filepath.Dir(typed)
		if baseOut == "." {
			baseOut = ""
		}
	}

	// 如果是相对路径（且不是 ~ 展开后的绝对路径），则以当前 cwd 作为基准。
	if !filepath.IsAbs(baseFS) && typedOutBase == "" {
		baseFS = filepath.Join(c.r.cwd, baseFS)
	}

	entries, err := os.ReadDir(baseFS)
	if err != nil {
		return nil
	}

	out := make([][]rune, 0, 32)
	for _, ent := range entries {
		name := ent.Name()
		if partial != "" && !strings.HasPrefix(name, partial) {
			continue
		}

		cand := name
		if ent.IsDir() {
			cand += string(filepath.Separator)
		}

		// 用“用户输入的路径风格”（相对路径 / ~ 前缀）重建候选，保证补全后看起来像真实终端。
		full := cand
		if strings.HasSuffix(baseOut, string(filepath.Separator)) || baseOut == "" {
			full = baseOut + cand
		} else if baseOut != "" {
			full = baseOut + string(filepath.Separator) + cand
		}

		// 保证候选以用户原始输入前缀开头，便于 readline 的 RetSegment 过滤逻辑工作。
		if !strings.HasPrefix(full, typed) {
			// If the user typed a relative path and baseOut is empty, full is just cand, which should still match.
			continue
		}

		out = append(out, []rune(full))
	}
	return out
}
