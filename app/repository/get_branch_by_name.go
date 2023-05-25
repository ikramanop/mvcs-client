package repository

import "github.com/ikramanop/mvcs-client/app/model"

func (r *ProjectsRepository) GetBranchByName(params *model.BranchParams) (*model.Branch, error) {
	row := r.sqlDB.QueryRow(
		`
select * from branch where name = ?;
`, params.Name)

	branch := model.Branch{}
	if err := row.Scan(
		&branch.Id,
		&branch.Name,
		&branch.ProjectId,
		&branch.UserIdentifier,
		&branch.ParentBranchId,
		&branch.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &branch, nil
}
