package config

import "github.com/spf13/cobra"

func BindArgConfig(cmd *cobra.Command, configFile *string, argConfig *Configuration) {
	if cmd == nil {
		return
	}
	cmd.Flags().StringVarP(configFile, "config", "c", "", "the config yaml file")
	cmd.Flags().StringVarP(&argConfig.Database.User, "user", "u", "", "database user")
	cmd.Flags().StringVarP(&argConfig.Database.Password, "password", "w", "", "database password")
	cmd.Flags().StringVarP(&argConfig.Database.DBName, "database", "d", "", "database name")
	cmd.Flags().StringVarP(&argConfig.Database.Schema, "schema", "e", "", "eta db schema")
	cmd.Flags().IntVarP(&argConfig.Server.Port, "srvport", "s", 0, "the port of web server")
	cmd.Flags().IntVarP(&argConfig.Monitor.Port, "mtport", "m", 0, "the port of web server")
}
