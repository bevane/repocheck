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
	"time"
)

type Repo struct {
	Name             string
	Path             string
	AbsPath          string
	LastModified     time.Time
	SyncedWithRemote bool
	SyncDescription  string
}

func GetReposWithDetails(root string) ([]Repo, error) {
	fsys := os.DirFS(root)
	repoPaths, err := listRepoPaths(fsys)
	repos := make([]Repo, len(repoPaths))
	if err != nil {
		return nil, err
	}
	for i, path := range repoPaths {
		absPath := filepath.Join(root, path)
		dirFS, err := fs.Sub(fsys, path)
		if err != nil {
			slog.Warn(fmt.Sprintf("Unable to get the filesystem at %v, %v", absPath, err))
			continue
		}
		lastModified, err := getContentLastModifiedTime(dirFS)
		// dont skip and only log a warning if lastmodified date could
		// not be calculated as it might still be possible for the the
		// directory to be a valid git repo
		if err != nil {
			slog.Warn(fmt.Sprintf("Unable get last modified time in %v, %v", absPath, err))
		}
		// skip this directory as it is most likely not a git repo
		syncedWithRemote, syncDescription, err := getSyncStatus(absPath)
		if err != nil {
			slog.Warn(fmt.Sprintf("Unable to run git command in %v, %v", absPath, err))
			continue
		}
		repos[i] = Repo{
			Name:             filepath.Base(path),
			Path:             path,
			AbsPath:          absPath,
			LastModified:     lastModified,
			SyncedWithRemote: syncedWithRemote,
			SyncDescription:  syncDescription,
		}
	}
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

func getContentLastModifiedTime(fileSystem fs.FS) (time.Time, error) {
	// returns the lastModified time of the most recently modified file/directory
	// in the given files system while ignoring the .git folder
	dirInfo, err := fs.Stat(fileSystem, ".")
	if err != nil {
		return time.Time{}, err
	}
	lastModified := dirInfo.ModTime()
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

func getSyncStatus(absPath string) (bool, string, error) {
	var statusDescription string
	// returns "" for repos that have all changes committed
	cmdCommitStatus := exec.Command("git", "status", "-s")
	cmdCommitStatus.Dir = absPath
	out, err := cmdCommitStatus.CombinedOutput()
	if err != nil {
		return false, "", errors.New(string(out))
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
