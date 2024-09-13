package app

import (
	"fmt"
	"github.com/clinaresl/table"
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

func GetRepos(root string) ([]Repo, error) {
	fsys := os.DirFS(root)
	repos, err := listRepoDirectories(fsys)
	if err != nil {
		return nil, err
	}
	for i := range repos {
		absPath := filepath.Join(root, repos[i].Path)
		repos[i].AbsPath = absPath
		repos[i].SyncedWithRemote, repos[i].SyncDescription, err = getSyncStatus(repos[i].AbsPath)
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

func ConstructSummary(repos []Repo, root string) string {
	countRepos := len(repos)
	var countUnsynced int
	for _, repo := range repos {
		if !repo.SyncedWithRemote {
			countUnsynced++
		}
	}
	return fmt.Sprintf(
		"%v repos found in %v: %v repo(s) are not synced",
		countRepos,
		root,
		countUnsynced,
	)
}

func ConstructTable(repos []Repo) (*table.Table, error) {
	t, err := table.NewTable("| C{15} | L{20} | c | c | L{27} |")
	if err != nil {
		return nil, err
	}
	t.AddThickRule()
	t.AddRow("Repo", "Path", "Last Modified", "Synced?", "Sync Details")
	t.AddThickRule()
	for i := range repos {
		year, month, day := repos[i].LastModified.Date()
		LastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		var isSynced string
		if repos[i].SyncedWithRemote {
			isSynced = "yes"
		} else {
			isSynced = "no"
		}

		t.AddRow(
			repos[i].Name,
			repos[i].AbsPath,
			LastModifiedDate,
			isSynced,
			repos[i].SyncDescription,
		)
		t.AddSingleRule()
	}
	return t, nil

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

func listRepoDirectories(fileSystem fs.FS) ([]Repo, error) {
	var repoDirectories []Repo
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
				return nil
			}
			if subDir.Name() == ".git" {
				dirFS, err := fs.Sub(fileSystem, path)
				check(err)
				lastModified := getContentLastModifiedTime(dirFS)
				repoDirectories = append(repoDirectories, Repo{
					Name:         d.Name(),
					Path:         path,
					LastModified: lastModified,
					valid:        true,
				})
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
	return repoDirectories, nil
}
