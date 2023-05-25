package cmd

import (
	"context"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/ikramanop/mvcs-client/app/utility"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

const (
	GetAllBranchesFromRemoteError = "ошибка при получении списка веток с сервера"
	GetAllBranchesError           = "ошибка при получении списка веток"
)

var GetBranchesCMD = &cobra.Command{
	Use:   "branches",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getBranches()
	},
}

func getBranches() error {
	c, err := app.NewContainer(true, false, false)
	if err != nil {
		return model.WrapError(CreateNewContainerError, err)
	}

	ctx, cancel := c.GetContext(true)
	defer cancel()

	wd, err := utility.GetWorkingDirectory()
	if err != nil {
		return model.WrapError(GetWorkingRepositoryError, err)
	}

	project, err := c.Projects.Find(ctx, wd)

	ctx, cancel = c.GetContext(true, project.Identifier)
	defer cancel()

	branches, err := GetAllBranches(c, ctx)
	if err != nil {
		return model.WrapError(GetAllBranchesError, err)
	}

	RenderBranches(branches)

	return nil
}

func GetAllBranches(c *app.Container, ctx context.Context) (*messages.GetBranchesResponse, error) {
	branches, err := c.Project.GetAllBranches(ctx)
	if err != nil {
		return nil, model.WrapError(GetAllBranchesFromRemoteError, err)
	}

	return branches, err
}

func RenderBranches(branches *messages.GetBranchesResponse) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", " Название"})
	for _, branch := range branches.Branches {
		t.AppendRow([]interface{}{branch.Id, branch.Name})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
