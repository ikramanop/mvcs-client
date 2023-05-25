package repository

import "github.com/ikramanop/mvcs-client/app/model"

func (r *ProjectsRepository) GetFileByPath(params *model.FileParams) (*model.File, error) {
	row := r.sqlDB.QueryRow(
		`
select * from file where file_path = ?;
`, params.FilePath)

	file := model.File{}
	if err := row.Scan(
		&file.Id,
		&file.BranchId,
		&file.FilePath,
		&file.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &file, nil
}
