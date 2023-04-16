package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
)

const ProjectCreateProcessError = "ошибка в процессе создания проекта"

var ProjectCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return create()
	},
}

func create() error {
	_, err := app.NewContainer(true, false)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Println("TO DO later")

	return nil
}
