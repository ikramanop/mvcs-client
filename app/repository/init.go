package repository

func (r *ProjectsRepository) Init() error {
	_, err := r.sqlDB.Exec(
		`
create table if not exists project
(
    id         integer primary key autoincrement  not null,
    identifier string                             not null,
    name       text unique                        not null,
    created_at datetime default CURRENT_TIMESTAMP not null
);

create table if not exists branch
(
    id               integer primary key autoincrement  not null,
    name             text unique                        not null,
    project_id       integer                            not null,
    user_identifier  text                               not null,
    parent_branch_id integer                            null,
    created_at       datetime default CURRENT_TIMESTAMP not null,
    foreign key (project_id) references project (id)
);

create table if not exists branch_version
(
    id              integer primary key autoincrement  not null,
    branch_id       integer                            not null,
    uuid            text unique                        not null,
    files_structure json                               not null,
    created_at      datetime default CURRENT_TIMESTAMP not null,
    foreign key (branch_id) references branch (id)
);

create table if not exists file
(
    id         integer primary key autoincrement  not null,
    branch_id  integer                            not null,
    file_path  text                               not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    foreign key (branch_id) references branch (id)
);

create table if not exists file_version
(
    id         integer primary key autoincrement  not null,
    file_id    integer                            not null,
    hash       text                               not null,
    content    blob                               not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    foreign key (file_id) references file (id)
);
	`,
	)

	if err != nil {
		return err
	}

	return nil
}
