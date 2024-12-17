package app

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Repo struct {
	Name             string    `json:"name"`
	Path             string    `json:"-"`
	AbsPath          string    `json:"path"`
	LastModified     time.Time `json:"lastModified"`
	SyncedWithRemote bool      `json:"synced"`
	SyncDetails      []string  `json:"syncDetails"`
	Author           string    `json:"author"`
}

// recursively traverses all paths in 'root' and returns a slice of local git Repos
// that were found along with additional details for each repo
// the fetch argument determines if a git fetch is ran for each repo before
// getting rest of the repo details
func GetReposWithDetails(root string, fetch bool) ([]Repo, error) {
	var wg sync.WaitGroup
	fsys := os.DirFS(root)
	// gather all the repo paths first so that each repo can be concurrently
	// processed below
	repoPaths, err := listRepoPaths(fsys)
	// set to same length as repos so that goroutines can write to each
	// index of repos safely without needing to use append
	repos := make([]Repo, len(repoPaths))
	if err != nil {
		return nil, err
	}
	// concurrency is necessary because git fetch is a lengthy blocking call
	for i, path := range repoPaths {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			absPath := filepath.Join(root, path)
			dirFS, err := fs.Sub(fsys, path)
			if err != nil {
				slog.Warn(fmt.Sprintf("Unable to get the filesystem at %v, %v", absPath, err))
				return
			}
			if fetch {
				err = gitFetch(absPath)
				if err != nil {
					// continue without returning because git fetch can fail due to
					// network issues and the rest of the repo details can likely be
					// gathered
					slog.Warn(fmt.Sprintf("Unable to run git fetch at %v, %v", absPath, err))
				}
			}
			lastModified, err := getContentLastModifiedTime(dirFS)
			// continue without returning if lastmodified date could
			// not be calculated as it might still be possible for the the
			// directory to be a valid git repo
			if err != nil {
				slog.Warn(fmt.Sprintf("Unable get last modified time in %v, %v", absPath, err))
			}
			syncedWithRemote, syncDescription, err := getSyncStatus(absPath)
			if err != nil {
				// skip this directory as it is most likely not a valid git repo if git status and
				// git for-each-ref cannot be run
				slog.Warn(fmt.Sprintf("Unable to run git commands in %v, %v", absPath, err))
				return
			}
			author, err := getLastCommitAuthor(absPath)
			if err != nil {
				slog.Warn(fmt.Sprintf("Unable to get commit author in %v, %v", absPath, err))
			}
			repos[i] = Repo{
				Name:             filepath.Base(path),
				Path:             path,
				AbsPath:          absPath,
				LastModified:     lastModified,
				SyncedWithRemote: syncedWithRemote,
				SyncDetails:      syncDescription,
				Author:           author,
			}

		}(i, path)
	}
	wg.Wait()
	// clean up indexes that were not set to a Repo due to errors
	validRepos := []Repo{}
	for i := range repos {
		if repos[i].Name == "" {
			continue
		}
		validRepos = append(validRepos, repos[i])
	}
	return validRepos, nil
}

// returns all paths to directories in a fileSystem that contain a .git folder
func listRepoPaths(fileSystem fs.FS) ([]string, error) {
	var repoPaths []string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		subDirs, err := fs.ReadDir(fileSystem, path)
		if err != nil {
			return err
		}
		for _, subDir := range subDirs {
			// avoid counting dirs with .git file as a repo
			if !subDir.IsDir() {
				continue
			}
			if subDir.Name() == ".git" {
				repoPaths = append(repoPaths, path)
				// Prevent recursing through a repository directory
				// to improve performance as it is unlikely for another
				// repository to exist inside a repository
				return fs.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return repoPaths, nil
}

// run git fetch in a directory under absPath
func gitFetch(absPath string) error {
	cmdFetch := exec.Command("git", "fetch", "-q")
	cmdFetch.Dir = absPath
	out, err := cmdFetch.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

// returns the lastModified time of the most recently modified file/directory
// in the given fileSystem while ignoring the .git folder
func getContentLastModifiedTime(fileSystem fs.FS) (time.Time, error) {
	dirInfo, err := fs.Stat(fileSystem, ".")
	if err != nil {
		return time.Time{}, err
	}
	lastModified := dirInfo.ModTime()
	// recursively traverse through the directory and update the last modified
	// time if any file or folder is found with an even later last modified
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ignore .git folder's last modified date since it can change
		// when running git status even though the repo's contents have
		// not changed
		if d.Name() == ".git" {
			return fs.SkipDir
		}
		subDirInfo, err := d.Info()
		if err != nil {
			return err
		}
		if subDirInfo.ModTime().Compare(lastModified) == 1 {
			lastModified = subDirInfo.ModTime()
		}

		return nil
	})
	if err != nil {
		return time.Time{}, err
	}
	return lastModified, nil
}

// return a slice of strings describing whether the git repo at absPath
// has uncommitted changes, branches that are ahead/behind and untracked branches
func getSyncStatus(absPath string) (bool, []string, error) {
	// initialize as non-nil empty slice so that json output after marshalling will be []
	// instead of null
	statusDescription := []string{}
	// git status returns "" for repos that have all changes committed
	cmdCommitStatus := exec.Command("git", "status", "-s")
	cmdCommitStatus.Dir = absPath
	out, err := cmdCommitStatus.CombinedOutput()
	if err != nil {
		return false, nil, errors.New(string(out))
	}
	allChangesCommitted, commitStatusDescription := evaluateCommitSyncStatus(string(out))
	if commitStatusDescription != "" {
		statusDescription = append(statusDescription, commitStatusDescription)
	}

	// this command will return an output where each line will contain
	// the status of the branch: "=" - synced, ">" - ahead, "<" - behind
	// "" - no remote branch/untracked branch
	cmdBranchStatus := exec.Command("git", "for-each-ref", "--format=%(upstream:trackshort)", "refs/heads")
	cmdBranchStatus.Dir = absPath
	out, err = cmdBranchStatus.CombinedOutput()
	if err != nil {
		return false, nil, errors.New(string(out))
	}
	allBranchesSynced, branchStatusDescription := evaluateBranchSyncStatus(string(out))
	if branchStatusDescription != nil {
		statusDescription = append(statusDescription, branchStatusDescription...)
	}
	syncedWithRemote := allBranchesSynced && allChangesCommitted
	return syncedWithRemote, statusDescription, nil
}

// return the author name of the last commit
func getLastCommitAuthor(absPath string) (string, error) {
	cmdFetch := exec.Command("git", "log", "-1", "--pretty=%an")
	cmdFetch.Dir = absPath
	out, err := cmdFetch.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	author := strings.TrimSuffix(string(out), "\n")
	return author, nil
}

func evaluateCommitSyncStatus(gitOut string) (bool, string) {
	if gitOut == "" {
		return true, ""
	} else {
		return false, "uncommitted changes"
	}

}

func evaluateBranchSyncStatus(gitOut string) (bool, []string) {
	var statusDescription []string
	branchesNoRemote := false
	branchesAhead := false
	branchesBehind := false

	// remove the trailing new line from the output because otherwise
	// it will result in a "" in slice when the output is split on \n
	// leading to "" being evaluated as a branch as well later
	outWithoutEndingNewLine := strings.TrimSuffix(gitOut, "\n")
	branches := strings.Split(outWithoutEndingNewLine, "\n")
	for _, branch := range branches {
		switch branch {
		case "":
			branchesNoRemote = true
		case ">":
			branchesAhead = true
		case "<":
			branchesBehind = true
		}

	}
	if branchesNoRemote {
		statusDescription = append(statusDescription, "untracked branch(es)")
	}
	if branchesAhead {
		statusDescription = append(statusDescription, "branch(es) ahead")
	}
	if branchesBehind {
		statusDescription = append(statusDescription, "branch(es) behind")
	}
	allBranchesSynced := !branchesNoRemote && !branchesAhead && !branchesBehind
	return allBranchesSynced, statusDescription
}
