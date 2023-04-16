package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/ikramanop/mvcs-client/app/utility"
)

const ConfigFilename = `config`
const SecretFilename = `secret`

var PathPrefix = func() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(dir, ".mvcs")
}()

type ClientConfig struct {
	ServerHost     string        `json:"server_host"`
	Auth           Auth          `json:"auth"`
	RequestTimeout time.Duration `json:"request_timeout"`
}

type Auth struct {
	Login string `json:"login"`
	Key   string `json:"key"`
}

func NewConfigFromOS() (*ClientConfig, error) {
	config := ClientConfig{RequestTimeout: 10}
	path := filepath.Join(PathPrefix, ConfigFilename)

	if _, err := os.Stat(PathPrefix); os.IsNotExist(err) {
		err := os.Mkdir(PathPrefix, 0777)

		return nil, err
	}

	content, err := utility.GetFileContent(path)
	if err != nil {
		marshaled, err := json.Marshal(&config)
		if err != nil {
			return nil, err
		}

		err = utility.CreateAndWriteFile(path, marshaled)
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, errors.New("ошибка получения конфигурации приложения: " + err.Error())
	}

	return &config, nil
}
