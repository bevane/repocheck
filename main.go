package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// helper function to streamline error checks
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	root := "/home/bevane/repos/"
	fileSystem := os.DirFS(root)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		check(err)
		if !d.IsDir() {
			return nil
		}
		subDirs, err := os.ReadDir(filepath.Join(root, path))
		check(err)
		for _, dir := range subDirs {
			if dir.Name() == ".git" {
				fmt.Println(path)
				// Prevent recursing through a repository directory
				// to improve performance as it is unlikely for another repository to exist
				// inside a repository
				return fs.SkipDir
			}
		}
		return nil
	})
}
