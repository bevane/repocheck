package cmd

import (
	"fmt"
	"github.com/bevane/repocheck/app"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var opt = app.NewQueries()

var rootCmd = &cobra.Command{
	Use:   "repocheck",
	Short: "repocheck is a cli tool to show repos in a directory and info about them",
	Long:  "repocheck is a cli tool to show repos in a directory and info about them - see info for each repo such as absolute path of repo, last modified date, whether repo is synced with remote (whether it has uncommited changes or branches that are ahead etc.)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  repocheckCmd,
}

func init() {
	rootCmd.Flags().StringVarP(&opt.Sort.Value, "sort", "s", "lastmodified", "Key to sort the results by. Example: '-s name'. Options: lastmodified | name | path | synced")
	rootCmd.Flags().StringVarP(&opt.Synced.Value, "synced", "S", "", "Filter results by synced status of repo. Example: '-S y' | '-S no'")
	rootCmd.Flags().StringVarP(&opt.LastModified.Value, "lastmodified", "L", "", "Filter results by last modified date of repo. Examples: '-L 2024-01-20' | '--lastmodified \"<2024-01-15\"' | '-L \">=2023-12-22\"'\nNote: surround any filters containing < or > with quotes")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func repocheckCmd(cmd *cobra.Command, args []string) error {
	var err error
	var root string
	// run validation of flag values in the beginning before proceeding
	// further
	err = app.ValidateQueries(opt)
	if err != nil {
		return fmt.Errorf("repocheck: %v", err)
	}
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
	if err != nil {
		return fmt.Errorf(
			"repocheck: cannot run check on '%v': %v",
			root,
			err,
		)
	}
	err = app.ApplyQueries(opt, &repos)
	if err != nil {
		return fmt.Errorf("repocheck: %v", err)
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
