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
// validateQueries and applyQueries functions, making the access of the query's
// value field value simpler
type query interface {
	value() string
}

type queries struct {
	// place sort at the end as there will be less elements to sort in the
	// slice after filtering
	LastModified lastModifiedFilter
	Synced       syncedFilter
	Author       authorFilter
	Sort         sorter
}

type sortFunc func([]Repo)

// sorts in alphabetical order ascending
func sortByName(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.Name, b.Name)
	})
}

// sorts in alphabetical order ascending
func sortByPath(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(a.AbsPath, b.AbsPath)
	})
}

// sorts in order of last modified datetime ascending
func sortByLastModified(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return a.LastModified.Compare(b.LastModified)
	})
}

// sorts false values first and then true values as it is likely user
// will want to see repos that are not synced first
func sortBySyncStatus(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		if a.SyncedWithRemote && !b.SyncedWithRemote {
			return 1
		} else if !a.SyncedWithRemote && b.SyncedWithRemote {
			return -1
		} else {
			return 0
		}
	})
}

// sorts in alphabetical order ascending
func sortByAuthor(repos []Repo) {
	slices.SortStableFunc(repos, func(a, b Repo) int {
		return strings.Compare(
			// compare lower case to make sort case insensitive
			// because author names will be a mix of capitalized
			// and non-capitalized
			strings.ToLower(a.Author), strings.ToLower(b.Author))
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
	_, ok := s.validOptions[strings.ToLower(s.Value)]
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
	// select the appropriate sort function based on flag value
	sort := s.validOptions[strings.ToLower(s.Value)]
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
	value := strings.ToLower(s.Value)
	if value != "yes" &&
		value != "y" &&
		value != "no" &&
		value != "n" {
		return fmt.Errorf("incorrect value for synced, value must be either 'yes', 'y', 'no' or 'n'")
	}
	return nil
}

func (s syncedFilter) apply(repos *[]Repo) error {
	value := strings.ToLower(s.Value)
	var queryBool bool
	if value == "yes" || value == "y" {
		queryBool = true
	} else {
		// no need to check for no explicitly as it is already covered
		// by validate() method
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
	// check for each prefix before trimming because some prefixes include
	// other prefixes. ex: <= includes <
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

type authorFilter struct {
	Value string
}

func (s authorFilter) value() string {
	return s.Value
}

func (s authorFilter) validate() error {
	// all querues must be validated through ValidateQueries() before applying
	// hence return nil to indicate any value for author is valid
	return nil
}

func (s authorFilter) apply(repos *[]Repo) error {
	var filteredRepos []Repo
	for _, repo := range *repos {
		// case insensitive check
		if strings.Contains(strings.ToLower(repo.Author), strings.ToLower(s.value())) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	*repos = filteredRepos
	return nil
}

// returns a pointer to queries struct and is used to set the filters and sort
// for repos
func NewQueries() *queries {
	return &queries{Sort: sorter{validOptions: map[string]sortFunc{
		// add possible sort flag values and their corresponding sort functions here
		"name":         sortByName,
		"path":         sortByPath,
		"lastmodified": sortByLastModified,
		"synced":       sortBySyncStatus,
		"author":       sortByAuthor,
	}}}
}

func ValidateQueries(q *queries) error {
	var err error
	v := reflect.ValueOf(*q)
	// loop through the fields in Queries and execute the validate method
	// for each query that has been set, meaning the flag was used by the
	// user
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

func ApplyQueries(q *queries, repos *[]Repo) error {
	var err error
	v := reflect.ValueOf(*q)
	// loop through the fields in Queries and execute the apply method
	// for each query that has been set, meaning the flag was used by the
	// user
	for i := 0; i < v.NumField(); i++ {
		query := v.Field(i).Interface().(query)
		// ignore queries where the value has not been set
		// indicating that the flag for the query was not used
		if query.value() == "" {
			continue
		}
		executor := v.Field(i).Interface().(executor)
		err = executor.apply(repos)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReverseSort(repos *[]Repo) error {
	// handle reverse sort separately from queries since the reverse flag is a
	// boolean and does not implement logic for flag values and validate method
	for i, j := 0, len(*repos)-1; i < j; i, j = i+1, j-1 {
		(*repos)[i], (*repos)[j] = (*repos)[j], (*repos)[i]
	}
	return nil
}
