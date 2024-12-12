// -*- coding: utf-8 -*-
// content.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"fmt"
	"regexp"
	"unicode"
)

// Processing a cell means transforming logical rows into physical ones by
// splitting its contents across several (physical) rows, and also adding blank
// lines so that the result satisfies the vertical format of the column where it
// has to be shown, if and only if the height of the corresponding row is larger
// than the number of physical rows necessary to display the contents of the
// cell. To properly process a cell it is necessary to get a pointer to the
// table, and also the integer indices to the (logical) row and column of the
// cell
func (c content) Process(t *Table, irow, jcol int) []formatter {

	// Processing a content involves splitting it across as many physical rows
	// as needed (e.g., if a "paragraph" alignment is given for this column). In
	// case the number of physical lines is strictly less than the number of
	// physical rows of this logical row, then compute the vertical alignment.
	// The result is given as a slice of contents so that they can be each
	// properly formatted.
	var result []content

	// If the rightmost column is required to be processed and it contains no
	// data, then just simply return the separator
	if jcol == len(t.columns)-1 && t.columns[jcol].hformat.alignment == 0 {

		// if the height of this row is known, then return as many copies of the
		// separator as required. Otherwise, return the separator only once
		for iline := 0; iline < max[int](1, t.rows[irow].height); iline++ {
			result = append(result, content(horizontal_empty))
		}
	} else {

		// aliasing
		col := t.columns[jcol]

		// get the number of physical rows of this logical row taking into account
		// that this logical row might not have been added to the table yet
		var nbrows int
		if irow < len(t.rows) {
			nbrows = t.rows[irow].height
		} else {
			nbrows = 0
		}

		// if a paragraph alignment (p, C, L, R) modifier is used for this specific
		// column, then split the content
		if col.hformat.alignment == 'p' ||
			col.hformat.alignment == 'C' ||
			col.hformat.alignment == 'L' ||
			col.hformat.alignment == 'R' {
			result = strToContent(splitParagraph(string(c), col.hformat.arg))
		} else {

			// if, on the other hand, a newline character has been provided, split the
			// content as well according to the newline characters
			re := regexp.MustCompile(newlineRegex)
			result = strToContent(re.Split(string(c), -1))
		}

		// if the number of physical rows of this logical row is strictly larger
		// than the number of lines necessary to display this content, then apply
		// the vertical format
		if nbrows > len(result) {

			var prefix, suffix int

			// computation of the prefix: number of blank lines before the
			// content
			if unicode.ToLower(rune(col.vformat.alignment)) == 'c' {
				prefix = (nbrows - len(result)) / 2
			}
			if unicode.ToLower(rune(col.vformat.alignment)) == 'b' {
				prefix = nbrows - len(result)
			}

			if unicode.ToLower(rune(col.vformat.alignment)) == 't' {
				suffix = nbrows - len(result)
			}
			if unicode.ToLower(rune(col.vformat.alignment)) == 'c' {
				suffix = (nbrows - len(result)) / 2
				suffix += (nbrows - len(result)) % 2
			}

			// computation of the suffix: number of blank lines after the
			// content
			if unicode.ToLower(rune(col.vformat.alignment)) == 'c' {
				suffix = (nbrows - len(result)) / 2
				suffix += (nbrows - len(result)) % 2
			}

			// and now add the corresponding number of blank lines as required
			for iline := 0; iline < prefix; iline++ {
				result = prepend(horizontal_empty, result)
			}
			for iline := 0; iline < suffix; iline++ {
				result = append(result, horizontal_empty)
			}
		}
	}

	// syntactical issue: explicitly transform the slice of contents into a
	// slice of formatters. Note this transformation guarantees casting back the
	// formatters into contents so that they can be properly formatted each
	var output []formatter
	for _, val := range result {
		output = append(output, formatter(val))
	}

	// and return them
	return output
}

// Cells are also formatted (physical) line by line where each physical line is
// the result of processing cell (irow, jcol) and should be given in the
// receiver of this method. Each invocation returns a string where each
// (physical) line is forrmatted according to the horizontal format
// specification of the j-th column.
func (c content) Format(t *Table, irow, jcol int) string {

	// Specifically speaking, formatting a content means computing the
	// horizontal format and to prefix the result with the separator of the
	// jcol-th column

	// aliasing
	col := t.columns[jcol]

	// in case it is necessary, the prefix and suffix contain a string of blank
	// characters to insert properly so that the contents satisfy the format of
	// this column
	prefix, suffix := justifyLine(string(c), rune(col.hformat.alignment), col.width)

	// get the separator to use
	sep := t.columns[jcol].sep

	// and return the concatenation of the prefix, the content and the suffix,
	// all prefixed with the horizontal separator of the jcol-th column
	return fmt.Sprintf("%v%v", sep, prefix+string(c)+suffix)
}
