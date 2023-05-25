package cmd

import (
	"github.com/spf13/cobra"
)

var ProjectCMD = &cobra.Command{
	Use:   "project",
	Short: "",
	Long:  "",
}

func init() {
	ProjectCMD.AddCommand(ProjectCreateCMD)
	ProjectCMD.AddCommand(GetBranchesCMD)
}
