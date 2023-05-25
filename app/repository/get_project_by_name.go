package repository

import "github.com/ikramanop/mvcs-client/app/model"

func (r *ProjectsRepository) GetProjectByName(params *model.ProjectParams) (*model.Project, error) {
	row := r.sqlDB.QueryRow(
		`
select id, identifier, name, created_at from project where name = ?;
`, params.Name)

	project := model.Project{}
	if err := row.Scan(
		&project.Id,
		&project.Identifier,
		&project.Name,
		&project.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &project, nil
}
