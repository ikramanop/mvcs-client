package repository

import (
	"github.com/ikramanop/mvcs-client/app/model"
)

func (r *ProjectsRepository) CreateFile(params *model.FileParams) error {
	_, err := r.sqlDB.Exec(
		`
insert into file (branch_id, file_path) values (?, ?);
`, params.BranchId, params.FilePath)
	if err != nil {
		return err
	}

	return nil
}
