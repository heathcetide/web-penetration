package cmd

import (
	"fmt"
	"web-penetration/internal/xss"

	"github.com/spf13/cobra"
)

var xssCmd = &cobra.Command{
	Use:   "xss",
	Short: "XSS漏洞检测",
	Long:  `检测目标URL是否存在跨站脚本(XSS)漏洞。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标URL")
			fmt.Println("用法: web-pen xss <url> [--parameter=<param>]")
			return
		}

		url := args[0]
		parameter, _ := cmd.Flags().GetString("parameter")

		fmt.Printf("目标URL: %s\n", url)
		fmt.Printf("测试参数: %s\n", parameter)

		// TODO: 实现XSS检测逻辑
		result, err := xss.Scan(url, parameter)
		if err != nil {
			fmt.Printf("XSS检测出错: %v\n", err)
			return
		}

		if result.Vulnerable {
			fmt.Printf("发现XSS漏洞! 类型: %s\n", result.Type)
			fmt.Printf("Payload: %s\n", result.Payload)
		} else {
			fmt.Println("未检测到XSS漏洞")
		}
	},
}

func init() {
	xssCmd.Flags().String("parameter", "search", "要测试的参数名")
}
