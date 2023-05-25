package repository

import (
	"github.com/ikramanop/mvcs-client/app/model"
)

func (r *ProjectsRepository) CreateProject(params *model.ProjectParams) error {
	_, err := r.sqlDB.Exec(
		`
insert into project (name, identifier) values (?, ?);
`, params.Name, params.Identifier)
	if err != nil {
		return err
	}

	return nil
}
