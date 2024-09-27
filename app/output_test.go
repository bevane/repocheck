package app

import (
	"testing"
)

func TestPlainOutput(t *testing.T) {
}

func getInputReposByKey(key string) []Repo {
	reposWithShortFields := []Repo{
		{
			Name:             "wheels",
			AbsPath:          "/home/repos/wheels",
			SyncedWithRemote: true,
			LastModified:     jan1,
		},
		{
			Name:             "engine",
			AbsPath:          "/home/repos/engine",
			SyncedWithRemote: true,
			LastModified:     jan2,
		},
	}
	keyToInputs := map[string][]Repo{
		"short": reposWithShortFields,
	}
	return keyToInputs[key]
}
