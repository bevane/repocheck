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
	Long:  "repocheck is a cli tool to show repos in a directory and info about them - see info for each repo such as absolute path of repo, last modified date, whether repo is synced with remote (whether it has uncommitted changes or branches that are ahead etc.)",
	Args:  cobra.MaximumNArgs(1),
	// allow shell to autocomplete file path for the first argument
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) < 1 {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: repocheckCmd,
}

func init() {
	rootCmd.AddCommand(completionCmd)
	LogWriter = bufio.NewWriter(os.Stderr)
	log.SetOutput(LogWriter)
	rootCmd.Flags().StringVarP(&opt.Sort.Value, "sort", "s", "lastmodified", "Sort results\noptions: lastmodified | name | path | synced")
	rootCmd.Flags().StringVarP(&opt.Synced.Value, "synced", "S", "", "Filter by synced status of repo\noptions: y | n")
	rootCmd.Flags().StringVarP(&opt.LastModified.Value, "lastmodified", "L", "", "Filter by last modified date of repo\noptions: yyyy-mm-dd | \">yyyy-mm-dd\" | \">=yyyy-mm-dd\"\nnote: surround any filters containing < or > with quotes")
	rootCmd.Flags().BoolVarP(&tsvOutput, "tsv", "t", false, "Output as tab separated values")
	rootCmd.Flags().BoolVarP(&noFetch, "no-fetch", "", false, "Run without doing a git fetch for each repo")
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
