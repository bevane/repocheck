package app

import (
	"fmt"
	"slices"
	"strings"
)

type SortFunc func([]Repo, bool)

func GetSortRepoFunc(key string) (SortFunc, error) {
	sortOptions := map[string]SortFunc{
		"name":         sortByRepoName,
		"path":         sortByRepoPath,
		"lastmodified": sortByRepoLastModified,
		"synced":       sortByRepoSyncStatus,
	}
	selectedSortFunc, ok := sortOptions[key]
	if !ok {
		var validOptions []string
		for key := range sortOptions {
			validOptions = append(validOptions, key)
		}
		// sort the keys to get a deterministic error message
		slices.SortFunc(validOptions, func(a, b string) int {
			return strings.Compare(a, b)
		})
		return nil, fmt.Errorf("%v is not a valid sort option. Options: %v", key, strings.Join(validOptions, " | "))
	}
	return selectedSortFunc, nil
}

func sortByRepoName(repos []Repo, reverse bool) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortByRepoPath(repos []Repo, reverse bool) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.AbsPath, b.AbsPath)
	})
}

func sortByRepoLastModified(repos []Repo, reverse bool) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return a.LastModified.Compare(b.LastModified)
	})
}

func sortByRepoSyncStatus(repos []Repo, reverse bool) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		// sort to show false values first as it is more likely the
		// user will want to see which repos are not synced when they
		// sort by sync status
		if a.SyncedWithRemote && !b.SyncedWithRemote {
			return 1
		} else if !a.SyncedWithRemote && b.SyncedWithRemote {
			return -1
		} else {
			return 0
		}
	})
}
