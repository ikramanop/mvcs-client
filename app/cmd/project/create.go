package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ikramanop/mvcs-client/app"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/protobuf/messages"
	"github.com/ikramanop/mvcs-client/app/utility"
	"github.com/spf13/cobra"
)

const (
	CreateNewContainerError       = "ошибка при создании сессии"
	GetWorkingRepositoryError     = "ошибка при получении рабочей директории"
	CheckRepositoryExistenceError = "ошибка при получении репозитория"
	GetProjectError               = "ошибка при получении проекта"
	ProjectAlreadyExistsOnServer  = "проект уже есть на remote"
	ProjectAlreadyExistsLocally   = "проект уже есть на local"
	ProjectCreateProcessError     = "ошибка в процессе создания проекта"
	BranchAlreadyExistsOnServer   = "ветка уже есть на remote"
	BranchAlreadyExistsLocally    = "ветка уже есть на local"
	BranchCreateProcessError      = "ошибка в процессе создания ветки"
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
		return model.WrapError(CreateNewContainerError, err)
	}

	wd, err := GetWorkingDirectory(c)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}
	if wd == "" {
		return nil
	}

	ctx, cancel := c.GetContext(false)
	defer cancel()

	// получаем local и remote проекты
	project, err := CreateProject(c, ctx, wd)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	// создаем master ветки в local и remote проектах
	branch, err := CreateBranch(c, project, ctx)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	//todo дальше отправить запросы на создание версий веток на сервере

	err = AddFilesToLocalDB(c, branch)
	if err != nil {
		return model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Println("Версии файлов зафиксированы для ветки master")

	return nil
}

func AddFilesToLocalDB(c *app.Container, branch *model.Branch) error {
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

	return nil
}

func CreateBranch(c *app.Container, project *model.Project, ctx context.Context) (*model.Branch, error) {
	// создаем ветку локально
	localBranch, err := CreateLocalBranch(c, project)
	if err != nil {
		return nil, model.WrapError(BranchCreateProcessError, err)
	}

	fmt.Printf("Инициализирована ветка %s на local\n", localBranch.Name)

	// создаем на сервере
	remoteBranch, err := c.Project.CreateBranch(ctx, localBranch.Name, nil)
	if err != nil {
		return nil, model.WrapError(BranchCreateProcessError, err)
	}
	// тут наверное можно по айди, а не по имени
	if remoteBranch.Name != localBranch.Name {
		return nil, model.WrapError(BranchCreateProcessError, err)
	}

	fmt.Printf("Инициализирована ветка %s на remote\n", localBranch.Name)

	return localBranch, nil
}

func CreateLocalBranch(c *app.Container, project *model.Project) (*model.Branch, error) {
	if err := c.DB.CreateBranch(&model.BranchParams{
		Name:           "master",
		ProjectId:      project.Id,
		UserIdentifier: c.Cfg.Auth.Login,
		ParentBranchId: sql.NullInt32{},
	}); err != nil {
		return nil, model.WrapError(BranchCreateProcessError, err)
	}

	branch, err := c.DB.GetBranchByName(&model.BranchParams{Name: "master"})
	if err != nil {
		return nil, model.WrapError(BranchCreateProcessError, err)
	}

	return branch, err
}

func CreateProject(c *app.Container, ctx context.Context, wd string) (*model.Project, error) {
	// чекаем на сервере наличие проекта с таким же именем
	// и создаем на сервере, если там такого проекта нет
	remoteProject, err := CreateRemoteProject(c, ctx, wd)
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	// создаем локально
	localProject, err := CreateLocalProject(c, wd)
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	if remoteProject.Name != localProject.Name {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	return localProject, nil
}

func CreateLocalProject(c *app.Container, wd string) (*model.Project, error) {
	if err := c.DB.CreateProject(&model.ProjectParams{Name: wd}); err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	project, err := c.DB.GetProjectByName(&model.ProjectParams{Name: wd})
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирован проект %s в директории %s на local\n", project.Name, wd)

	return project, nil
}

func CreateRemoteProject(c *app.Container, ctx context.Context, wd string) (*messages.Project, error) {
	project, err := CheckRemoteProjectExists(c, ctx, wd)
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}
	if project != nil {
		fmt.Printf("Проект %s уже инициализирован в директории %s на remote\n", project.Name, wd)
		return project, nil
	}

	project, err = c.Projects.Create(ctx, wd)
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}

	fmt.Printf("Инициализирован проект %s в директории %s на local\n", project.Name, wd)

	return project, nil
}

func CheckRemoteProjectExists(c *app.Container, ctx context.Context, wd string) (*messages.Project, error) {
	serverProject, err := c.Projects.Find(ctx, wd)
	if err != nil {
		return nil, model.WrapError(ProjectCreateProcessError, err)
	}
	if serverProject != nil {
		return nil, nil
	}
	return serverProject, nil
}

func GetWorkingDirectory(c *app.Container) (string, error) {
	wd, err := utility.GetWorkingDirectory()
	if err != nil {
		return "", model.WrapError(GetWorkingRepositoryError, err)
	}

	repositoryExists, err := c.DB.CheckRepositoryExistence()
	if err != nil {
		return "", model.WrapError(CheckRepositoryExistenceError, err)
	}

	if repositoryExists {
		_, err := c.DB.GetProjectByName(&model.ProjectParams{Name: wd})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return "", model.WrapError(GetProjectError, err)
		}
		if err == nil {
			fmt.Printf("Проект %s уже инициализирован\n", wd)

			return "", nil
		}
	} else {
		if err = c.DB.Init(); err != nil {
			return "", model.WrapError(ProjectCreateProcessError, err)
		}
	}
	return wd, nil
}
