package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ikramanop/mvcs-client/app"
	appConfig "github.com/ikramanop/mvcs-client/app/config"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/utility"
)

const ConfigServerProcessError = "ошибка в процессе смены хоста сервера"

var ConfigServerCMD = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return server()
	},
}

func server() error {
	c, err := app.NewContainer(false, true, true)
	if err != nil {
		return model.WrapError(ConfigRegisterProcessError, err)
	}

	reader := bufio.NewReader(os.Stdin)

	if c.Cfg.ServerHost != "" {
		for {
			fmt.Print("Связь уже установлена для сервера: '", c.Cfg.ServerHost, "'. Изменить хост сервера? (y/n): ")
			confirmation, err := utility.ReadStdin(reader)
			if err != nil {
				return model.WrapError(ConfigServerProcessError, err)
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

	fmt.Print("Введите адрес сервера: ")
	input, err := utility.ReadStdin(reader)
	if err != nil {
		return model.WrapError(ConfigServerProcessError, err)
	}

	_, err = url.Parse(input)
	if err != nil {
		return model.WrapError(ConfigServerProcessError, model.WrapError(model.ValidationError, err))
	}

	c.Cfg.ServerHost = input

	content, err := json.Marshal(c.Cfg)
	if err != nil {
		return model.WrapError(ConfigServerProcessError, model.WrapError(model.ConfigurationSaveError, err))
	}

	err = utility.WriteFileContent(filepath.Join(appConfig.PathPrefix, appConfig.ConfigFilename), content)
	if err != nil {
		return model.WrapError(ConfigServerProcessError, model.WrapError(model.ConfigurationSaveError, err))
	}

	fmt.Printf("\nКонфигурация сервера: %+v\n", c.Cfg)

	return nil
}
