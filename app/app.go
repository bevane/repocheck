package app

import (
	"io/fs"
)

// helper function to streamline error checks
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ListRepoDirectories(fileSystem fs.FS) []string {
	// returns a slice of relative paths for each repo directory found
	var repoDirectories []string
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		check(err)
		if !d.IsDir() {
			return nil
		}
		subDirs, err := fs.ReadDir(fileSystem, path)
		check(err)
		for _, dir := range subDirs {
			if dir.Name() == ".git" {
				repoDirectories = append(repoDirectories, path)
				// Prevent recursing through a repository directory
				// to improve performance as it is unlikely for another
				// repository to exist inside a repository
				return fs.SkipDir
			}
		}
		return nil
	})
	return repoDirectories
}
