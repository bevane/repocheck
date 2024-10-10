package app

import (
	"fmt"
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
