package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"

	"github.com/ikramanop/mvcs-client/app"
	appConfig "github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/utility"
)

const (
	ConfigRegisterProcessError = "ошибка в процессе регистрации пользователя"
)

var ConfigRegisterCMD = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return register()
	},
}

func register() error {
	c, err := app.NewContainer(false, true)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, err)
	}

	reader := bufio.NewReader(os.Stdin)

	if c.Cfg.Auth.Login != "" || c.Cfg.Auth.Key != "" {
		for {
			fmt.Print("Авторизационные данные присутствуют. Повторить регистрацию? (y/n): ")
			confirmation, err := utility.ReadStdin(reader)
			if err != nil {
				return model.WrapError(ConfigRegisterProcessError, err)
			}

			if strings.ToLower(confirmation) == "y" {
				fmt.Println()

				break
			} else if strings.ToLower(confirmation) == "n" {
				return nil
			} else {
				fmt.Println("Введите 'y' или 'n'.")
			}
		}
	}

	fmt.Print("Введите желаемое имя пользователя: ")
	login, err := utility.ReadStdin(reader)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, err)
	}

	c.Cfg.Auth.Login = login

	key, err := utility.GenerateRandomString(256)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, err)
	}

	c.Cfg.Auth.Key = key

	ctx, cancel := c.GetContext(false)
	defer cancel()
	err = c.Auth.Register(ctx, login, key)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, err)
	}

	content, err := json.Marshal(c.Cfg)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, model.WrapError(model.ConfigurationSaveError, err))
	}

	err = utility.WriteFileContent(filepath.Join(appConfig.PathPrefix, appConfig.ConfigFilename), content)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, model.WrapError(model.ConfigurationSaveError, err))
	}

	fmt.Println("\nВаш ключ для доступа к серверу:", key)

	return nil
}
