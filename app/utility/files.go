package utility

import (
	"os"

	"github.com/ikramanop/mvcs-client/app/model"
)

const (
	OpenFileError   = "ошибка при открытии файла"
	CreateFileError = "ошибка при создании файла"
	WriteFileError  = "ошибка при записи файла"
)

func GetFileContent(filename string) ([]byte, error) {
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, model.WrapError(OpenFileError, err)
	}

	return body, nil
}

func CreateAndWriteFile(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return model.WrapError(CreateFileError, err)
	}

	if content != nil {
		_, err = file.Write(content)
		if err != nil {
			return model.WrapError(WriteFileError, err)
		}
	}

	return nil
}

func WriteFileContent(filename string, content []byte) error {
	err := os.WriteFile(filename, content, 0777)
	if err != nil {
		return model.WrapError(WriteFileError, err)
	}

	return nil
}
