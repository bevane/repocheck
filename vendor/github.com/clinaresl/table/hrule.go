// -*- coding: utf-8 -*-
// hrule.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Processing a cell means transforming logical rows into physical ones by
// splitting its contents across several (physical) rows, and also adding blank
// lines so that the result satisfies the vertical format of the column where it
// has to be shown, if and only if the height of the corresponding row is larger
// than the number of physical rows necessary to display the contents of the
// cell. To properly process a cell it is necessary to get a pointer to the
// table, and also the integer indices to the row and column of the cell
func (h hrule) Process(t *Table, irow, jcol int) []formatter {

	// Processing a horizontal rule involves computing a single string which
	// contains the separator of the j-th column where vertical separators are
	// just substituted by blanks because these intersections are computed
	// later.
	var splitters string

	// get the vertical separator to process
	sep := t.columns[jcol].sep

	// search for ANSI color escape sequences
	re := regexp.MustCompile(ansiColorRegex)
	colindexes := re.FindAllStringIndex(sep, -1)

	// position at the first color and annotate how many have been found
	colind, nbcolors := 0, len(colindexes)

	// the following value should be equal to -1 if we have not found a vertical
	// separator yet and 0 otherwise. It is used to decide what horizontal rule
	// to use (either the one from the preceding column or the one following
	// immediately after) in case a rune found in the separator has to be
	// substituted
	offset := -1

	// process all runes in the current separator
	for idx, irune := range sep {

		// ANSI color escape sequences have to be directly copied to the
		// splitters
		if colind < nbcolors && idx >= colindexes[colind][0] {

			// if the ANSI color escape sequence starts right here then copy it
			// to the splitter
			if idx == colindexes[colind][0] {
				splitters += sep[colindexes[colind][0]:colindexes[colind][1]]
			}

			// if this position ends the entire ANSI color sequence, then move
			// to the next color
			if idx == colindexes[colind][1]-1 {
				colind++
			}

			// and skip the treatment of this rune (character)
			continue
		}

		// in addition, in case this rune is a vertical separator then make sure
		// the offset is 0. Note that if the string used in the separator
		// contains more than one vertical separator, the next column is
		// considered immediately after the first vertical separator, in spite
		// of the number of them
		if isVerticalSeparator(irune) {
			offset = 0
		}

		// otherwise, take the horizontal rule either before or after a
		// vertical separator, if any has been found. In particular, if we
		// are before the first column and no vertical separator has been
		// found yet, or if we are at the last column after any vertical
		// separator, then just copy this rune
		if (offset == -1 && jcol == 0) ||
			(offset == 0 && jcol == len(t.columns)-1) {
			splitters += string(irune)
		} else {

			// if, on the other hand, we are anywhere between the first
			// column after a vertical separator and the last column before
			// a vertical separator, then take the horizontal rule used in
			// the corresponding cell
			brkrule, _ := utf8.DecodeRuneInString(string(t.cells[irow][jcol+offset].(hrule)))
			splitters += string(brkrule)
		}
	}

	// and return the string computed so far but in the form of a slice
	// containing only one line. Mind the trick: the splitters (which are a
	// standard string) are casted into a hrule to enable the specific format of
	// horizontal rules
	return []formatter{hrule(splitters)}
}

// Cells are also formatted (physical) line by line where each physical line is
// the result of processing cell (irow, jcol) and should be given in the
// receiver of this method. Each invocation returns a string where each
// (physical) line is forrmatted according to the horizontal format
// specification of the j-th column.
func (h hrule) Format(t *Table, irow, jcol int) string {

	// The result of formatting a horizontal rule consists of prefixing the
	// horizontal rule used in the data column with the splitters given in this
	// horizontal rule

	// the only task to do consists of repeating the separator as many times as
	// the width of this column after the splitters
	rule, ok := t.cells[irow][jcol].(hrule)
	if !ok {
		log.Fatalf(" The formatter in location (%v, %v) could not be casted into a rule!", irow, jcol)
	}

	return string(h) + strings.Repeat(string(rule), t.columns[jcol].width)
}
