package repository

import "errors"

func (r *ProjectsRepository) CheckRepositoryExistence() (bool, error) {
	row := r.sqlDB.QueryRow(
		`
SELECT count(1)
FROM sqlite_master
WHERE type = 'table'
  AND name in ('project', 'branch', 'branch_version', 'file', 'file_version');
`)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	if count < 5 {
		return false, errors.New("БД проекта сломана, проверьте схему")
	}

	return true, nil
}
