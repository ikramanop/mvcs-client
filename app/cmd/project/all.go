package cmd

import (
	"fmt"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/spf13/cobra"
)

const GetProjectsListFromRemote = "ошибка при получении списка веток с remote"

var GetProjectsListCMD = &cobra.Command{
	Use:   "all",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getProjectsList()
	},
}

func getProjectsList() error {
	c, err := app.NewContainer(true, false, false)
	if err != nil {
		return model.WrapError(CreateNewContainerError, err)
	}

	ctx, cancel := c.GetContext(true)
	defer cancel()

	projectsList, err := c.Projects.GetAll(ctx)
	if err != nil {
		return model.WrapError(GetProjectsListFromRemote, err)
	}

	fmt.Println(projectsList)

	return nil
}
