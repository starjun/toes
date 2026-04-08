// Package main 是 toes API 服务器的入口点。
//
// 该包提供命令行入口，使用 Cobra 框架解析命令行参数，
// 并启动 API 服务器。
//
// 使用示例:
//
//	go run cmd/apiserver/main.go -c configs/apiserver.yaml
//
// 参见:
//   - https://github.com/spf13/cobra
//   - https://godoc.org/github.com/spf13/cobra
package main

import (
	"github.com/spf13/cobra"

	"toes/internal/apiserver"
)

func main() {
	command := apiserver.NewAppCommand()
	cobra.CheckErr(command.Execute())
}
