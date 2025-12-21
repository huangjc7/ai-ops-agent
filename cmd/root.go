package cmd

import (
	"ai-ops-agent/internal/ask"
	"ai-ops-agent/internal/repl"
	"ai-ops-agent/internal/version"
	"ai-ops-agent/pkg/i18n"
	"ai-ops-agent/pkg/precheck"
	"ai-ops-agent/pkg/shell"
	"time"

	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	showVersion bool
)

var rootCmd = &cobra.Command{
	Use:   "ai",
	Short: "Start the Linux AI assistant",
	Long:  i18n.T("CmdLong"),
	PreRunE: func(cmd *cobra.Command, args []string) error {

		ok, err := precheck.CheckVersion()
		if err == nil {
			if !ok {
				fmt.Println(i18n.T("NewVersion"))
				time.Sleep(time.Second * 4)
			}
		}

		if err := precheck.CheckVar(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n",
				version.Version, version.Commit, version.BuildDate)
			os.Exit(0)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// 检查管道输入（非交互模式，只支持 ask）
		input := shell.GetInputFromPipe()
		if input != "" {
			ask.Run(input)
			return
		}

		// 交互模式 REPL
		r := repl.New()
		if err := r.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version and exit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
