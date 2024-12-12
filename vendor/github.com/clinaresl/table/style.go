// -*- coding: utf-8 -*-
// style.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"errors"
	"regexp"
	"strconv"
)

// ----------------------------------------------------------------------------
// Style
// ----------------------------------------------------------------------------

// Functions
// ----------------------------------------------------------------------------

// newStyle creates a new style with information on the column alignment and,
// optionally, an argument. It is initialized automatically from a string
// extracted from the column specification
func newStyle(spec string) (*style, error) {

	// first things first, verify that the given specification is correct
	re := regexp.MustCompile(columnSpecRegex)
	if !re.MatchString(spec) {
		return &style{}, errors.New("invalid style specification")
	}

	// now, check for the special case of the qualifiers 'p/C/L/R' which accept
	// a numerical argument
	re = regexp.MustCompile(pRegex)
	pmatch := re.FindStringIndex(spec)
	if pmatch == nil {

		// if nono of them is not present, then return a style with the
		// character used instead and no numerical argument
		return &style{alignment: spec[0]}, nil
	}

	// if 'p/C/L/R' has been given, then process separately the numerical
	// argument. Note that the bounds given by FindStringIndex are proerly
	// shifted to get only the digits
	arg, err := strconv.Atoi(spec[2+pmatch[0] : pmatch[1]-1])
	if err != nil || arg <= 0 {
		return &style{}, errors.New("invalid numerical argument in a 'p' column")
	}
	return &style{alignment: spec[0],
		arg: arg}, nil
}
