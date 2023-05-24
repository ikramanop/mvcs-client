package model

import "database/sql"

type ProjectParams struct {
	Name       string
	Identifier string
}

type BranchParams struct {
	Name           string
	ProjectId      int
	UserIdentifier string
	ParentBranchId sql.NullInt32
}

type FileParams struct {
	BranchId int
	FilePath string
}

type FileVersionParams struct {
	FileId  int
	Hash    string
	Content []byte
}

type BranchVersionParams struct {
	BranchId       int
	Uuid           string
	FilesStructure string
}
