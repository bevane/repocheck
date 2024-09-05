package app_test

import (
	"github.com/bevane/rpchk/app"
	"slices"
	"testing"
	"testing/fstest"
)

func TestListRepoDirectories(t *testing.T) {
	testFsys := fstest.MapFS{
		"repo1/.git/.keep":         {},
		"repo2/.git/.keep":         {},
		"repo2/test/.keep":         {},
		"norepo1/test/.keep":       {},
		"norepo2/test/.keep":       {},
		"norepo2/repo3/.git/.keep": {},
	}
	want := []string{
		"norepo2/repo3",
		"repo1",
		"repo2",
	}
	got := app.ListRepoDirectories(testFsys)
	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
