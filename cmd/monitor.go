package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func monitor() {
	fmt.Println("== Server System Information ==")
	fmt.Println("Timestamp:", time.Now().Format(time.RFC1123))
	fmt.Println()

	printHostInfo()
	switch runtime.GOOS {
	case "linux":
		printUptimeAndLoad()
		printCPUInfo()
		printMemInfo()
		printDiskUsage()
		printNetworkStats()
	case "darwin":
		fmt.Println("Mac OS will be supported in future versions.")
	default:
		fmt.Println("Unsupported OS")
	}

}

// Get basic host info
func printHostInfo() {
	fmt.Println("== Host Info ==")
	hostname, _ := os.Hostname()
	fmt.Println("Hostname:", hostname)
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("Architecture:", runtime.GOARCH)

	kernel := getCmdOutput("uname", "-r")
	fmt.Println("Kernel Version:", kernel)

	if content, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				fmt.Println("OS Release:", strings.Trim(line[13:], `"`))
				break
			}
		}
	}
	fmt.Println()
}

func printUptimeAndLoad() {
	fmt.Println("== Uptime & Load ==")
	if content, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(content))
		if len(fields) >= 1 {
			uptimeSec := parseFloat(fields[0])
			fmt.Printf("Uptime: %.2f hours\n", uptimeSec/3600)
		}
	}

	if content, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(content))
		if len(fields) >= 3 {
			fmt.Printf("Load Average (1/5/15 min): %s %s %s\n", fields[0], fields[1], fields[2])
		}
	}
	fmt.Println()
}

func printCPUInfo() {
	fmt.Println("== CPU Info ==")
	content, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			fmt.Println(line)
			break
		}
	}
	fmt.Println("CPU Cores:", runtime.NumCPU())

	// CPU usage (top line from /proc/stat)
	content, err = os.ReadFile("/proc/stat")
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "cpu ") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					user := parseFloat(fields[1])
					nice := parseFloat(fields[2])
					system := parseFloat(fields[3])
					idle := parseFloat(fields[4])
					total := user + nice + system + idle
					usage := (user + nice + system) / total * 100
					fmt.Printf("CPU Usage (approx): %.2f%%\n", usage)
				}
				break
			}
		}
	}
	fmt.Println()
}

func printMemInfo() {
	fmt.Println("== Memory Info ==")
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	info := make(map[string]float64)
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			val := parseFloat(fields[1])
			info[fields[0]] = val / 1024 // convert to MB
		}
	}
	total := info["MemTotal:"]
	free := info["MemFree:"] + info["Buffers:"] + info["Cached:"]
	used := total - free
	fmt.Printf("Total: %.2f MB\n", total)
	fmt.Printf("Used : %.2f MB\n", used)
	fmt.Printf("Free : %.2f MB\n", free)
	fmt.Println()
}

func printDiskUsage() {
	fmt.Println("== Disk Usage ==")
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	total := float64(stat.Blocks*uint64(stat.Bsize)) / 1024 / 1024 / 1024
	free := float64(stat.Bfree*uint64(stat.Bsize)) / 1024 / 1024 / 1024
	used := total - free
	fmt.Printf("Total: %.2f GB\n", total)
	fmt.Printf("Used : %.2f GB\n", used)
	fmt.Printf("Free : %.2f GB\n", free)
	fmt.Println()
}

func printNetworkStats() {
	fmt.Println("== Network Interfaces ==")
	content, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines[2:] {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			iface := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 8 {
				rx := parseFloat(fields[0]) / 1024
				tx := parseFloat(fields[8]) / 1024
				fmt.Printf("Interface %s: RX %.2f KB, TX %.2f KB\n", iface, rx, tx)
			}
		}
	}
	fmt.Println()
}

// Helpers
func getCmdOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
