package difftool

import "os"

type Tool interface {
	GetFileStructureWithHashesAndInfo() (map[string]FileScam, error)
	ExtractFileTree(tree map[string]FileScam) map[string]string
}

type FileScam struct {
	Hash []byte
	Info os.FileInfo
}

type DiffTool struct {
}

func NewDiffTool() (*DiffTool, error) {
	tool := DiffTool{}

	return &tool, nil
}
