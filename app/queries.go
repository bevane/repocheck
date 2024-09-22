package app

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

type validator interface {
	validate() error
}

type executor interface {
	apply(*[]Repo) error
}

// this interface is needed to get the value for each query in the
// validateQueries and applyQueries functions to make accessing the query's
// value field value simpler
type query interface {
	value() string
}

type queries struct {
	Sort         sorter
	LastModified lastModifiedFilter
	Synced       syncedFilter
}

type sortFunc func([]Repo)

func sortByName(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortByPath(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.AbsPath, b.AbsPath)
	})
}

func sortByLastModified(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return a.LastModified.Compare(b.LastModified)
	})
}

func sortBySyncStatus(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		// sort to show false values first as it is more likely the
		// user will want to see which repos are not synced when they
		// sort by sync status
		if a.SyncedWithRemote && !b.SyncedWithRemote {
			return 1
		} else if !a.SyncedWithRemote && b.SyncedWithRemote {
			return -1
		} else {
			return 0
		}
	})
}

type sorter struct {
	Value        string
	validOptions map[string]sortFunc
}

func (s sorter) value() string {
	return s.Value
}

func (s sorter) validate() error {
	_, ok := s.validOptions[s.Value]
	if !ok {
		var validOptions []string
		for key := range s.validOptions {
			validOptions = append(validOptions, key)
		}
		// sort the keys to get a deterministic error message
		slices.SortFunc(validOptions, func(a, b string) int {
			return strings.Compare(a, b)
		})
		return fmt.Errorf("%v is not a valid sort option. Options: %v", s.Value, strings.Join(validOptions, " | "))
	}
	return nil
}

func (s sorter) apply(repos *[]Repo) error {
	err := s.validate()
	if err != nil {
		return err
	}
	// select the appropriate sort function based on flag value
	sort := s.validOptions[s.Value]
	sort(*repos)
	return nil
}

type syncedFilter struct {
	Value string
}

func (s syncedFilter) value() string {
	return s.Value
}

func (s syncedFilter) validate() error {
	if s.Value != "yes" &&
		s.Value != "y" &&
		s.Value != "no" &&
		s.Value != "n" {
		return fmt.Errorf("incorrect value for synced, value must be either 'yes', 'y', 'no' or 'n'")
	}
	return nil
}

func (s syncedFilter) apply(repos *[]Repo) error {
	err := s.validate()
	if err != nil {
		return err
	}
	var queryBool bool
	if s.Value == "yes" || s.Value == "y" {
		queryBool = true
	} else {
		queryBool = false
	}
	var filteredRepos []Repo
	for _, repo := range *repos {
		if repo.SyncedWithRemote == queryBool {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	*repos = filteredRepos
	return nil
}

type lastModifiedFilter struct {
	Value string
}

func (s lastModifiedFilter) value() string {
	return s.Value
}

func (l lastModifiedFilter) validate() error {
	var dateString string
	switch {
	case strings.HasPrefix(l.Value, "<="):
		dateString = strings.TrimPrefix(l.Value, "<=")
	case strings.HasPrefix(l.Value, ">="):
		dateString = strings.TrimPrefix(l.Value, ">=")
	case strings.HasPrefix(l.Value, "<"):
		dateString = strings.TrimPrefix(l.Value, "<")
	case strings.HasPrefix(l.Value, ">"):
		dateString = strings.TrimPrefix(l.Value, ">")
	default:
		dateString = l.Value
	}
	_, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		return fmt.Errorf("unexpected date %v, date must be in the format yyyy-mm-dd and can only be prefixed with '<=', '>=', '<' or '>'", dateString)
	}
	return nil
}

func (l lastModifiedFilter) apply(repos *[]Repo) error {
	err := l.validate()
	if err != nil {
		return err
	}
	var dateString string
	var queryDate time.Time
	var filteredRepos []Repo
	switch {
	case strings.HasPrefix(l.Value, "<="):
		dateString = strings.TrimPrefix(l.Value, "<=")
		queryDate, _ = time.Parse(time.DateOnly, dateString)
		for _, repo := range *repos {
			// compare string representations of date to exclude time in comparison
			if repo.LastModified.Format(time.DateOnly) <= queryDate.Format(time.DateOnly) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	case strings.HasPrefix(l.Value, ">="):
		dateString = strings.TrimPrefix(l.Value, ">=")
		queryDate, _ = time.Parse(time.DateOnly, dateString)
		for _, repo := range *repos {
			// compare string representations of date to exclude time in comparison
			if repo.LastModified.Format(time.DateOnly) >= queryDate.Format(time.DateOnly) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	case strings.HasPrefix(l.Value, "<"):
		dateString = strings.TrimPrefix(l.Value, "<")
		queryDate, _ = time.Parse(time.DateOnly, dateString)
		for _, repo := range *repos {
			// compare string representations of date to exclude time in comparison
			if repo.LastModified.Format(time.DateOnly) < queryDate.Format(time.DateOnly) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	case strings.HasPrefix(l.Value, ">"):
		dateString = strings.TrimPrefix(l.Value, ">")
		queryDate, _ = time.Parse(time.DateOnly, dateString)
		for _, repo := range *repos {
			// compare string representations of date to exclude time in comparison
			if repo.LastModified.Format(time.DateOnly) > queryDate.Format(time.DateOnly) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
	default:
		dateString = l.Value
		queryDate, _ = time.Parse(time.DateOnly, dateString)
		for _, repo := range *repos {
			// compare string representations of date to exclude time in comparison
			if repo.LastModified.Format(time.DateOnly) == queryDate.Format(time.DateOnly) {
				filteredRepos = append(filteredRepos, repo)
			}
		}

	}
	*repos = filteredRepos
	return nil
}

func NewQueries() queries {
	return queries{Sort: sorter{validOptions: map[string]sortFunc{
		"name":         sortByName,
		"path":         sortByPath,
		"lastmodified": sortByLastModified,
		"synced":       sortBySyncStatus,
	}}}
}

func ValidateQueries(q queries) error {
	var err error
	v := reflect.ValueOf(q)
	for i := 0; i < v.NumField(); i++ {
		query := v.Field(i).Interface().(query)
		// ignore queries where the value has not been set
		// indicating that the flag for the query was not used
		if query.value() == "" {
			continue
		}
		validator := v.Field(i).Interface().(validator)
		err = validator.validate()
		if err != nil {
			return err
		}
	}
	return nil
}
