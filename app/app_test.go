package app

import (
	"fmt"
	"reflect"
	"testing"
	"testing/fstest"
	"time"
)

// Functions that do file io and require actual directories to exist
// are not covered in these unit tests. Those functions are covered in the main
// e2e test

func TestListRepoPaths(t *testing.T) {
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
	want := []string{
		"norepo2/repo3",
		"repo1",
		"repo2",
	}
	got, _ := listRepoPaths(testFsys)
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
	got, _ := getContentLastModifiedTime(testFsys)
	if !got.Equal(want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestEvaluateCommitSyncStatus(t *testing.T) {
	var tests = []struct {
		gitOut     string
		wantBool   bool
		wantString string
	}{
		{
			"",
			true,
			"",
		},
		{
			"M app/app_test.go",
			false,
			"- uncommitted changes\n",
		},
		{
			"M  .jest.config.json\nM  package.json\nA  main.js",
			false,
			"- uncommitted changes\n",
		},
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("%v", tt.gitOut)
		t.Run(testname, func(t *testing.T) {
			gotBool, gotString := evaluateCommitSyncStatus(tt.gitOut)
			if gotBool != tt.wantBool || gotString != tt.wantString {
				t.Errorf(
					"got (%v, %v) , want (%v, %v)",
					gotBool, gotString,
					tt.wantBool, tt.wantString,
				)
			}
		})
	}
}

func TestEvaluateBranchSyncStatus(t *testing.T) {
	var tests = []struct {
		gitOut     string
		wantBool   bool
		wantString string
	}{
		{
			"=",
			true,
			"",
		},
		{
			"",
			false,
			"- untracked branch(es)\n",
		},
		{
			">",
			false,
			"- branch(es) ahead\n",
		},
		{
			"<",
			false,
			"- branch(es) behind\n",
		},
		{
			"\n=",
			false,
			"- untracked branch(es)\n",
		},
		{
			"\n=\n>",
			false,
			"- untracked branch(es)\n- branch(es) ahead\n",
		},
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("%v", tt.gitOut)
		t.Run(testname, func(t *testing.T) {
			gotBool, gotString := evaluateBranchSyncStatus(tt.gitOut)
			if gotBool != tt.wantBool || gotString != tt.wantString {
				t.Errorf(
					"got (%v, %v) , want (%v, %v)",
					gotBool, gotString,
					tt.wantBool, tt.wantString,
				)
			}
		})
	}
}
