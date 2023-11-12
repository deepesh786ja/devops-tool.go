package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/fatih/color"
)

func main() {
	rootCmd := &cobra.Command{Use: "mytool"}
	rootCmd.AddCommand(cpuCommand())
	rootCmd.AddCommand(memoryCommand())
	rootCmd.AddCommand(diskSpaceCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cpuCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cpu",
		Short: "Display and continuously update CPU usage in percentage",
		Run:   runCpuUsage,
	}
}

func memoryCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "memory",
		Short: "Display memory usage in MB and GB",
		Run:   getMemoryUsage,
	}
}

func diskSpaceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diskspace",
		Short: "Display disk space usage",
		Run:   getDiskSpace,
	}
}

func runCpuUsage(cmd *cobra.Command, args []string) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				getCpuUsage(nil, nil)
			}
		}
	}()

	// Wait for user to exit
	fmt.Println("Press Ctrl+C to exit...")
	<-stopChan

	// Clean up goroutine
	close(done)
}

func getCpuUsage(cmd *cobra.Command, args []string) {
	cpuInfo, err := cpu.Percent(0, false)
	if err != nil {
		fmt.Println("Error getting CPU usage:", err)
		return
	}

	color.Cyan("CPU Usage: %.2f%%\n", cpuInfo[0])
}

func getMemoryUsage(cmd *cobra.Command, args []string) {
	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error getting memory usage:", err)
		return
	}

	color.Green("Total Memory: %.2f GB\n", float64(memoryInfo.Total)/1024/1024/1024)
	color.Yellow("Free Memory: %.2f GB\n", float64(memoryInfo.Free)/1024/1024/1024)
	color.Red("Used Memory: %.2f GB\n", float64(memoryInfo.Used)/1024/1024/1024)
}

func getDiskSpace(cmd *cobra.Command, args []string) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		fmt.Println("Error getting disk partitions:", err)
		return
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			fmt.Println("Error getting disk space usage:", err)
			continue
		}
		fmt.Printf("Disk: %s\n", partition.Mountpoint)
		color.Cyan("Total Space: %.2f GB\n", float64(usage.Total)/1024/1024/1024)
		color.Green("Free Space: %.2f GB\n", float64(usage.Free)/1024/1024/1024)
		color.Red("Used Space: %.2f GB\n", float64(usage.Used)/1024/1024/1024)
		fmt.Println()
	}
}
