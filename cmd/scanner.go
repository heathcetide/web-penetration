package cmd

import (
	"fmt"
	"web-penetration/internal/scanner"

	"github.com/spf13/cobra"
)

var scannerCmd = &cobra.Command{
	Use:   "scan",
	Short: "端口扫描功能",
	Long:  `对指定主机进行端口扫描，检测开放端口和服务信息。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标主机")
			fmt.Println("用法: web-pen scan <target> [--ports=<port-range>]")
			return
		}

		target := args[0]
		ports, _ := cmd.Flags().GetString("ports")

		fmt.Printf("扫描目标: %s\n", target)
		fmt.Printf("端口范围: %s\n", ports)

		// TODO: 实现端口扫描逻辑
		results, err := scanner.ScanPorts(target, ports)
		if err != nil {
			fmt.Printf("扫描出错: %v\n", err)
			return
		}

		fmt.Println("扫描结果:")
		for _, result := range results {
			fmt.Printf("  端口 %d: %s\n", result.Port, result.Status)
		}
	},
}

func init() {
	scannerCmd.Flags().String("ports", "1-1000", "要扫描的端口范围，如 '1-1000' 或 '80,443,8080'")
}
