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
	fmt.Println(repoDirectories)

}
