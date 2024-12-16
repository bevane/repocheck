package app

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var (
	jan1, _  = time.Parse(time.DateOnly, "2024-01-01")
	jan2, _  = time.Parse(time.DateOnly, "2024-01-02")
	jan3, _  = time.Parse(time.DateOnly, "2024-01-03")
	jan3a, _ = time.Parse(time.DateTime, "2024-01-03 10:00:00")
	jan4, _  = time.Parse(time.DateOnly, "2024-01-04")
)

var testQueries = NewQueries()

var sortTests = []struct {
	key  string
	want []Repo
}{
	{
		"name",
		getSortedOutput("name"),
	},
	{
		"path",
		getSortedOutput("path"),
	},
	{
		"lastmodified",
		getSortedOutput("lastmodified"),
	},
	{
		"synced",
		getSortedOutput("synced"),
	},
}

func TestSort(t *testing.T) {
	for _, test := range sortTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			testQueries.Sort.Value = test.key
			repos := getInputRepos()
			err := testQueries.Sort.apply(&repos)
			// the apply function  mutates the input hence the input itself is compared with want
			if !reflect.DeepEqual(repos, test.want) || err != nil {
				t.Errorf(
					"got (%v, %v)\nwant (%v, %v)",
					repos, err,
					test.want, nil,
				)
			}
		})
	}
}

func TestSortError(t *testing.T) {
	wantE := fmt.Errorf("invalid is not a valid sort option. Options: lastmodified | name | path | synced")
	testQueries.Sort.Value = "invalid"
	gotE := testQueries.Sort.validate()
	if gotE == nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"got (%v)\nwant (%v)",
			gotE, wantE,
		)
	}
}

var syncedFilterTests = []struct {
	key  string
	want []Repo
}{
	{
		"yes",
		getFilteredOutputSynced("yes"),
	},
	{
		"y",
		getFilteredOutputSynced("yes"),
	},
	{
		"no",
		getFilteredOutputSynced("no"),
	},
	{
		"n",
		getFilteredOutputSynced("no"),
	},
}

func TestSyncedFilter(t *testing.T) {
	for _, test := range syncedFilterTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			testQueries.Synced.Value = test.key
			repos := getInputRepos()
			err := testQueries.Synced.apply(&repos)
			// the apply function  mutates the input hence the input itself is compared with want
			if !reflect.DeepEqual(repos, test.want) || err != nil {
				t.Errorf(
					"got (%v, %v)\nwant (%v, %v)",
					repos, err,
					test.want, nil,
				)
			}
		})
	}
}

func TestSyncedFilterError(t *testing.T) {
	wantE := fmt.Errorf("incorrect value for synced, value must be either 'yes', 'y', 'no' or 'n'")
	testQueries.Synced.Value = "invalid"
	gotE := testQueries.Synced.validate()
	if gotE == nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"got (%v)\nwant (%v)",
			gotE, wantE,
		)
	}
}

var lastmodifiedFilterTests = []struct {
	key  string
	want []Repo
}{
	{
		"2024-01-03",
		getFilteredOutputLastModified("2024-01-03"),
	},
	{
		"<=2024-01-03",
		getFilteredOutputLastModified("<=2024-01-03"),
	},
	{
		">=2024-01-03",
		getFilteredOutputLastModified(">=2024-01-03"),
	},
	{
		"<2024-01-03",
		getFilteredOutputLastModified("<2024-01-03"),
	},
	{
		">2024-01-03",
		getFilteredOutputLastModified(">2024-01-03"),
	},
}

func TestLastModifiedFilter(t *testing.T) {
	for _, test := range lastmodifiedFilterTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			testQueries.LastModified.Value = test.key
			repos := getInputRepos()
			err := testQueries.LastModified.apply(&repos)
			// the apply function  mutates the input hence the input itself is compared with want
			if !reflect.DeepEqual(repos, test.want) || err != nil {
				t.Errorf(
					"got (%v, %v)\nwant (%v, %v)",
					repos, err,
					test.want, nil,
				)
			}
		})
	}
}

func TestLastModifiedFilterError(t *testing.T) {
	wantE := fmt.Errorf("unexpected date invalid, date must be in the format yyyy-mm-dd and can only be prefixed with '<=', '>=', '<' or '>'")
	testQueries.LastModified.Value = "invalid"
	gotE := testQueries.LastModified.validate()
	if gotE == nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"got (%v)\nwant (%v)",
			gotE, wantE,
		)
	}
}

var authorFilterTests = []struct {
	key  string
	want []Repo
}{
	{
		"author ab",
		getFilteredOutputAuthor("author ab"),
	},
	{
		"author cd",
		getFilteredOutputAuthor("author cd"),
	},
	{
		"author e",
		getFilteredOutputAuthor("author e"),
	},
	{
		"author z",
		getFilteredOutputAuthor("author z"),
	},
}

func TestAuthorFilter(t *testing.T) {
	for _, test := range authorFilterTests {
		testname := fmt.Sprintf("%v", test.key)
		t.Run(testname, func(t *testing.T) {
			testQueries.Author.Value = test.key
			repos := getInputRepos()
			err := testQueries.Author.apply(&repos)
			// the apply function  mutates the input hence the input itself is compared with want
			if !reflect.DeepEqual(repos, test.want) || err != nil {
				t.Errorf(
					"got (%v, %v)\nwant (%v, %v)",
					repos, err,
					test.want, nil,
				)
			}
		})
	}
}

func TestValidateQueries(t *testing.T) {
	emptyQueries := NewQueries()
	validQueries := NewQueries()
	validQueries.Sort.Value = "name"
	validQueries.Synced.Value = "y"
	validQueries.LastModified.Value = ">=2024-01-01"
	emptyAndValidQueries := NewQueries()
	emptyAndValidQueries.Synced.Value = "no"

	var tests = []*queries{
		emptyQueries,
		validQueries,
		emptyAndValidQueries,
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%v", test)
		t.Run(testname, func(t *testing.T) {
			err := ValidateQueries(test)
			if err != nil {
				t.Errorf("got (%v)\nwant (%v)", err, nil)
			}
		})
	}
}

func TestValidateQueriesError(t *testing.T) {
	invalidQueries := NewQueries()
	invalidQueries.LastModified.Value = ">=2024-23-01"

	wantE := fmt.Errorf("unexpected date 2024-23-01, date must be in the format yyyy-mm-dd and can only be prefixed with '<=', '>=', '<' or '>'")
	gotE := ValidateQueries(invalidQueries)
	if gotE == nil || gotE.Error() != wantE.Error() {
		t.Errorf("got (%v)\nwant (%v)", gotE, wantE)
	}
}

func TestAppplyQueries(t *testing.T) {
	queries := NewQueries()
	queries.Synced.Value = "y"
	queries.LastModified.Value = ">=2024-01-02"
	queries.Sort.Value = "name"

	repos := getInputRepos()
	want := getApplyQueriesResult()

	err := ApplyQueries(queries, &repos)

	// ApplyQueries mutates the repos passed in
	if !reflect.DeepEqual(repos, want) || err != nil {
		t.Errorf(
			"got (%v, %v)\nwant (%v, %v)",
			repos, err,
			want, nil,
		)
	}
}

func getInputRepos() []Repo {
	// helper to provide fake input to test the sortFunc
	return []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
}

func getSortedOutput(key string) []Repo {
	// helper to provide expected outputs for each sortFunc
	sortedByName := []Repo{
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
	}
	sortedByAbsPath := []Repo{
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
	}
	sortedByLastModified := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
	}
	sortedBySynced := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
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

func getFilteredOutputSynced(key string) []Repo {
	// helper to provide expected outputs for each filter strategy apply
	filteredBySyncYes := []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
	}
	filteredBySyncNo := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
	outputOptions := map[string][]Repo{
		"yes": filteredBySyncYes,
		"no":  filteredBySyncNo,
	}
	return outputOptions[key]
}

func getFilteredOutputLastModified(key string) []Repo {
	filteredByLastModifiedEQjan3 := []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
	filteredByLastModifiedLEQjan3 := []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
	filteredByLastModifiedGEQjan3 := []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
	filteredByLastModifiedLESjan3 := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
	}
	filteredByLastModifiedGRTjan3 := []Repo{
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
	}
	outputOptions := map[string][]Repo{
		"2024-01-03":   filteredByLastModifiedEQjan3,
		"<=2024-01-03": filteredByLastModifiedLEQjan3,
		">=2024-01-03": filteredByLastModifiedGEQjan3,
		"<2024-01-03":  filteredByLastModifiedLESjan3,
		">2024-01-03":  filteredByLastModifiedGRTjan3,
	}
	return outputOptions[key]
}

func getFilteredOutputAuthor(key string) []Repo {
	filteredByAuthorAB := []Repo{
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			"author ab",
		},
	}
	filteredByAuthorCD := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			"author cd",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
	}
	filteredByAuthorE := []Repo{
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
	}
	var filteredByAuthorZ []Repo
	outputOptions := map[string][]Repo{
		"author ab": filteredByAuthorAB,
		"author cd": filteredByAuthorCD,
		"author e":  filteredByAuthorE,
		"author z":  filteredByAuthorZ,
	}
	return outputOptions[key]
}

func getApplyQueriesResult() []Repo {
	// Filtered by synced yes
	// Filtered by lastmodified >= jan2
	// Sorted by Name
	return []Repo{
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			"author ab",
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			"author cd",
		},
		{
			"e",
			"repos/e",
			"/home/user/repos/x/e",
			jan3a,
			true,
			"",
			"author e",
		},
	}
}
