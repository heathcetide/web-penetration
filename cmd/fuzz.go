package cmd

import (
	"fmt"
	"web-penetration/internal/fuzzer"

	"github.com/spf13/cobra"
)

var fuzzCmd = &cobra.Command{
	Use:   "fuzz",
	Short: "Web模糊测试功能",
	Long:  `对目标URL进行模糊测试，检测隐藏文件、目录和参数。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标URL")
			fmt.Println("用法: web-pen fuzz <url> [--wordlist=<file>]")
			return
		}

		url := args[0]
		wordlist, _ := cmd.Flags().GetString("wordlist")

		fmt.Printf("目标URL: %s\n", url)
		fmt.Printf("字典文件: %s\n", wordlist)

		// TODO: 实现模糊测试逻辑
		results, err := fuzzer.Fuzz(url, wordlist)
		if err != nil {
			fmt.Printf("模糊测试出错: %v\n", err)
			return
		}

		fmt.Println("发现的项目:")
		for _, result := range results {
			fmt.Printf("  %s -> %d\n", result.Path, result.StatusCode)
		}
	},
}

func init() {
	fuzzCmd.Flags().String("wordlist", "wordlist.txt", "字典文件路径")
}
