package app

import (
	"reflect"
	"testing"
	"testing/fstest"
	"time"
)

func TestListRepoDirectories(t *testing.T) {
	modTime, _ := time.Parse(time.RFC3339, "2024-01-01T15:00:00Z")
	testFsys := fstest.MapFS{
		"repo1/.git/.keep":         {ModTime: modTime},
		"repo1/.keep":              {ModTime: modTime},
		"repo2/.git/.keep":         {ModTime: modTime},
		"repo2/test/.keep":         {ModTime: modTime},
		"norepo1/test/.keep":       {ModTime: modTime},
		"norepo2/test/.keep":       {ModTime: modTime},
		"norepo2/repo3/.git/.keep": {ModTime: modTime},
		"norepo2/repo3/.keep":      {ModTime: modTime},
	}
	want := []Repo{
		{
			Name:         "repo3",
			Path:         "norepo2/repo3",
			LastModified: modTime,
		},
		{
			Name:         "repo1",
			Path:         "repo1",
			LastModified: modTime,
		},
		{
			Name:         "repo2",
			Path:         "repo2",
			LastModified: modTime,
		},
	}
	got := ListRepoDirectories(testFsys)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetContentLastModifiedDate(t *testing.T) {
	tOld, _ := time.Parse(time.RFC3339, "2024-01-01T15:00:00Z")
	tNew, _ := time.Parse(time.RFC3339, "2024-01-02T13:00:00Z")
	testFsys := fstest.MapFS{
		"test/test1/testfile.test": {ModTime: tOld},
		"test/test2/testfile.test": {ModTime: tNew},
		"test/testfile.test":       {ModTime: tOld},
	}
	want := tNew
	got := getContentLastModifiedTime(testFsys)
	if !got.Equal(want) {
		t.Errorf("got %v want %v", got, want)
	}
}
