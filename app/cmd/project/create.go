package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ikramanop/mvcs-client/app/utility"

	"github.com/spf13/cobra"

	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
)

const ProjectCreateProcessError = "ошибка в процессе создания проекта"

var ProjectCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return create()
	},
}

func create() error {
	c, err := app.NewContainer(true, false)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	wd, err := utility.GetWorkingDirectory()

	created, err := c.DB.CheckRepositoryExistence()
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	if created {
		_, err := c.DB.GetProjectByName(&model.ProjectParams{Name: wd})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return model.WrapError(ProjectCreateProcessError, err)
		}
		if err == nil {
			fmt.Printf("Проект %s уже инициализирован\n", wd)

			return nil
		}
	} else {
		if err = c.DB.Init(); err != nil {
			return model.WrapError(ProjectCreateProcessError, err)
		}
	}

	// тут надо проверить, что на сервере есть проект
	// если есть, то нужно закрывать процесс, тк мы не можем создавать два проекта с одним именем
	// если нет, то продолжаем процесс

	if err = c.DB.CreateProject(&model.ProjectParams{Name: wd}); err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	project, err := c.DB.GetProjectByName(&model.ProjectParams{Name: wd})
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирован проект %s в директории %s\n", project.Name, wd)

	if err = c.DB.CreateBranch(&model.BranchParams{
		Name:           "master",
		ProjectId:      project.Id,
		UserIdentifier: c.Cfg.Auth.Login,
		ParentBranchId: sql.NullInt32{},
	}); err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	branch, err := c.DB.GetBranchByName(&model.BranchParams{Name: "master"})
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирована ветка %s\n", branch.Name)

	fileStructure, err := c.Diff.GetFileStructureWithHashesAndInfo()
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	for file, scam := range fileStructure {
		if err := c.DB.CreateFile(&model.FileParams{
			BranchId: branch.Id,
			FilePath: file,
		}); err != nil {
			return model.WrapError(ProjectCreateProcessError, err)
		}

		fileDB, err := c.DB.GetFileByPath(&model.FileParams{FilePath: file})
		if err != nil {
			return model.WrapError(ProjectCreateProcessError, err)
		}

		content, err := utility.GetFileContent(file)
		if err != nil {
			return model.WrapError(ProjectCreateProcessError, err)
		}

		if err := c.DB.CreateFileVersion(&model.FileVersionParams{
			FileId:  fileDB.Id,
			Hash:    string(scam.Hash),
			Content: content,
		}); err != nil {
			return model.WrapError(ProjectCreateProcessError, err)
		}
	}

	extractedTree := c.Diff.ExtractFileTree(fileStructure)

	extractedTreeJson, err := json.Marshal(extractedTree)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	if err := c.DB.CreateBranchVersion(&model.BranchVersionParams{
		BranchId:       branch.Id,
		Uuid:           uuid.New().String(),
		FilesStructure: string(extractedTreeJson),
	}); err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Println("Версии файлов зафиксированы для ветки master")

	// тут нужно добавить отправку файлов на сервер
	// методы тут https://github.com/Amir0715/MVCS-Server/blob/main/MVCS/MVCS.Presentation.gRPC/Protos/projects.proto
	// нужно для каждого написать оболочку в client, прокинуть в container и вызывать здесь обёрнутые методы

	// НУЖНО ДОБАВИТЬ ОТПРАВКУ ДАННЫХ НА СЕРВЕР, ЛОКАЛЬНО ВСЁ ГОТОВО

	return nil
}
