package cmd

import (
	"bufio"
	"fmt"
	"github.com/bevane/repocheck/app"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"time"
)

var opt = app.NewQueries()
var tsvOutput bool
var noFetch bool
var LogWriter *bufio.Writer

var rootCmd = &cobra.Command{
	Use:   "repocheck",
	Short: "repocheck is a cli tool to show repos in a directory and info about them",
	Long:  "repocheck is a cli tool to show repos in a directory and info about them - see info for each repo such as absolute path of repo, last modified date, whether repo is synced with remote (whether it has uncommited changes or branches that are ahead etc.)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  repocheckCmd,
}

func init() {
	LogWriter = bufio.NewWriter(os.Stderr)
	log.SetOutput(LogWriter)
	rootCmd.Flags().StringVarP(&opt.Sort.Value, "sort", "s", "lastmodified", "Key to sort the results by. Example: '-s name'. Options: lastmodified | name | path | synced")
	rootCmd.Flags().StringVarP(&opt.Synced.Value, "synced", "S", "", "Filter results by synced status of repo. Example: '-S y' | '-S no'")
	rootCmd.Flags().StringVarP(&opt.LastModified.Value, "lastmodified", "L", "", "Filter results by last modified date of repo. Examples: '-L 2024-01-20' | '--lastmodified \"<2024-01-15\"' | '-L \">=2023-12-22\"'\nNote: surround any filters containing < or > with quotes")
	rootCmd.Flags().BoolVarP(&tsvOutput, "tsv", "t", false, "Output results as tab separated values")
	rootCmd.Flags().BoolVarP(&noFetch, "no-fetch", "", false, "Don't run  git fetch for each repo")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func repocheckCmd(cmd *cobra.Command, args []string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Start()
	var err error
	var root string
	// run validation of flag values in the beginning before proceeding
	// further
	err = app.ValidateQueries(opt)
	if err != nil {
		s.Stop()
		return fmt.Errorf("repocheck: %v", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		s.Stop()
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
	repos, err := app.GetReposWithDetails(root, !noFetch)
	if err != nil {
		s.Stop()
		return fmt.Errorf(
			"repocheck: cannot run check on '%v': %v",
			root,
			err,
		)
	}
	err = app.ApplyQueries(opt, &repos)
	if err != nil {
		s.Stop()
		return fmt.Errorf("repocheck: %v", err)
	}
	var output string
	switch {
	case tsvOutput:
		output = app.ConstructTSVOutput(repos)
	default:
		table, err := app.ConstructTable(repos)
		if err != nil {
			s.Stop()
			return fmt.Errorf(
				"repocheck: error constructing table: %v",
				err,
			)
		}
		summary := app.ConstructSummary(repos, root)
		output = fmt.Sprintf("%v\n%v\n", table, summary)
	}
	s.Stop()
	LogWriter.Flush()
	fmt.Print(output)
	return nil
}
