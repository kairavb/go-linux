package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	runMonitor bool
	dryRun     bool
)

var rootCmd = &cobra.Command{
	Use:   "lihelp [description]",
	Short: "lihelp helps new Linux users by generating commands from natural language",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if runMonitor {
			monitor()
		}

		userQuery := strings.Join(args, " ")

		command, explanation, err := GetCommandFromAI(userQuery)
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}

		fmt.Println("\n🔹 Command:\n   ", command)
		fmt.Println("\n📘 Explanation:\n   ", explanation)

		logCommand(userQuery, command, explanation)

		if dryRun {
			fmt.Println("\n🚫 Dry Run Mode: Command not executed.")
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n❓ Run this command? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))

		if confirm == "y" {
			if needsSudo(command) && os.Geteuid() != 0 {
				fmt.Print("🔐 This command may require sudo. Proceed with sudo? (y/n): ")
				sudoConfirm, _ := reader.ReadString('\n')
				sudoConfirm = strings.TrimSpace(strings.ToLower(sudoConfirm))
				if sudoConfirm == "y" {
					command = "sudo " + command
				}
			}
			runCommand(command)
		} else {
			fmt.Println("✅ Cancelled.")
		}
	},
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show and interact with command history",
	Run: func(cmd *cobra.Command, args []string) {
		showHistory()
	},
}

func Execute() {
	rootCmd.Flags().BoolVarP(&runMonitor, "monitor", "m", false, "Show system monitor info")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview command without running")
	rootCmd.AddCommand(historyCmd)
	rootCmd.Execute()
}
