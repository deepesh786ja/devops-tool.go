// main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sysinfo",
	Short: "A CLI tool to display system information",
	Run:   displaySystemInfo,
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func displaySystemInfo(cmd *cobra.Command, args []string) {
	selectedInfo, _ := cmd.Flags().GetString("info")

	switch selectedInfo {
	case "cpu":
		displayCPUInfo()
	case "memory":
		displayMemoryInfo()
	case "disk":
		displayDiskInfo()
	default:
		fmt.Println("Invalid option. Please use 'cpu', 'memory', or 'disk'.")
	}
}

func displayCPUInfo() {
	cpuUsage, err := cpu.Percent(time.Second, true)
	if err != nil {
		fmt.Println("Error fetching CPU usage:", err)
		return
	}
	displayGraph("CPU Usage", cpuUsage, "\033[1;34m") // Blue
}

func displayMemoryInfo() {
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error fetching memory usage:", err)
		return
	}
	displayGraph("Memory Usage", []float64{float64(memUsage.UsedPercent)}, "\033[1;32m", formatBytes(memUsage.Used), formatBytes(memUsage.Total)) // Green
}

func displayDiskInfo() {
	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("Error fetching disk partitions:", err)
		return
	}
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Printf("Error fetching disk usage for %s: %v\n", partition.Mountpoint, err)
			continue
		}
		displayGraph(fmt.Sprintf("Disk Usage (%s)", partition.Mountpoint), []float64{float64(usage.UsedPercent)}, "\033[1;31m", formatBytes(usage.Used), formatBytes(usage.Total)) // Red
	}
}

func displayGraph(title string, data []float64, color string, extraInfo ...string) {
	fmt.Printf("\n%s%s\033[0m\n", color, title) // Color for the title
	if len(extraInfo) > 0 {
		fmt.Printf("Used: %s, Total: %s\n", extraInfo[0], extraInfo[1])
	}
	graph := asciigraph.Plot(data, asciigraph.Height(10))
	fmt.Printf("\033[92m%s\033[0m\n", graph) // Green color for the graph
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	rootCmd.Flags().StringP("info", "i", "cpu", "Specify the system information to display (cpu, memory, disk)")
}

func main() {
	execute()
}
