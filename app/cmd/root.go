package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	config "github.com/ikramanop/mvcs-client/app/cmd/config"
	project "github.com/ikramanop/mvcs-client/app/cmd/project"
)

var rootCmd = &cobra.Command{
	Use:   "mvcs",
	Short: "Система контроля версий МИФИ",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(project.ProjectCMD)
	rootCmd.AddCommand(config.ConfigCMD)
}
