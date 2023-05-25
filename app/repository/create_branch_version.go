package repository

import (
	"github.com/ikramanop/mvcs-client/app/model"
)

func (r *ProjectsRepository) CreateBranchVersion(params *model.BranchVersionParams) error {
	_, err := r.sqlDB.Exec(
		`
insert into branch_version (branch_id, uuid, files_structure) values (?, ?, ?);
`, params.BranchId, params.Uuid, params.FilesStructure)
	if err != nil {
		return err
	}

	return nil
}
