package cmd

import (
	"fmt"
	"web-penetration/internal/sqli"

	"github.com/spf13/cobra"
)

var sqliCmd = &cobra.Command{
	Use:   "sqli",
	Short: "SQL注入检测",
	Long:  `检测目标URL是否存在SQL注入漏洞。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标URL")
			fmt.Println("用法: web-pen sqli <url> [--parameter=<param>]")
			return
		}

		url := args[0]
		parameter, _ := cmd.Flags().GetString("parameter")

		fmt.Printf("目标URL: %s\n", url)
		fmt.Printf("测试参数: %s\n", parameter)

		// TODO: 实现SQL注入检测逻辑
		result, err := sqli.Scan(url, parameter)
		if err != nil {
			fmt.Printf("SQL注入检测出错: %v\n", err)
			return
		}

		if result.Vulnerable {
			fmt.Printf("发现SQL注入漏洞! 类型: %s\n", result.Type)
		} else {
			fmt.Println("未检测到SQL注入漏洞")
		}
	},
}

func init() {
	sqliCmd.Flags().String("parameter", "id", "要测试的参数名")
}
