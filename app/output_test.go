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
	reposWithLongFields := []Repo{
		{
			Name:             "blink-frost-dune-glimmer",
			AbsPath:          "/home/repos/blink-frost-dune-glimmer",
			SyncedWithRemote: false,
			LastModified:     jan1,
			SyncDescription:  "- has uncommitted changes\n- has branch(es) with no remote branch",
		},
		{
			Name:             "stone-drift-moon-sparkle-breeze",
			AbsPath:          "/home/repos/stone-drift-moon-sparkle-breeze",
			SyncedWithRemote: false,
			LastModified:     jan2,
			SyncDescription:  "- has uncommitted changes\n- has branch(es) with no remote branch\n- has branch(es) that are ahead",
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
		`Name	Path	LastModified	SyncedWithRemote	SyncDescription
wheels	/home/repos/wheels	2024-01-01	true	
engine	/home/repos/engine	2024-01-02	true	
`
	outWithLongFields :=
		`Name	Path	LastModified	SyncedWithRemote	SyncDescription
blink-frost-dune-glimmer	/home/repos/blink-frost-dune-glimmer	2024-01-01	false	- has uncommitted changes - has branch(es) with no remote branch
stone-drift-moon-sparkle-breeze	/home/repos/stone-drift-moon-sparkle-breeze	2024-01-02	false	- has uncommitted changes - has branch(es) with no remote branch - has branch(es) that are ahead
`

	keyToOutputs := map[string]string{
		"short": outWithShortFields,
		"long":  outWithLongFields,
	}
	return keyToOutputs[key]
}
