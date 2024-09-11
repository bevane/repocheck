package app

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// helper function to streamline error checks
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Repo struct {
	Name             string
	Path             string
	AbsPath          string
	LastModified     time.Time
	SyncedWithRemote bool
	SyncDescription  string
}

func CLI() int {
	flag.Parse()
	pathArg := flag.Arg(0)
	var root string
	var err error
	if pathArg == "" {
		root, err = os.Getwd()
		check(err)
	} else {
		root = pathArg
	}
	fsys := os.DirFS(root)
	repos := ListRepoDirectories(fsys)
	for _, repo := range repos {
		absPath := filepath.Join(root, repo.Path)
		repo.AbsPath = absPath
		syncedwithRemote, syncDescription := getSyncStatus(repo.AbsPath)

		fmt.Printf("%v %v %v %v %v \n", repo.Name, repo.AbsPath, repo.LastModified.String(), syncedwithRemote, syncDescription)
	}
	return 0
}

func getSyncStatus(absPath string) (bool, string) {
	var statusDescription string
	cmdCommitStatus := exec.Command("git", "status")
	cmdCommitStatus.Dir = absPath
	out, err := cmdCommitStatus.Output()
	check(err)
	allChangesCommited := strings.Contains(string(out), "nothing to commit")
	if !allChangesCommited {
		statusDescription += "- has uncommitted changes"

	}
	return allChangesCommited, statusDescription
}

func getContentLastModifiedTime(fileSystem fs.FS) time.Time {
	// returns the lastModified time of the most recently modified file/directory
	// in the given files system while ignoring the .git folder
	dirInfo, err := fs.Stat(fileSystem, ".")
	check(err)
	lastModified := dirInfo.ModTime()
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		check(err)
		// ignore .git folder's last modified date since it can change
		// when running git status even though the repo's contents have
		// not changed
		if d.Name() == ".git" {
			return fs.SkipDir
		}
		subDirInfo, err := d.Info()
		check(err)
		if subDirInfo.ModTime().Compare(lastModified) == 1 {
			lastModified = subDirInfo.ModTime()
		}

		return nil
	})
	return lastModified
}

func ListRepoDirectories(fileSystem fs.FS) []Repo {
	var repoDirectories []Repo
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		check(err)
		if !d.IsDir() {
			return nil
		}
		subDirs, err := fs.ReadDir(fileSystem, path)
		check(err)
		for _, subDir := range subDirs {
			if subDir.Name() == ".git" {
				dirFS, err := fs.Sub(fileSystem, path)
				check(err)
				lastModified := getContentLastModifiedTime(dirFS)
				repoDirectories = append(repoDirectories, Repo{
					Name:         d.Name(),
					Path:         path,
					LastModified: lastModified,
				})
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
