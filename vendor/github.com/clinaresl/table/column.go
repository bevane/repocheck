// -*- coding: utf-8 -*-
// column.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"errors"
	"regexp"
)

// ----------------------------------------------------------------------------
// Column
// ----------------------------------------------------------------------------

// Functions
// ----------------------------------------------------------------------------

// newColumn returns an instance of a column which is obtained by processing the
// given (single) column specification. It updates the horizontal formatting
// style and sets the vertical formatting style to "top" by default
func newColumn(spec string) (*column, error) {

	// first things first, verify that the given (single) column specification
	// is syntactically correct
	re := regexp.MustCompile(colSpecRegex)
	if !re.MatchString(spec) {
		return &column{}, errors.New("invalid (single) column specification")
	}

	// now, look for the column format
	re = regexp.MustCompile(columnSpecRegex)
	smatch := re.FindStringIndex(spec)
	cstyle, err := newStyle(spec[smatch[0]:smatch[1]])
	if err != nil {
		return &column{}, err
	}

	// so far, distinguish the separator from the format which is processed
	// separately
	return &column{sep: spec[0:smatch[0]],
		hformat: *cstyle,
		vformat: style{alignment: 't'}}, nil
}
