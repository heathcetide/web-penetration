package cmd

import (
	"fmt"
	"web-penetration/internal/filescan"

	"github.com/spf13/cobra"
)

var fileScanCmd = &cobra.Command{
	Use:   "filescan",
	Short: "敏感文件扫描",
	Long:  `扫描Web应用常见的敏感文件和备份文件。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标URL")
			fmt.Println("用法: web-pen filescan <url>")
			return
		}

		url := args[0]

		fmt.Printf("扫描目标: %s\n", url)

		// TODO: 实现文件扫描逻辑
		results, err := filescan.Scan(url)
		if err != nil {
			fmt.Printf("文件扫描出错: %v\n", err)
			return
		}

		fmt.Println("发现的敏感文件:")
		for _, file := range results {
			fmt.Printf("  %s (%s)\n", file.Path, file.Status)
		}
	},
}
