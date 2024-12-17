package app

import (
	"encoding/json"
	"fmt"
	"github.com/clinaresl/table"
	"strings"
)

func ConstructJSONOutput(repos []Repo) string {
	jsonOutput, _ := json.MarshalIndent(&repos, "", "\t")
	// add new line at the end because Marshal does not end the output with newline
	return string(jsonOutput) + "\n"
}

func ConstructTSVOutput(repos []Repo) string {
	output := "Name\tPath\tAuthor\tLastModified\tSynced\tSyncDetails\n"
	for _, repo := range repos {
		// skip any invalid repos
		if repo.Name == "" {
			continue
		}
		year, month, day := repo.LastModified.Date()
		lastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		row := fmt.Sprintf("%s\t%s\t%s\t%s\t%t\t%s\n", repo.Name, repo.AbsPath, repo.Author, lastModifiedDate, repo.SyncedWithRemote, strings.Join(repo.SyncDetails, ", "))
		output += row
	}
	return output
}

func ConstructSummary(repos []Repo, root string) string {
	countRepos := len(repos)
	var countUnsynced int
	for _, repo := range repos {
		// skip any invalid repos
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
	t, err := table.NewTable("| C{15} | L{20} | L{10} | c | c | L{23} |")
	if err != nil {
		return nil, err
	}
	t.AddThickRule()
	t.AddRow("Repo", "Path", "Author", "Last Modified", "Synced", "Sync Details")
	t.AddThickRule()
	for i, repo := range repos {
		// skip any invalid repos
		if repo.Name == "" {
			continue
		}
		year, month, day := repos[i].LastModified.Date()
		LastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		prettySyncDetails := ""
		// format sync details so that each detail is in its own line
		for _, line := range repo.SyncDetails {
			prettySyncDetails += "- " + line + "\n"
		}

		t.AddRow(
			repos[i].Name,
			repos[i].AbsPath,
			repos[i].Author,
			LastModifiedDate,
			repos[i].SyncedWithRemote,
			prettySyncDetails,
		)
		t.AddSingleRule()
	}
	return t, nil

}
