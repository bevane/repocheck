// End to end tests for repocheck
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// temp location to store the fake repos
const root = "/tmp/repochecktest"

func TestMain(m *testing.M) {
	var err error
	var cmd *exec.Cmd

	// setup
	err = setup(root)
	if err != nil {
		// call teardown directly because defer will not work for log.Fatal
		teardown(root)
		log.Fatal(err)
	}

	cmd = exec.Command("go", "build")
	_, err = cmd.CombinedOutput()
	if err != nil {
		// call teardown directly because defer will not work for log.Fatal
		teardown(root)
		log.Fatal(err)
	}

	// teardown
	defer teardown(root)

	m.Run()

}

func TestRepoCheckNoFlags(t *testing.T) {
	// root contains the fake repos that have been set up in TestMain
	// ./repocheck is also built in setup
	cmd := exec.Command("./repocheck", root)
	out, _ := cmd.CombinedOutput()
	got := string(out)
	want := getCLIOutSnapshot()
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}

}

func setup(root string) error {
	var err error
	err = initFakeRepos(root)
	if err != nil {
		return err
	}
	err = setupRepoA(root)
	if err != nil {
		return err
	}
	err = setupRepoB(root)
	if err != nil {
		return err
	}
	err = setupRepoC(root)
	if err != nil {
		return err
	}
	return nil
}

func teardown(root string) {
	os.RemoveAll(root)
	os.Remove("repocheck")
}

func initFakeRepos(root string) error {
	// creates and initializes local and remote repos at root
	var err error
	var cmd *exec.Cmd
	var out []byte
	repos := []string{"a", "b", "c"}
	for _, repo := range repos {
		remotePath := filepath.Join(root, "remote", repo)
		localPath := filepath.Join(root, "local", repo)
		err = os.MkdirAll(remotePath, 0755)
		if err != nil {
			return err
		}

		err = os.MkdirAll(localPath, 0755)
		if err != nil {
			return err
		}

		cmd = exec.Command("git", "config", "user.name")
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %v", err, string(out))
		}

		if string(out) == "" {
			// set test credentials if no credentials are found in git config
			// to avoid errors due to missing credentials
			cmd = exec.Command("git", "config", "--global", "user.name", "Test")
			cmd.Output()
			cmd = exec.Command("git", "config", "--global", "user.email", "test@github.com")
			cmd.Output()
		}

		cmd = exec.Command("git", "init", "--bare")
		cmd.Dir = remotePath
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %v", err, string(out))
		}

		cmd = exec.Command("git", "clone", remotePath, localPath)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %v", err, string(out))
		}

	}
	return nil
}

func setupRepoA(root string) error {
	// repo A setup details
	// last modified : 2024-01-01
	// fully synced
	var err error
	var cmd *exec.Cmd
	var out []byte
	localPath := filepath.Join(root, "local", "a")
	// change the last modified date of the file and directory
	// so the results of test are deterministic
	combinedCommands := `touch file1 &&
			     touch -t 202401011000 file1 &&
			     touch -t 202401011000 . &&
			     git add . &&
			     git commit -m 'add file' &&
			     git push`

	cmd = exec.Command("sh", "-c", combinedCommands)
	cmd.Dir = localPath
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %v", err, string(out))
	}
	return nil
}

func setupRepoB(root string) error {
	// repo B setup details
	// last modified : 2024-01-02
	// has uncommitted changes
	var err error
	var cmd *exec.Cmd
	var out []byte
	localPath := filepath.Join(root, "local", "b")
	// change the last modified date of the file and directory
	// so the results of test are deterministic
	combinedCommands := `touch file1 &&
			     touch -t 202401011000 file1 &&
			     git add . &&
			     git commit -m 'add file' &&
			     git push
			     touch file2 &&
			     touch -t 202401021000 file2
			     touch -t 202401011000 .`

	cmd = exec.Command("sh", "-c", combinedCommands)
	cmd.Dir = localPath
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %v", err, string(out))
	}
	return nil
}

func setupRepoC(root string) error {
	// repo C setup details
	// last modified : 2024-01-02
	// has branches that are ahead and branches with no remote branch
	var err error
	var cmd *exec.Cmd
	var out []byte
	localPath := filepath.Join(root, "local", "c")
	// change the last modified date of the file and directory
	// so the results of test are deterministic
	combinedCommands := `touch file1 &&
			     touch -t 202401011000 file1 &&
			     git add . &&
			     git commit -m 'add file' &&
			     git push
			     touch file2 &&
			     touch -t 202401031000 file2 &&
			     touch -t 202401031000 . &&
			     git add . &&
			     git commit -m 'add file' &&
			     git switch -c newbranch`

	cmd = exec.Command("sh", "-c", combinedCommands)
	cmd.Dir = localPath
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %v", err, string(out))
	}
	return nil
}

func getCLIOutSnapshot() string {
	return `┍━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━┯━━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┑
│      Repo       │ Path                 │ Last Modified │ Synced? │ Sync Details                │
┝━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━━━┿━━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥
│        a        │ /tmp/repochecktest/l │  2024-01-01   │   yes   │                             │
│                 │ ocal/a               │               │         │                             │
├─────────────────┼──────────────────────┼───────────────┼─────────┼─────────────────────────────┤
│        b        │ /tmp/repochecktest/l │  2024-01-02   │   no    │ - has uncommitted changes   │
│                 │ ocal/b               │               │         │                             │
├─────────────────┼──────────────────────┼───────────────┼─────────┼─────────────────────────────┤
│        c        │ /tmp/repochecktest/l │  2024-01-03   │   no    │ - has branch(es) with no    │
│                 │ ocal/c               │               │         │ remote branch               │
│                 │                      │               │         │ - has branch(es) that are   │
│                 │                      │               │         │ ahead                       │
└─────────────────┴──────────────────────┴───────────────┴─────────┴─────────────────────────────┘
3 repos found in /tmp/repochecktest: 2 repo(s) are not synced
`
}
