package model

import (
	"database/sql"
	"time"
)

type Project struct {
	Id        int       `sql:"id"`
	Name      string    `sql:"name"`
	CreatedAt time.Time `sql:"created_at"`
}

type Branch struct {
	Id             int           `sql:"id"`
	Name           string        `sql:"name"`
	ProjectId      int           `sql:"project_id"`
	UserIdentifier string        `sql:"user_identifier"`
	ParentBranchId sql.NullInt32 `sql:"parent_branch_id"`
	CreatedAt      time.Time     `sql:"created_at"`
}

type File struct {
	Id        int       `sql:"id"`
	BranchId  int       `sql:"branch_id"`
	FilePath  string    `sql:"file_path"`
	CreatedAt time.Time `sql:"created_at"`
}
