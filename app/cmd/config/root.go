package cmd

import (
	"github.com/spf13/cobra"
)

var ConfigCMD = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config()
	},
}

func init() {
	ConfigCMD.AddCommand(ConfigServerCMD)
	ConfigCMD.AddCommand(ConfigRegisterCMD)
	ConfigCMD.AddCommand(ConfigMeCMD)
}

func config() error {
	return nil
}
