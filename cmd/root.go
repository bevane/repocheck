package cmd

import (
	"fmt"
	"github.com/bevane/repocheck/app"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var SortKey string

var rootCmd = &cobra.Command{
	Use:   "repocheck",
	Short: "repocheck is a cli tool show repos and info about them in a directory",
	Long:  "repocheck is a cli tool show repos and info about them in a directory - see info for each repo such as absolute path of repo, last modified date, whether repo is synced with remote (whether it has uncommited changtes or branches that are ahead etc.)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  repocheckCmd,
}

func init() {
	rootCmd.Flags().StringVarP(&SortKey, "sort", "s", "lastmodified", "Key to sort the results by. Options: lastmodified | name | path | synced")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func repocheckCmd(cmd *cobra.Command, args []string) error {
	var selectedSortFunc app.SortFunc
	if SortKey != "" {
		var err error
		selectedSortFunc, err = app.GetSortRepoFunc(SortKey)
		if err != nil {
			return fmt.Errorf("repocheck: %v", err)
		}
	}
	var root string
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("repocheck: error getting working dir: %v", err)
	}
	if len(args) == 0 {
		root = wd
	} else {
		pathArg := args[0]
		if filepath.IsAbs(pathArg) {
			root = pathArg
		} else {
			root = filepath.Join(wd, pathArg)
		}
	}
	repos, err := app.GetRepos(root)
	switch SortKey {
	case "name":
	}
	if err != nil {
		return fmt.Errorf(
			"repocheck: cannot run check on '%v': %v",
			root,
			err,
		)
	}
	if selectedSortFunc != nil {
		selectedSortFunc(repos, false)
	}
	table, err := app.ConstructTable(repos)
	if err != nil {
		return fmt.Errorf(
			"repocheck: error constructing table: %v",
			err,
		)
	}
	summary := app.ConstructSummary(repos, root)
	fmt.Printf("%v\n%v\n", table, summary)
	return nil
}
