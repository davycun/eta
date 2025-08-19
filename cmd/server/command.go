package server

import (
	"github.com/spf13/cobra"
)

var (
	StartCommand = &cobra.Command{
		Use:   "server",
		Short: "启动一个eta服务",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	return callLifeCycle()
}
