package app

import (
	"fmt"
	"github.com/clinaresl/table"
	"strings"
)

func ConstructTSVOutput(repos []Repo) string {
	output := "Name\tPath\tAuthor\tLastModified\tSyncedWithRemote\tSyncDescription\n"
	for _, repo := range repos {
		if repo.Name == "" {
			continue
		}
		year, month, day := repo.LastModified.Date()
		lastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		oneLineDescription := strings.Join(
			strings.Split(repo.SyncDescription, "\n"), " ",
		)
		row := fmt.Sprintf("%s\t%s\t%s\t%s\t%t\t%s\n", repo.Name, repo.AbsPath, repo.Author, lastModifiedDate, repo.SyncedWithRemote, oneLineDescription)
		output += row
	}
	return output
}

func ConstructSummary(repos []Repo, root string) string {
	countRepos := len(repos)
	var countUnsynced int
	for _, repo := range repos {
		if repo.Name == "" {
			continue
		}
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
	t, err := table.NewTable("| C{15} | L{20} | C{12} | c | c | L{27} |")
	if err != nil {
		return nil, err
	}
	t.AddThickRule()
	t.AddRow("Repo", "Path", "Author", "Last Modified", "Synced?", "Sync Details")
	t.AddThickRule()
	for i, repo := range repos {
		if repo.Name == "" {
			continue
		}
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
			repos[i].Author,
			LastModifiedDate,
			isSynced,
			repos[i].SyncDescription,
		)
		t.AddSingleRule()
	}
	return t, nil

}
