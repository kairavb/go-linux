package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func showHistory() {
	file, err := os.Open("~/.lihelp_commands.log")
	if err != nil {
		fmt.Println("⚠️ Could not open log file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []string
	var entry strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			if entry.Len() > 0 {
				entries = append(entries, entry.String())
				entry.Reset()
			}
			continue
		}
		entry.WriteString(line + "\n")
	}
	if entry.Len() > 0 {
		entries = append(entries, entry.String())
	}

	if len(entries) == 0 {
		fmt.Println("📭 No history found.")
		return
	}

	fmt.Println("📜 Command History:")
	for i, e := range entries {
		fmt.Printf("[%d]\n%s\n", i+1, e)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🔢 Enter history number to re-run or press Enter to exit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("👋 Exiting history view.")
		return
	}

	index := -1
	fmt.Sscanf(input, "%d", &index)
	if index <= 0 || index > len(entries) {
		fmt.Println("❌ Invalid choice.")
		return
	}

	selected := entries[index-1]
	fmt.Println("\n📦 Selected Entry:\n", selected)

	fmt.Print("\n❓ Run this command again? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm == "y" {
		lines := strings.Split(selected, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "💻 Command:") {
				cmd := strings.TrimPrefix(line, "💻 Command:")
				runCommand(strings.TrimSpace(cmd))
				break
			}
		}
	} else {
		fmt.Println("✅ Cancelled.")
	}
}
