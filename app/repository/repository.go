package repository

import (
	"database/sql"
	"github.com/ikramanop/mvcs-client/app/model"
	"github.com/ikramanop/mvcs-client/app/utility"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const DBName = ".project"

const (
	DatabaseCreateError = "ошибка создания БД"
	DatabaseOpenError   = "ошибка открытия БД"
)

type Repository interface {
	Init() error
	CheckRepositoryExistence() (bool, error)
	CreateProject(params *model.ProjectParams) error
	GetProjectByName(params *model.ProjectParams) (*model.Project, error)
	CreateBranch(params *model.BranchParams) error
	GetBranchByName(params *model.BranchParams) (*model.Branch, error)
	CreateFile(params *model.FileParams) error
	GetFileByPath(params *model.FileParams) (*model.File, error)
	CreateFileVersion(params *model.FileVersionParams) error
	CreateBranchVersion(params *model.BranchVersionParams) error
}

type ProjectsRepository struct {
	sqlDB *sql.DB
}

func NewRepository(init bool) (*ProjectsRepository, error) {
	if init {
		_, err := os.Stat(DBName)
		if err != nil {
			err = utility.CreateAndWriteFile(DBName, nil)
			if err != nil {
				return nil, model.WrapError(DatabaseCreateError, err)
			}
		}
	}

	sqlDB, err := sql.Open("sqlite3", DBName)
	if err != nil {
		return nil, model.WrapError(DatabaseOpenError, err)
	}

	return &ProjectsRepository{
		sqlDB: sqlDB,
	}, err
}
