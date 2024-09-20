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

var getFilterStrategyTests = []struct {
	query string
	want  reflect.Type
}{
	{
		"synced=yes",
		reflect.TypeOf(SyncedEQFilterStrategy{}),
	},
	{
		"synced=no",
		reflect.TypeOf(SyncedEQFilterStrategy{}),
	},
	{
		"lastmodified=2024-01-03",
		reflect.TypeOf(LastModifiedEQFilterStrategy{}),
	},
	{
		"lastmodified<=2024-01-03",
		reflect.TypeOf(LastModifiedLEQFilterStrategy{}),
	},
	{
		"lastmodified>=2024-01-03",
		reflect.TypeOf(LastModifiedGEQFilterStrategy{}),
	},
	{
		"lastmodified<2024-01-03",
		reflect.TypeOf(LastModifiedLESFilterStrategy{}),
	},
	{
		"lastmodified>2024-01-03",
		reflect.TypeOf(LastModifiedGRTFilterStrategy{}),
	},
	{
		"invalidQuery",
		nil,
	},
}

func TestGetFilterStrategy(t *testing.T) {
	for _, test := range getFilterStrategyTests {
		filterStrategy, _ := GetFilterStrategy(test.query)
		if reflect.TypeOf(filterStrategy) != test.want {
			t.Errorf(
				"got (%v) , want (%v)",
				reflect.TypeOf(filterStrategy),
				test.want,
			)
		}
	}
}

func TestGetFilterStrategyError(t *testing.T) {
	wantE := fmt.Errorf("invalidQuery is not a valid filter option. Examples of options: synced=no | lastmodified=2024-01-20 | \"lastmodified<2024-01-15\" | \"lastmodified>=2023-12-22\"")
	got, gotE := GetFilterStrategy("invalidQuery")
	if got != nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"got (%v, %v) , want (%v, %v)",
			got, gotE.Error(),
			nil, wantE.Error(),
		)
	}
}

var syncedEQFilterStrategyTests = []struct {
	query string
	want  []Repo
}{
	{
		"synced=yes",
		getFilteredOutputByQuery("synced=yes"),
	},
	{
		"synced=y",
		getFilteredOutputByQuery("synced=yes"),
	},
	{
		"synced=no",
		getFilteredOutputByQuery("synced=no"),
	},
	{
		"synced=n",
		getFilteredOutputByQuery("synced=no"),
	},
}

func TestSyncedEQFilterStrategy(t *testing.T) {
	for _, test := range syncedEQFilterStrategyTests {
		syncedEQFilter := SyncedEQFilterStrategy{}
		repos := getUnfilteredRepos()
		got, _ := syncedEQFilter.Apply(repos, test.query)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(
				"got (%v)\nwant (%v)",
				got,
				test.want,
			)
		}
	}
}

func TestSyncedEQFilterStrategyError(t *testing.T) {
	wantE := fmt.Errorf("incorrect value for synced, value must be either 'yes', 'y', 'no' or 'n'")
	repos := getUnfilteredRepos()
	got, gotE := SyncedEQFilterStrategy{}.Apply(repos, "synced=invalid")
	if got != nil || gotE.Error() != wantE.Error() {
		t.Errorf(
			"got (%v, %v) , want (%v, %v)",
			got, gotE.Error(),
			nil, wantE.Error(),
		)
	}

}

func TestLastModifiedEQFilterStrategy(t *testing.T) {
	want := getFilteredOutputByQuery("lastmodified=2024-01-03")
	repos := getUnfilteredRepos()
	got, _ := LastModifiedEQFilterStrategy{}.Apply(repos, "lastmodified=2024-01-03")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"got (%v)\nwant (%v)",
			got,
			want,
		)
	}
}

func TestLastModifiedLEQFilterStrategy(t *testing.T) {
	want := getFilteredOutputByQuery("lastmodified<=2024-01-03")
	repos := getUnfilteredRepos()
	got, _ := LastModifiedLEQFilterStrategy{}.Apply(repos, "lastmodified<=2024-01-03")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"got (%v)\nwant (%v)",
			got,
			want,
		)
	}
}

func TestLastModifiedGEQFilterStrategy(t *testing.T) {
	want := getFilteredOutputByQuery("lastmodified>=2024-01-03")
	repos := getUnfilteredRepos()
	got, _ := LastModifiedGEQFilterStrategy{}.Apply(repos, "lastmodified>=2024-01-03")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"got (%v)\nwant (%v)",
			got,
			want,
		)
	}
}

func TestLastModifiedLESFilterStrategy(t *testing.T) {
	want := getFilteredOutputByQuery("lastmodified<2024-01-03")
	repos := getUnfilteredRepos()
	got, _ := LastModifiedLESFilterStrategy{}.Apply(repos, "lastmodified<2024-01-03")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"got (%v)\nwant (%v)",
			got,
			want,
		)
	}
}

func TestLastModifiedGRTFilterStrategy(t *testing.T) {
	want := getFilteredOutputByQuery("lastmodified>2024-01-03")
	repos := getUnfilteredRepos()
	got, _ := LastModifiedGRTFilterStrategy{}.Apply(repos, "lastmodified>2024-01-03")
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"got (%v)\nwant (%v)",
			got,
			want,
		)
	}
}

var lastModifiedFilterStrategyErrorTests = []struct {
	query    string
	strategy FilterStrategy
	wantE    error
}{
	{
		"lastmodified=invalid",
		LastModifiedEQFilterStrategy{},
		fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query lastmodified=invalid with date invalid"),
	},
	{
		"lastmodified<=invalid",
		LastModifiedLEQFilterStrategy{},
		fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query lastmodified<=invalid with date invalid"),
	},
	{
		"lastmodified>=invalid",
		LastModifiedGEQFilterStrategy{},
		fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query lastmodified>=invalid with date invalid"),
	},
	{
		"lastmodified<invalid",
		LastModifiedLESFilterStrategy{},
		fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query lastmodified<invalid with date invalid"),
	},
	{
		"lastmodified>invalid",
		LastModifiedGRTFilterStrategy{},
		fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query lastmodified>invalid with date invalid"),
	},
}

func TestLastModifiedFilterStrategyError(t *testing.T) {
	for _, test := range lastModifiedFilterStrategyErrorTests {
		lastModifiedFilterStrategy := test.strategy
		gotE := lastModifiedFilterStrategy.ValidateQuery(test.query)
		if gotE == nil || gotE.Error() != test.wantE.Error() {
			t.Errorf(
				"got (%v)\nwant (%v)",
				gotE,
				test.wantE,
			)
		}
	}
}

func getUnfilteredRepos() []Repo {
	// helper to provide fake input to test filter strategy
	return []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3a,
			true,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			true,
		},
	}
}

func getFilteredOutputByQuery(query string) []Repo {
	// helper to provide expected outputs for each filter strategy apply
	filteredBySyncYes := []Repo{
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3a,
			true,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			true,
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
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			true,
		},
	}
	filteredByLastModifiedEQjan3 := []Repo{
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3a,
			true,
			"",
			true,
		},
	}
	filteredByLastModifiedLEQjan3 := []Repo{
		{
			"c",
			"repos/c",
			"/home/user/repos/z/c",
			jan1,
			false,
			"",
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3a,
			true,
			"",
			true,
		},
	}
	filteredByLastModifiedGEQjan3 := []Repo{
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3,
			false,
			"",
			true,
		},
		{
			"a",
			"repos/a",
			"/home/user/repos/x/a",
			jan3a,
			true,
			"",
			true,
		},
		{
			"b",
			"repos/b",
			"/home/user/repos/y/b",
			jan4,
			true,
			"",
			true,
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
			true,
		},
		{
			"d",
			"repos/d",
			"/home/user/repos/w/d",
			jan2,
			true,
			"",
			true,
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
			true,
		},
	}
	outputOptions := map[string][]Repo{
		"synced=yes":               filteredBySyncYes,
		"synced=no":                filteredBySyncNo,
		"lastmodified=2024-01-03":  filteredByLastModifiedEQjan3,
		"lastmodified<=2024-01-03": filteredByLastModifiedLEQjan3,
		"lastmodified>=2024-01-03": filteredByLastModifiedGEQjan3,
		"lastmodified<2024-01-03":  filteredByLastModifiedLESjan3,
		"lastmodified>2024-01-03":  filteredByLastModifiedGRTjan3,
	}
	return outputOptions[query]
}
