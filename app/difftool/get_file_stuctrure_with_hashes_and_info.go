package difftool

import (
	"crypto/sha256"
	"github.com/ikramanop/mvcs-client/app/repository"
	"io"
	"os"
	"path/filepath"
)

func (d *DiffTool) GetFileStructureWithHashesAndInfo() (map[string]FileScam, error) {
	dir := "."

	fileMap := make(map[string]FileScam)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && path != repository.DBName {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			hash := sha256.New()

			if _, err := io.Copy(hash, file); err != nil {
				return err
			}

			scam := FileScam{
				Hash: hash.Sum(nil),
				Info: info,
			}

			fileMap[path] = scam
		}

		return nil
	})

	return fileMap, err
}
