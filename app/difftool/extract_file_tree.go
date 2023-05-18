package difftool

func (d *DiffTool) ExtractFileTree(tree map[string]FileScam) map[string]string {
	extractedTree := make(map[string]string)

	for file, scam := range tree {
		extractedTree[file] = string(scam.Hash)
	}

	return extractedTree
}
