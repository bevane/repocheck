package app

import (
	"fmt"
	"strings"
	"time"
)

type FilterStrategy interface {
	ValidateQuery(string) error
	Apply([]Repo, string) ([]Repo, error)
}

type SyncedEQFilterStrategy struct{}

func (s SyncedEQFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "synced=")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'synced' but got query: %v", query)
	}
	if queryValue != "yes" &&
		queryValue != "y" &&
		queryValue != "no" &&
		queryValue != "n" {
		return fmt.Errorf("incorrect value for synced, value must be either 'yes', 'y', 'no' or 'n'")
	}
	return nil
}

func (s SyncedEQFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "synced=")
	var queryBool bool
	if queryValue == "yes" || queryValue == "y" {
		queryBool = true
	} else {
		queryBool = false
	}
	var filteredRepos []Repo
	for _, repo := range repos {
		if repo.SyncedWithRemote == queryBool {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

type LastModifiedEQFilterStrategy struct{}

func (s LastModifiedEQFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "lastmodified=")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'lastmodified' but got query: %v", query)
	}
	_, err := time.Parse(time.DateOnly, queryValue)
	if err != nil {
		return fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query %v with date %v", query, queryValue)
	}
	return nil
}

func (s LastModifiedEQFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "lastmodified=")
	queryDate, _ := time.Parse(time.DateOnly, queryValue)
	var filteredRepos []Repo
	for _, repo := range repos {
		// compare string representations of date to exclude time in comparison
		if repo.LastModified.Format(time.DateOnly) == queryDate.Format(time.DateOnly) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

type LastModifiedLEQFilterStrategy struct{}

func (s LastModifiedLEQFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "lastmodified<=")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'lastmodified' but got query: %v", query)
	}
	_, err := time.Parse(time.DateOnly, queryValue)
	if err != nil {
		return fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query %v with date %v", query, queryValue)
	}
	return nil
}

func (s LastModifiedLEQFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "lastmodified<=")
	queryDate, _ := time.Parse(time.DateOnly, queryValue)
	var filteredRepos []Repo
	for _, repo := range repos {
		// compare string representations of date to exclude time in comparison
		if repo.LastModified.Format(time.DateOnly) <= queryDate.Format(time.DateOnly) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

type LastModifiedGEQFilterStrategy struct{}

func (s LastModifiedGEQFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "lastmodified>=")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'lastmodified' but got query: %v", query)
	}
	_, err := time.Parse(time.DateOnly, queryValue)
	if err != nil {
		return fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query %v with date %v", query, queryValue)
	}
	return nil
}

func (s LastModifiedGEQFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "lastmodified>=")
	queryDate, _ := time.Parse(time.DateOnly, queryValue)
	var filteredRepos []Repo
	for _, repo := range repos {
		// compare string representations of date to exclude time in comparison
		if repo.LastModified.Format(time.DateOnly) >= queryDate.Format(time.DateOnly) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

type LastModifiedLESFilterStrategy struct{}

func (s LastModifiedLESFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "lastmodified<")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'lastmodified' but got query: %v", query)
	}
	_, err := time.Parse(time.DateOnly, queryValue)
	if err != nil {
		return fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query %v with date %v", query, queryValue)
	}
	return nil
}

func (s LastModifiedLESFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "lastmodified<")
	queryDate, _ := time.Parse(time.DateOnly, queryValue)
	var filteredRepos []Repo
	for _, repo := range repos {
		// compare string representations of date to exclude time in comparison
		if repo.LastModified.Format(time.DateOnly) < queryDate.Format(time.DateOnly) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}

type LastModifiedGRTFilterStrategy struct{}

func (s LastModifiedGRTFilterStrategy) ValidateQuery(query string) error {
	queryValue, ok := strings.CutPrefix(query, "lastmodified>")
	if !ok {
		return fmt.Errorf("query mismatch, expected query to start with 'lastmodified' but got query: %v", query)
	}
	_, err := time.Parse(time.DateOnly, queryValue)
	if err != nil {
		return fmt.Errorf("unexpected date, date must be in the format yyyy-mm-dd but got query %v with date %v", query, queryValue)
	}
	return nil
}

func (s LastModifiedGRTFilterStrategy) Apply(repos []Repo, query string) ([]Repo, error) {
	err := s.ValidateQuery(query)
	if err != nil {
		return nil, err
	}
	queryValue, _ := strings.CutPrefix(query, "lastmodified>")
	queryDate, _ := time.Parse(time.DateOnly, queryValue)
	var filteredRepos []Repo
	for _, repo := range repos {
		// compare string representations of date to exclude time in comparison
		if repo.LastModified.Format(time.DateOnly) > queryDate.Format(time.DateOnly) {
			filteredRepos = append(filteredRepos, repo)
		}
	}
	return filteredRepos, nil
}
func GetFilterStrategy(query string) (FilterStrategy, error) {
	switch {
	case strings.HasPrefix(query, "synced="):
		return SyncedEQFilterStrategy{}, nil
	case strings.HasPrefix(query, "lastmodified="):
		return LastModifiedEQFilterStrategy{}, nil
	case strings.HasPrefix(query, "lastmodified<="):
		return LastModifiedLEQFilterStrategy{}, nil
	case strings.HasPrefix(query, "lastmodified>="):
		return LastModifiedGEQFilterStrategy{}, nil
	case strings.HasPrefix(query, "lastmodified<"):
		return LastModifiedLESFilterStrategy{}, nil
	case strings.HasPrefix(query, "lastmodified>"):
		return LastModifiedGRTFilterStrategy{}, nil
	default:
		return nil, fmt.Errorf("%v is not a valid filter option. Examples of options: synced=no | lastmodified=2024-01-20 | \"lastmodified<2024-01-15\" | \"lastmodified>=2023-12-22\"", query)
	}
}
