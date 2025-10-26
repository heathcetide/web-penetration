package cmd

import (
	"fmt"
	"web-penetration/internal/bruteforce"

	"github.com/spf13/cobra"
)

var bruteForceCmd = &cobra.Command{
	Use:   "brute",
	Short: "暴力破解功能",
	Long:  `对HTTP基础认证、表单登录等进行暴力破解测试。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("错误: 请指定目标URL")
			fmt.Println("用法: web-pen brute <url> --username=<user> --password-list=<file>")
			return
		}

		url := args[0]
		username, _ := cmd.Flags().GetString("username")
		passwordList, _ := cmd.Flags().GetString("password-list")

		fmt.Printf("目标URL: %s\n", url)
		fmt.Printf("用户名: %s\n", username)
		fmt.Printf("密码字典: %s\n", passwordList)

		// TODO: 实现暴力破解逻辑
		success, err := bruteforce.Attack(url, username, passwordList)
		if err != nil {
			fmt.Printf("暴力破解出错: %v\n", err)
			return
		}

		if success != nil {
			fmt.Printf("成功! 用户名: %s, 密码: %s\n", success.Username, success.Password)
		} else {
			fmt.Println("未找到有效凭据")
		}
	},
}

func init() {
	bruteForceCmd.Flags().String("username", "", "目标用户名")
	bruteForceCmd.Flags().String("password-list", "passwords.txt", "密码字典文件")
	bruteForceCmd.MarkFlagRequired("username")
}
