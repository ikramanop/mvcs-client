package model

import "errors"

const (
	ConfigurationSaveError = "ошибка при сохранении конфигурации"
	ValidationError        = "ошибка при валидации"
	StdinError             = "ошибка при вводе"
)

func WrapError(baseError string, addition error) error {
	if addition == nil {
		return errors.New(baseError)
	}

	return errors.New(baseError + ": " + addition.Error())
}
