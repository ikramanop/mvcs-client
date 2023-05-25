package repository

import (
	"github.com/ikramanop/mvcs-client/app/model"
)

func (r *ProjectsRepository) CreateBranch(params *model.BranchParams) error {
	_, err := r.sqlDB.Exec(
		`
insert into branch (name, project_id, user_identifier, parent_branch_id) values (?, ?, ?, ?);
`, params.Name, params.ProjectId, params.UserIdentifier, params.ParentBranchId)
	if err != nil {
		return err
	}

	return nil
}
