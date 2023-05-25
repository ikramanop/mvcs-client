package cmd

import (
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
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

	RenderProjectsList(projectsList)

	return nil
}

func RenderProjectsList(projectsList *messages.ProjectList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", " Название"})
	for _, project := range projectsList.Projects {
		t.AppendRow([]interface{}{project.Identifier, project.Name})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
