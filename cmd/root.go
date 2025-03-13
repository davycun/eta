package cmd

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/cmd/server"
	"github.com/davycun/eta/cmd/stats"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "eta",
	Short: "d",
	Long:  `eta is a very fast city element web app`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New("use help or -h to see usage")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func init() {
	rootCmd.AddCommand(server.StartCommand)
	rootCmd.AddCommand(stats.StartCommand)
}
func tip() {
	fmt.Printf("欢迎使用eta, 您可以通过%s 来查看帮助!\n", utils.FmtTextRed("eta -h"))
}

func Registry(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
