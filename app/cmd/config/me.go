package cmd

import (
	"fmt"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/spf13/cobra"
)

const (
	ConfigCheckProcessError = "ошибка в процессе проверки пользователя"
)

var ConfigMeCMD = &cobra.Command{
	Use:   "me",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return me()
	},
}

func me() error {
	c, err := app.NewContainer(false, true, false)
	if err != nil {
		return model.WrapError(ConfigCheckProcessError, err)
	}

	if c.Cfg.Auth.Login == "" || c.Cfg.Auth.Key == "" {
		fmt.Println("Авторизационные данные отсутствуют")

		return nil
	}

	ctx, cancel := c.GetContext(true)
	defer cancel()
	message, err := c.Auth.Me(ctx)
	if err != nil {
		return model.WrapError(ConfigCheckProcessError, err)
	}

	fmt.Println("Логин пользователя:", message.Username)

	return nil
}
