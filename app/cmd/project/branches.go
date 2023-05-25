package cmd

import (
	"context"
	"fmt"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/ikramanop/mvcs-client/app/utility"
	"github.com/spf13/cobra"
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

	fmt.Println(branches)

	return nil
}

func GetAllBranches(c *app.Container, ctx context.Context) (*messages.GetBranchesResponse, error) {
	branches, err := c.Project.GetAllBranches(ctx)
	if err != nil {
		return nil, model.WrapError(GetAllBranchesFromRemoteError, err)
	}

	return branches, err
}
