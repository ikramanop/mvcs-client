package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/utility"
	"github.com/spf13/cobra"
)

const (
	ProjectAlreadyExistsOnServer = "проект уже есть на remote"
	ProjectAlreadyExistsLocally  = "проект уже есть на local"
	ProjectCreateProcessError    = "ошибка в процессе создания проекта"
	BranchAlreadyExistsOnServer  = "ветка уже есть на remote"
	BranchAlreadyExistsLocally   = "ветка уже есть на local"
	BranchCreateProcessError     = "ошибка в процессе создания ветки"
)

type Project struct {
	Identifier string `protobuf:"bytes,1,opt,name=identifier,proto3" json:"identifier,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

var ProjectCreateCMD = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return create()
	},
}

func create() error {
	c, err := app.NewContainer(true, false, false)
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

	ctx, cancel := c.GetContext(false)
	defer cancel()

	// чекаем на сервере наличие проекта с таким же именем
	serverProject, err := c.Projects.Find(ctx, wd)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}
	if serverProject != nil {
		return model.WrapError(ProjectAlreadyExistsOnServer, err)
	}

	// создаем локально
	if err = c.DB.CreateProject(&model.ProjectParams{Name: wd}); err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	dbProject, err := c.DB.GetProjectByName(&model.ProjectParams{Name: wd})
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирован проект %s в директории %s на local\n", dbProject.Name, wd)

	// создаем на сервере
	serverProject, err = c.Projects.Create(ctx, dbProject.Name)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирован проект %s в директории %s на remote\n", dbProject.Name, wd)

	// создаем ветку локально
	if err = c.DB.CreateBranch(&model.BranchParams{
		Name:           "master",
		ProjectId:      dbProject.Id,
		UserIdentifier: c.Cfg.Auth.Login,
		ParentBranchId: sql.NullInt32{},
	}); err != nil {
		return model.WrapError(BranchCreateProcessError, err)
	}

	dbBranch, err := c.DB.GetBranchByName(&model.BranchParams{Name: "master"})
	if err != nil {
		return model.WrapError(BranchCreateProcessError, err)
	}

	fmt.Printf("Инициализирована ветка %s на local\n", dbBranch.Name)

	// создаем на сервере
	_, err = c.Project.CreateBranch(ctx, dbBranch.Name, nil)
	if err != nil {
		return model.WrapError(BranchCreateProcessError, err)
	}

	fmt.Printf("Инициализирована ветка %s на remote\n", dbProject.Name)

	//todo дальше отправить запросы на создание версий веток на сервере

	fileStructure, err := c.Diff.GetFileStructureWithHashesAndInfo()
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	for file, scam := range fileStructure {
		if err := c.DB.CreateFile(&model.FileParams{
			BranchId: dbBranch.Id,
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
		BranchId:       dbBranch.Id,
		Uuid:           uuid.New().String(),
		FilesStructure: string(extractedTreeJson),
	}); err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Println("Версии файлов зафиксированы для ветки master")

	return nil
}
