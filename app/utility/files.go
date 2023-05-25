package utility

import (
	"os"
	"path/filepath"

	"github.com/ikramanop/mvcs-client/app/model"
)

const (
	OpenFileError         = "ошибка при открытии файла"
	CreateFileError       = "ошибка при создании файла"
	WriteFileError        = "ошибка при записи файла"
	WorkingDirectoryError = "ошибка при получении рабочей директории"
)

func GetWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", model.WrapError(WorkingDirectoryError, err)
	}

	relDir := filepath.Base(dir)

	return relDir, nil
}

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
