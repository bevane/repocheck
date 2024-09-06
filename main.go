package main

import (
	"fmt"
	"os"

	"github.com/bevane/rpchk/app"
)

func main() {
	root := "."
	fsys := os.DirFS(root)
	repoDirectories := app.ListRepoDirectories(fsys)
	for _, repo := range repoDirectories {
		fmt.Printf("%v %v %v \n", repo.Name, repo.Path, repo.LastModified.String())
	}
}
