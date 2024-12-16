package app

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var TSVTests = []struct {
	key   string
	input []Repo
	want  string
}{
	{
		"short",
		getInputReposByKey("short"),
		getTSVOutputByKey("short"),
	},
	{
		"long",
		getInputReposByKey("long"),
		getTSVOutputByKey("long"),
	},
}

func TestTSVOutput(t *testing.T) {
	for _, test := range TSVTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			repos := test.input
			got := ConstructTSVOutput(repos)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("-want +got:\n%s", diff)
			}
		})
	}
}

var JSONTests = []struct {
	key   string
	input []Repo
	want  string
}{
	{
		"short",
		getInputReposByKey("short"),
		getJSONOutputByKey("short"),
	},
	{
		"long",
		getInputReposByKey("long"),
		getJSONOutputByKey("long"),
	},
}

func TestJSONOutput(t *testing.T) {
	for _, test := range JSONTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			repos := test.input
			got := ConstructJSONOutput(repos)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("-want +got:\n%s", diff)
			}
		})
	}
}

func getInputReposByKey(key string) []Repo {
	reposWithShortFields := []Repo{
		{
			Name:             "wheels",
			AbsPath:          "/home/repos/wheels",
			SyncedWithRemote: true,
			LastModified:     jan1,
			Author:           "Test Author",
		},
		{
			Name:             "engine",
			AbsPath:          "/home/repos/engine",
			SyncedWithRemote: true,
			LastModified:     jan2,
			Author:           "Test Author",
		},
	}
	reposWithLongFields := []Repo{
		{
			Name:             "blink-frost-dune-glimmer",
			AbsPath:          "/home/repos/blink-frost-dune-glimmer",
			SyncedWithRemote: false,
			LastModified:     jan1,
			SyncDescription:  "- uncommitted changes\n- untracked branch(es)",
			Author:           "Test Author",
		},
		{
			Name:             "stone-drift-moon-sparkle-breeze",
			AbsPath:          "/home/repos/stone-drift-moon-sparkle-breeze",
			SyncedWithRemote: false,
			LastModified:     jan2,
			SyncDescription:  "- uncommitted changes\n- untracked branch(es)\n- branch(es) ahead",
			Author:           "Test Author",
		},
	}
	keyToInputs := map[string][]Repo{
		"short": reposWithShortFields,
		"long":  reposWithLongFields,
	}
	return keyToInputs[key]
}

func getTSVOutputByKey(key string) string {
	outWithShortFields :=
		`Name	Path	Author	LastModified	Synced	SyncDetails
wheels	/home/repos/wheels	Test Author	2024-01-01	true	
engine	/home/repos/engine	Test Author	2024-01-02	true	
`
	outWithLongFields :=
		`Name	Path	Author	LastModified	Synced	SyncDetails
blink-frost-dune-glimmer	/home/repos/blink-frost-dune-glimmer	Test Author	2024-01-01	false	- uncommitted changes - untracked branch(es)
stone-drift-moon-sparkle-breeze	/home/repos/stone-drift-moon-sparkle-breeze	Test Author	2024-01-02	false	- uncommitted changes - untracked branch(es) - branch(es) ahead
`

	keyToOutputs := map[string]string{
		"short": outWithShortFields,
		"long":  outWithLongFields,
	}
	return keyToOutputs[key]
}

func getJSONOutputByKey(key string) string {
	outWithShortFields := `[
	{
		"name": "wheels",
		"path": "/home/repos/wheels",
		"lastModified": "2024-01-01T00:00:00Z",
		"synced": true,
		"syncDetails": "",
		"author": "Test Author"
	},
	{
		"name": "engine",
		"path": "/home/repos/engine",
		"lastModified": "2024-01-02T00:00:00Z",
		"synced": true,
		"syncDetails": "",
		"author": "Test Author"
	}
]
`
	outWithLongFields := `[
	{
		"name": "blink-frost-dune-glimmer",
		"path": "/home/repos/blink-frost-dune-glimmer",
		"lastModified": "2024-01-01T00:00:00Z",
		"synced": false,
		"syncDetails": "- uncommitted changes\n- untracked branch(es)",
		"author": "Test Author"
	},
	{
		"name": "stone-drift-moon-sparkle-breeze",
		"path": "/home/repos/stone-drift-moon-sparkle-breeze",
		"lastModified": "2024-01-02T00:00:00Z",
		"synced": false,
		"syncDetails": "- uncommitted changes\n- untracked branch(es)\n- branch(es) ahead",
		"author": "Test Author"
	}
]
`

	keyToOutputs := map[string]string{
		"short": outWithShortFields,
		"long":  outWithLongFields,
	}
	return keyToOutputs[key]
}
