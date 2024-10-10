package app

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func ConstructPlainTable(repos []Repo) string {
	// var tableBuilder strings.Builder
	// w := tabwriter.NewWriter(&tableBuilder, 0, 0, 1, ' ', 0)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tPath\tLastModified\tSynced\tSyncDescription")
	for _, repo := range repos {
		year, month, day := repo.LastModified.Date()
		lastModifiedDate := fmt.Sprintf("%04d-%02d-%02d", year, int(month), day)
		var isSynced string
		if repo.SyncedWithRemote {
			isSynced = "yes"
		} else {
			isSynced = "no"
		}
		oneLineDescription := strings.Join(
			strings.Split(repo.SyncDescription, "\n"), " ",
		)
		row := fmt.Sprintf("%s\t%s\t%s\t%s\t%s", repo.Name, repo.AbsPath, lastModifiedDate, isSynced, oneLineDescription)
		fmt.Fprintln(w, row)
	}
	w.Flush()
	// return tableBuilder.String()
	return ""
}

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
