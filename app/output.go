package app

import (
	"fmt"
	"github.com/clinaresl/table"
	"strings"
)

func ConstructTSVOutput(repos []Repo) string {
	output := "Name\tPath\tLastModified\tSyncedWithRemote\tSyncDescription\n"
	for _, repo := range repos {
		year, month, day := repo.LastModified.Date()
		lastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		oneLineDescription := strings.Join(
			strings.Split(repo.SyncDescription, "\n"), " ",
		)
		row := fmt.Sprintf("%s\t%s\t%s\t%t\t%s\n", repo.Name, repo.AbsPath, lastModifiedDate, repo.SyncedWithRemote, oneLineDescription)
		output += row
	}
	return output
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
