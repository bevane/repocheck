package app

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	Jan1, _ = time.Parse(time.DateOnly, "2024-01-01")
	Jan2, _ = time.Parse(time.DateOnly, "2024-01-02")
	Jan3, _ = time.Parse(time.DateOnly, "2024-01-03")
	Jan4, _ = time.Parse(time.DateOnly, "2024-01-04")
)

func TestGetSortRepoFuncValidKey(t *testing.T) {
	var tests = []struct {
		key   string
		want  []Repo
		wantE error
	}{
		{
			"name",
			getSortedOutputByKey("name"),
			nil,
		},
		{
			"path",
			getSortedOutputByKey("path"),
			nil,
		},
		{
			"lastmodified",
			getSortedOutputByKey("lastmodified"),
			nil,
		},
		{
			"synced",
			getSortedOutputByKey("synced"),
			nil,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.key)
		t.Run(testname, func(t *testing.T) {
			sortFunc, gotE := GetSortRepoFunc(tt.key)
			repos := getUnsortedRepos()
			// false refers to the reverse flag passed in which
			// is false by default
			sortFunc(repos, false)
			// the sortFunc itself needs to be tested to ensure that
			// GetSortRepoFunc returned the correct function
			// because there is no way to deterministically compare
			// return values that are functions
			if !reflect.DeepEqual(repos, tt.want) || gotE != tt.wantE {
				t.Errorf(
					"sortFunc results\ngot (%v, %v) , want (%v, %v)",
					repos, gotE,
					tt.want, tt.wantE,
				)
			}
		})
	}
}

func TestGetSortRepoFuncInvalidKey(t *testing.T) {
	// separate logic for testing invalid key because the sort func itself
	// will be nil for an invalid key and wont be tested
	wantE := fmt.Errorf("invalidkey is not a valid sort option. Options: lastmodified | name | path | synced")
	gotFunc, gotE := GetSortRepoFunc("invalidkey")
	if gotFunc != nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"GetSortRepoFunc results\ngot (%v, %v) , want (%v, %v)",
			gotFunc, gotE.Error(),
			nil, wantE.Error(),
		)
	}

}

func getUnsortedRepos() []Repo {
	// helper to provide fake input to test the sortFunc
	return []Repo{
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			Jan4,
			true,
			"",
			true,
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			Jan1,
			false,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			Jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			Jan3,
			false,
			"",
			true,
		},
	}
}

func getSortedOutputByKey(key string) []Repo {
	// helper to provide expected outputs for each sortFunc
	sortedByName := []Repo{
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			Jan3,
			false,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			Jan4,
			true,
			"",
			true,
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			Jan1,
			false,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			Jan2,
			true,
			"",
			true,
		},
	}
	sortedByAbsPath := []Repo{
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			Jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			Jan3,
			false,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			Jan4,
			true,
			"",
			true,
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			Jan1,
			false,
			"",
			true,
		},
	}
	sortedByLastModified := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			Jan1,
			false,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			Jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			Jan3,
			false,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			Jan4,
			true,
			"",
			true,
		},
	}
	sortedBySynced := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			Jan1,
			false,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			Jan3,
			false,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			Jan4,
			true,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			Jan2,
			true,
			"",
			true,
		},
	}
	outputOptions := map[string][]Repo{
		"name":         sortedByName,
		"path":         sortedByAbsPath,
		"lastmodified": sortedByLastModified,
		"synced":       sortedBySynced,
	}
	return outputOptions[key]
}
