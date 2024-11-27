package app

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// helper function to streamline error checks for unhandled errors
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
	valid            bool
}

func GetReposWithDetails(root string) ([]Repo, error) {
	fsys := os.DirFS(root)
	repoPaths, err := listRepoPaths(fsys)
	repos := make([]Repo, len(repoPaths))
	if err != nil {
		return nil, err
	}
	for i, path := range repoPaths {
		dirFS, err := fs.Sub(fsys, path)
		check(err)
		lastModified := getContentLastModifiedTime(dirFS)
		absPath := filepath.Join(root, path)
		syncedWithRemote, syncDescription, err := getSyncStatus(absPath)
		repos[i] = Repo{
			Name:             filepath.Base(path),
			Path:             path,
			AbsPath:          absPath,
			LastModified:     lastModified,
			SyncedWithRemote: syncedWithRemote,
			SyncDescription:  syncDescription,
			valid:            true,
		}
		// error from getSyncStatus most likely means the directory
		// is not a git repository, hence set valid to false to mark
		// it as a non-repo
		if err != nil {
			repos[i].valid = false
		}
	}
	// clean up invalid repos
	repos = slices.DeleteFunc(repos, func(repo Repo) bool {
		return !repo.valid
	})
	return repos, nil
}

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
		check(err)
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

func getSyncStatus(absPath string) (bool, string, error) {
	var statusDescription string
	// returns "" for repos that have all changes committed
	cmdCommitStatus := exec.Command("git", "status", "-s")
	cmdCommitStatus.Dir = absPath
	out, err := cmdCommitStatus.CombinedOutput()
	if err != nil {
		notARepo := strings.Contains(string(out), "not a git repository")
		if notARepo {
			return false, "", err
		}
		return false, fmt.Sprintf("Error: %v", string(out)), nil
	}
	allChangesCommitted, commitStatusDescription := evaluateCommitSyncStatus(string(out))
	statusDescription += commitStatusDescription

	// this command will return an output where each line will contain
	// the status of the branch: "=" - synced, ">" - ahead, "<" - behind
	// "" - no remote branch
	cmdBranchStatus := exec.Command("git", "for-each-ref", "--format=%(upstream:trackshort)", "refs/heads")
	cmdBranchStatus.Dir = absPath
	out, err = cmdBranchStatus.CombinedOutput()
	if err != nil {
		return false, fmt.Sprintf("Error: %v", string(out)), nil
	}
	allBranchesSynced, branchStatusDescription := evaluateBranchSyncStatus(string(out))
	statusDescription += branchStatusDescription

	statusDescription = strings.TrimSuffix(string(statusDescription), "\n")
	syncedWithRemote := allBranchesSynced && allChangesCommitted
	return syncedWithRemote, statusDescription, nil
}

func evaluateCommitSyncStatus(gitOut string) (bool, string) {
	if gitOut == "" {
		return true, ""
	} else {
		return false, "- has uncommitted changes\n"
	}

}

func evaluateBranchSyncStatus(gitOut string) (bool, string) {
	var statusDescription string
	branchesNoRemote := false
	branchesAhead := false
	branchesBehind := false

	// remove the trailing new line from the output because otherwise
	// it will result in a "" in slice when the output is split on \n
	// leading to "" being evaluated as a branch as well later
	// remove trailing new line in description for prettier table printing
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
		statusDescription += "- has branch(es) with no remote branch\n"
	}
	if branchesAhead {
		statusDescription += "- has branch(es) that are ahead\n"
	}
	if branchesBehind {
		statusDescription += "- has branch(es) that are behind\n"
	}
	allBranchesSynced := !branchesNoRemote && !branchesAhead && !branchesBehind
	return allBranchesSynced, statusDescription
}
