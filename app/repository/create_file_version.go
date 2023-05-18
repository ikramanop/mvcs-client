package repository

import (
	"github.com/ikramanop/mvcs-client/app/model"
)

func (r *ProjectsRepository) CreateFileVersion(params *model.FileVersionParams) error {
	_, err := r.sqlDB.Exec(
		`
insert into file_version (file_id, hash, content) values (?, ?, ?);
`, params.FileId, params.Hash, params.Content)
	if err != nil {
		return err
	}

	return nil
}
