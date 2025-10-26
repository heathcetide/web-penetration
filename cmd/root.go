package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "web-pen",
	Short: "一个专业的Web渗透测试工具集",
	Long: `Web Penetration Testing Tool - 专业的Web渗透测试工具集合
包含多种常见的Web安全测试功能，用于安全审计和漏洞检测。
请仅将此工具用于合法的安全测试目的。`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(scannerCmd)
	rootCmd.AddCommand(fuzzCmd)
	rootCmd.AddCommand(bruteForceCmd)
	rootCmd.AddCommand(sqliCmd)
	rootCmd.AddCommand(xssCmd)
	rootCmd.AddCommand(fileScanCmd)
}
