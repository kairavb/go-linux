package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	runMonitor bool
	dryRun     bool
	logFile    = "lihelp_commands.log"
)

func main() {
	// Parse flags
	args := os.Args[1:]
	var userQueryParts []string

	for _, arg := range args {
		if arg == "-m" {
			runMonitor = true
		} else if arg == "-d" {
			dryRun = true
		} else {
			userQueryParts = append(userQueryParts, arg)
		}
	}

	if len(userQueryParts) == 0 {
		fmt.Println("Usage: lihelp [-m] [-d] \"your linux task description here\"")
		return
	}

	if runMonitor {
		monitor()
	}

	userQuery := strings.Join(userQueryParts, " ")

	// Call AI handler
	command, explanation, err := GetCommandFromAI(userQuery)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("\n🔹 Command:")
	fmt.Println("   ", command)

	fmt.Println("\n📘 Explanation:")
	fmt.Println("   ", explanation)

	// Log the command
	logCommand(userQuery, command, explanation)

	// Ask user if they want to execute it
	if dryRun {
		fmt.Println("\n🚫 Dry Run Mode: Command not executed.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n❓ Run this command? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		// Ask for sudo if needed
		if needsSudo(command) && os.Geteuid() != 0 {
			fmt.Print("🔐 This command may require sudo. Proceed with sudo? (y/n): ")
			sudoConfirm, _ := reader.ReadString('\n')
			sudoConfirm = strings.TrimSpace(strings.ToLower(sudoConfirm))

			if sudoConfirm == "y" || sudoConfirm == "yes" {
				command = "sudo " + command
			} else {
				fmt.Println("⚠️ Skipping sudo. Attempting to run without elevated privileges.")
			}
		}
		runCommand(command)
	} else {
		fmt.Println("✅ Cancelled. Command not run.")
	}
}

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
	// Check for some common privileged commands
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
