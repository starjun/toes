package apiserver

import (
	"fmt"

	"github.com/spf13/cobra"

	"toes/global"
)

func NewAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "toes",
		Short:        "A good Go api-server project",
		Long:         `A good Go api-server project, by st.......`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 初始化日志
			global.InitLog(&global.Cfg.Log)
			defer global.LogSync()
			// Sync 将缓存中的日志刷新到磁盘文件中

			return Run()
			// return internal.TestRun()
			// return internal.Runendless()

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
		"The path to the configuration file. Default: ./conf/apiserver.yaml")

	return cmd
}

func main() {
	command := NewAppCommand()
	cobra.CheckErr(command.Execute())
}
