package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var logFile = "~/.lihelp_commands.log"

// 🛠 Run the command using shell
func runCommand(cmd string) {
	if strings.TrimSpace(cmd) == "" {
		fmt.Println("❌ Empty command.")
		return
	}

	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	fmt.Println("\n🏃 Running:", cmd)
	err := execCmd.Run()
	if err != nil {
		fmt.Println("❌ Error running command:", err)
	} else {
		fmt.Println("✅ Command executed successfully.")
	}
}

// 🧠 Detect if command might need sudo
func needsSudo(cmd string) bool {
	privileged := []string{"sudo", "apt", "dnf", "yum", "pacman", "rm -rf /", "systemctl", "reboot", "shutdown", "chmod", "chown", "mount", "umount"}

	for _, keyword := range privileged {
		if strings.Contains(cmd, keyword) {
			return true
		}
	}
	return false
}

// 📒 Log command + explanation to file
func logCommand(userQuery, command, explanation string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("⚠️ Could not retrieve home directory:", err)
		return
	}
	logFile := filepath.Join(homeDir, ".lihelp_commands.log")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("⚠️ Could not log command:", err)
		return
	}
	defer f.Close()

	entry := fmt.Sprintf(
		"\n---\n🕒 %s\n📝 Query: %s\n💻 Command: %s\n📘 Explanation: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		userQuery,
		command,
		explanation,
	)

	if _, err := f.WriteString(entry); err != nil {
		fmt.Println("⚠️ Error writing to log file:", err)
	}
}
