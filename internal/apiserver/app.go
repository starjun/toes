// Package apiserver 提供 API 服务器的核心功能。
//
// 该包包含服务器启动、配置初始化、命令行接口等功能。
// 使用 Cobra 框架提供命令行界面，支持配置文件加载和日志初始化。
//
// 主要功能:
//   - 命令行参数解析
//   - 配置文件加载
//   - 日志系统初始化
//   - 服务器启动
//
// 使用示例:
//
//	cmd := apiserver.NewAppCommand()
//	cobra.CheckErr(command.Execute())
package apiserver

import (
	"fmt"
	"github.com/spf13/cobra"

	"toes/global"
)

func NewAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "toes",
		Short:        "A good Go api-apiserver project",
		Long:         `A good Go api-apiserver project, by st.......`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 初始化日志
			global.InitLog(&global.Cfg.Log)
			defer global.LogSync()
			// Sync 将缓存中的日志刷新到磁盘文件中

			// windows
			return Run()
			// return internal.TestRun()

			// default run
			//return Runendless()

		},
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// 以下设置，使得 InitConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(global.InitConfig)

	// 在这里您将定义标志和配置设置。

	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&global.CfgFile, "config", "c", "",
		"The path to the configuration file. Default: ./configs/apiserver.yaml")

	return cmd
}
