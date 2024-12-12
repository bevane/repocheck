// -*- coding: utf-8 -*-
// helpers.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

// Functions
// ----------------------------------------------------------------------------

// Return the minimum of two ordered items
func min[T constraints.Ordered](n, m T) T {
	if n > m {
		return m
	}
	return n
}

// Return the maximum of two ordered items
func max[T constraints.Ordered](n, m int) int {
	if n > m {
		return n
	}
	return m
}

// process the given column specification and return a slice of instances of
// columns properly initialized. In case the parsing was not possible an error
// is returned
func getColumns(colspec string) ([]column, error) {

	// --initialization
	var columns []column

	// the specification is processed with a regular expression which should be
	// used to consume the whole string
	re := regexp.MustCompile(colSpecRegex)
	for {

		// get the next column and, if none is found, then exit
		recol := re.FindStringIndex(colspec)
		if recol == nil {
			break
		}

		// in case creating the new column raises an error then return the
		// provisional columns and the error
		nxtcol, err := newColumn(colspec[recol[0]:recol[1]])
		if err != nil {
			return []column{}, err
		}

		// add the new column to the slice of columns to return
		columns = append(columns, *nxtcol)

		// and now move forward in the column specification string
		colspec = colspec[recol[1]:]
	}

	// maybe the column specification string is not empty here. Any remainings
	// are interpreted as the separator of a last column which contains no text
	// and which has no format
	if colspec != "" {
		columns = append(columns,
			column{sep: colspec,
				hformat: style{},
				vformat: style{}})
	}

	// return the slice of columns along with no error
	return columns, nil
}

// replace the ASCII vertical bars found in the given string by the
// corresponding UTF-8 vertical separators
func separatorToUTF8(input *string) {
	*input = strings.ReplaceAll(*input, "|||", "┃")
	*input = strings.ReplaceAll(*input, "||", "║")
	*input = strings.ReplaceAll(*input, "|", "│")
}

// process the given specification according to the specified regex (which must
// match either a column or row specification) and return: first, a new one
// which has removed the specification of the last column/row if and only if a
// last column/row with no column/row specifier was given; second, the separator
// of the last column/row that was removed in the first place
func stripLastSeparator(colspec string, rexp string) (string, string) {

	// -- initialization
	var output string

	// the specification is processed with a regular expression. Only those
	// parts matching the regular expression are returned so that if a last
	// column with no specifier is given, it is not added to the result
	re := regexp.MustCompile(rexp)
	for {

		// get the next column and, if none is found, then exit
		recol := re.FindStringIndex(colspec)
		if recol == nil {
			break
		}

		// copy this part into the output
		output += colspec[recol[0]:recol[1]]

		// and move forward in the column specification string
		colspec = colspec[recol[1]:]
	}

	// and return the string computed so far substituting the last separator by
	// its corresponding UTF-8 runes
	separatorToUTF8(&colspec)
	return output, colspec
}

// return a pointer to the preceding multicell in the i-th row which reaches the
// j-th column or nil if no multicell is found.
func getPreviousHorizontalMerger(t *Table, i, j int) *multicell {

	// this function just iterates over all cells of the i-th row until the j-th
	// column
	idx := 0
	for idx < j && idx < len(t.cells[i]) {

		// if a multicell is found at this location
		if cell, ok := t.cells[i][idx].(multicell); ok {

			// and it ends precisely at the j-th column
			if cell.getColumnInit()+cell.getNbColumns() == j {
				return &cell
			}

			// otherwise, move forward
			idx += cell.getNbColumns()
		} else {

			// otherwise, just move one cell forward
			idx += 1
		}
	}

	// at this point, no multirow has been found so that just return nil
	return nil
}

// return true if and only if the given rune is recognized as a vertical
// separator as defined in this package and false otherwise
func isVerticalSeparator(r rune) bool {
	return r == '│' || r == '║' || r == '┃'
}

// Just cast a slice of strings into a slice of contents
func strToContent(input []string) (output []content) {

	for _, str := range input {
		output = append(output, content(str))
	}
	return
}

// Return the number of runes in the given string which are both printable and
// graphic. It also skips color ANSI codes
func countPrintableRuneInString(s string) (count int) {

	// -- initialization: idx is used to count physical runes, i.e., the
	// physical location of each rune considering also the ANSI color codes
	idx := 0

	// regular expression used to recognize ANSI color codes
	re := regexp.MustCompile(ansiColorRegex)

	// get the indices to all matches of the regular expression for recognizing
	// ANSI color codes, and go then over all runes
	for colind, colindexes := 0, re.FindAllStringIndex(s, -1); idx < len(s); {

		// verify if a ANSI color code starts right at this position
		if colind < len(colindexes) && idx == colindexes[colind][0] {

			// then jump to the first location after the regular expression, and
			// move to the next match of the ANSI color codes
			idx = colindexes[colind][1]
			colind++
		} else {

			// get the rune at the current position, and count it in case it is
			// both printable and graphic
			r, size := utf8.DecodeRuneInString(s[idx:])
			if unicode.IsGraphic(r) && unicode.IsPrint(r) {
				count++
			}

			// and move forward
			idx += size
		}
	}

	return
}

// the following function returns a slice of strings with the same contents than
// the input string (with some spaces removed) such that the length of each
// string is the larger one less or equal than the given width
func splitParagraph(str string, width int) (result []string) {

	// iterate over all runes of the input string
	for len(str) > 0 {

		// while processing a substring, keep track of the number of runes in it
		// and also the location of the last byte to add to it. In addition, it
		// is required to store the position of the rune to start considering in
		// the next cycle
		var nbrunes, end, nxt int
		for pos, rune := range str {

			// accept this rune
			nbrunes++

			// in case this is a space (including utf-8 spaces) then remember
			// the location of the last position to include in the current
			// substring
			if unicode.IsSpace(rune) {
				end, nxt = pos, utf8.RuneLen(rune)

				// and, in case this is a newline character, then exit
				// immediately from the inner loop
				if rune == '\n' {
					break
				}
			}

			// If the maximum number of runes to add has been reached then break
			// avoiding adding more runes
			if nbrunes >= width {

				// if no breaking point has been found before then add all runes
				// until the current location
				if end == 0 {
					end, nxt = pos+utf8.RuneLen(rune), 0
				}

				// If the character immediately after this one is a space then
				// add all runes until this location also
				nxtrune, _ := utf8.DecodeRuneInString(str[pos+utf8.RuneLen(rune):])
				if unicode.IsSpace(nxtrune) {
					end, nxt = pos+utf8.RuneLen(rune), utf8.RuneLen(rune)
				}

				break
			}

			// Finally, if the whole string has been exhausted, then add it
			// until the end
			if pos+utf8.RuneLen(rune) >= len(str) {
				end, nxt = len(str), 0
			}
		}

		// add the substring from the beginning of the input string until the
		// end
		result = append(result, str[:end])

		// and move forward in the string
		str = str[end+nxt:]
	}

	return
}

// A (physical) line is just a string and they can be justified in various ways
// according to the alignment parameter: 'l', 'c', 'r', ... To get the desired
// effect, the contents of the line have to be preceded and continued by a
// prefix and suffix of white spaces which are returned in the output params
// prefix and suffix respectively
func justifyLine(line string, alignment rune, width int) (prefix, suffix string) {

	// compute the prefix to use for representing this line
	if unicode.ToLower(rune(alignment)) == 'c' {
		prefix = strings.Repeat(string(horizontal_blank), (width-countPrintableRuneInString(line))/2)
	}
	if unicode.ToLower(rune(alignment)) == 'r' {
		prefix = strings.Repeat(string(horizontal_blank), width-countPrintableRuneInString(line))
	}

	// compute the suffix to use for representing the contents of this line
	if unicode.ToLower(rune(alignment)) == 'c' {

		// note that in this case an additional character is added, i.e.,
		// centered strings are ragged left in case the difference is and odd
		// number
		suffix = strings.Repeat(string(horizontal_blank), (width-countPrintableRuneInString(line))/2)
		suffix += strings.Repeat(" ", (width-countPrintableRuneInString(line))%2)
	}
	if unicode.ToLower(rune(alignment)) == 'l' || alignment == 'p' {
		suffix = strings.Repeat(string(horizontal_blank), width-countPrintableRuneInString(line))
	}

	// and return the prefix and suffix computed so far
	return
}

// return the rune that splits the four regions north-west, north-east,
// south-west and south-east as stored in the map of splitters with no error. In
// case that any of the runes given to the west, east, north and south is not
// defined in the map of runes, then it is automatically substituted by none
func getSingleSplitter(west, east, north, south rune) rune {

	// check for the existence of the west rune. In case it does not exist,
	// take none
	if _, ok := splitterUTF8[west]; !ok {
		west = none
	}

	// east
	if _, ok := splitterUTF8[west][east]; !ok {
		east = none
	}

	// north
	if _, ok := splitterUTF8[west][east][north]; !ok {
		north = none
	}

	// south
	if _, ok := splitterUTF8[west][east][north][south]; !ok {
		south = none
	}

	// and return the corresponding splitter which, at this point, is guaranteed
	// to exist
	return splitterUTF8[west][east][north][south]
}

// return a slice of vertical specifications as a slice of styles. In case the
// row specification is incorrect, an error is returned and the contents of the
// result are undetermined
func getVerticalStyles(rowspec string) ([]style, error) {

	var result []style

	// while the row specification is not empty. Yeah, the row specification
	// should not consist of runes but just simple ascii characters. Still, we
	// traverse the string as runes
	for _, rune := range rowspec {
		switch rune {
		case 't', 'b', 'c':
			result = append(result, style{alignment: byte(rune)})
		default:
			return result, fmt.Errorf("'%v' is an incorrect vertical format", string(rune))
		}
	}

	return result, nil
}

// The following function prepends the given argument to the slice of contents
// given second
func prepend(item content, data []content) []content {

	// just add an item to the slice, copy all items shifting them all by one
	// position to the right and overwrite the first item
	data = append(data, item)
	copy(data[1:], data)
	data[0] = item

	return data
}

// Evenly increment the width of all columns given in the slice of columns so
// that their accumulated sum is incremented by n
func distributeColumns(n int, columns []column) {

	// compute first the quotient (the amount of space to add to all columns)
	// and the remainder (the additional space to add to a subset of the
	// columns)
	quotient, remainder := n/len(columns), n%len(columns)

	// if and only if the space left to distribute is strictly larger or equal
	// than the number of columns
	if n >= len(columns) {

		// distribute the quotient among all columns
		for idx, _ := range columns {
			columns[idx].width += quotient
		}
	}

	// and now distribute the remainder among the first columns
	for idx := 0; idx < remainder; idx++ {
		columns[idx].width++
	}
}

// Evenly increment the height of all rows given in the slice of rows so that
// their accumulated sum is incremented by n
func distributeRows(n int, rows []row) {

	// compute first the quotient (the amount of space to add to all rows) and
	// the remainder (the additional space to add to a subset of the rows)
	quotient, remainder := n/len(rows), n%len(rows)

	// if and only if the space left to distribute is strictly larger or equal
	// than the number of rows
	if n >= len(rows) {

		// distribute the quotient among all rows
		for idx, _ := range rows {
			rows[idx].height += quotient
		}
	}

	// and now distribute the remainder among the first columns
	for idx := 0; idx < remainder; idx++ {
		rows[idx].height++
	}
}

// return the pi-th physical rune which is known to take the li-th logical
// position. A position is said to be physical if and only if it also takes into
// account control codes such as ANSI color codes; it is logical otherwise.
//
// If such position does not exist it returns -1 unless force is True in which
// case the string is extended to have li logical positions and its physical
// position is then returned.
//
// Because the input string might have been modified or not, it returns the
// resulting string after seeking the physical location of the li-th logical
// position
func logicalToPhysical(s string, li int, force bool) (pi int, sout string) {

	// -- initialization: idx is used to count logical runes---i.e., without
	// considering ANSI color codes
	idx := 0

	// regular expression used to recognize ANSI color codes
	re := regexp.MustCompile(ansiColorRegex)

	// get the indices to all matches of the regular expression for recognizing
	// ANSI color codes, and go then over all runes in the given string until
	// the current logical location goes beyond the logical location requested
	for colind, colindexes := 0, re.FindAllStringIndex(s, -1); pi < len(s) && idx <= li; {

		// verify if a ANSI color code starts right at this position
		if colind < len(colindexes) && pi == colindexes[colind][0] {

			// then jump to the first physical location after the regular
			// expression, and move to the next match of the ANSI color codes
			pi = colindexes[colind][1]
			colind++
		} else {

			// if this is the rune taking the li-th logical position then return
			// its physical location without modifying the input string
			if idx == li {
				return pi, s
			}

			// get the rune at the current position
			_, size := utf8.DecodeRuneInString(s[pi:])

			// and move forward
			pi += size
			idx++
		}
	}

	// if we get here is because the given logical location has not been found.
	// If force has been given, then extend the string so that it contains
	// exactly li logical positions
	if force && li >= 0 {

		// compute the number of extra spaces that have to be added to the
		// string
		diff := 1 + li - countPrintableRuneInString(s)

		// and return the physical location of the newly *created* logical
		// location li along with the new string
		return pi + diff - 1, s + strings.Repeat(string(horizontal_blank), diff)
	}

	// If force is false, an impossible value is returned as a token to signal
	// this case withouth modifying the input string
	return -1, s
}

// return the i-th printable and graphic rune in the given string, if it exists.
// Otherwise, return an emtpy rune along with an error. It skips color ANSI
// codes
func getRune(s string, i int) (rune, error) {

	// -- initialization: idx is used to count physical runes, i.e., the
	// physical location of each rune considering also the ANSI color codes,
	// whereas li is used to count logical runes, i.e., those after disregarding
	// the color ANSI codes
	idx, li := 0, 0

	// regular expression used to recognize ANSI color codes
	re := regexp.MustCompile(ansiColorRegex)

	// get the indices to all matches of the regular expression for recognizing
	// ANSI color codes, and go then over all runes
	for colind, colindexes := 0, re.FindAllStringIndex(s, -1); idx < len(s); {

		// verify if a ANSI color code starts right at this position
		if colind < len(colindexes) && idx == colindexes[colind][0] {

			// then jump to the first location after the regular expression, and
			// move to the next match of the ANSI color codes
			idx = colindexes[colind][1]
			colind++
		} else {

			// get the rune at the current position
			r, size := utf8.DecodeRuneInString(s[idx:])

			// if this is the rune taking the i-th logical position then return
			// it immediately
			if li == i {
				return r, nil
			}

			// and move forward both physically and logically
			idx += size
			li++
		}
	}

	// if we exited from the main loop then no rune exists at the specified
	// location
	return rune(0), fmt.Errorf("there is no rune at location %v in string '%v'", i, s)
}

// modify the given string by replacing the i-th physical location of the rune by
// the given rune r, i.e., it also counts color ANSI codes
func insertRune(s string, i int, r rune) string {

	var sb strings.Builder

	// safety checking
	if i < 0 || i >= len(s) {
		return s
	}

	// for all runes in the string
	for idx := 0; idx < len(s); {

		// get the rune at the current position
		ir, size := utf8.DecodeRuneInString(s[idx:])

		// if this rune is not the i-th rune then add it to the result
		if idx != i {
			sb.WriteRune(ir)
		} else {

			// otherwise, insert the given rune
			sb.WriteRune(r)
		}

		// and move forward
		idx += size
	}

	// and finally return the string computed so far
	return sb.String()
}

// Insert a single splitter in the physical location (i, jp) of the table that
// has been already drawn using String () that corresponds to the logical
// location (i, jl).
//
// Note that "physical location" is interpreted as follows: i is the i-th slice
// of the textual representation of the table (so that it is both a physical and
// logical coordinate); jl is the j-th *rune* printable+graphic non ANSI color
// code in the string, whereas jl is the j-th *rune* in the string
func addSplitter(tab []string, i, jp, jl int) {

	// define variables for storing the runes to the west, east, north and south
	// of the current location
	var west, east, north, south rune = none, none, none, none

	// west
	if jl > 0 {
		west, _ = getRune(tab[i], jl-1)
	}

	// east
	if jl < countPrintableRuneInString(tab[i])-1 {
		east, _ = getRune(tab[i], jl+1)
	}

	// north
	if i > 0 {
		north, _ = getRune(tab[i-1], jl)
	}

	// south
	if i < len(tab)-1 {
		south, _ = getRune(tab[i+1], jl)
	}

	// now, in case there is a splitter for this combination of west, east,
	// north and south, then insert it and otherwise do nothing
	if splitter := getSingleSplitter(west, east, north, south); splitter != none {

		tab[i] = insertRune(tab[i], jp, splitter)
	}
}

// Add splitters to a table that has been already drawn using String () and
// returns a slice of strings, each representing one line of the table
func addSplitters(tab []string) {

	// store the physical location of a logical position of any string
	var pi int

	// -- Initialization: regular expression used to recognize ANSI color codes
	re := regexp.MustCompile(ansiColorRegex)

	// To do this, the contents of the table are examined (physical) line by
	// line and all positions adjacent to a vertical separator are processed to
	// see whether a splitter has to be added there or not
	for i := 0; i < len(tab); i++ {

		// idx is used to count physical runes, i.e., the physical location of
		// each rune considering also the ANSI color codes, whereas j is the
		// logical location of the physical location idx
		idx, j := 0, 0

		// make a copy of the i-th line of the table
		s := tab[i]

		// get the indices to all matches of the regular expression for
		// recognizing ANSI color codes, and go then over all runes in the given
		// string
		for colind, colindexes := 0, re.FindAllStringIndex(s, -1); idx < len(s); {

			// verify if a ANSI color code starts right at this position
			if colind < len(colindexes) && idx == colindexes[colind][0] {

				// then jump to the first location after the regular expression,
				// and move to the next match of the ANSI color codes
				idx = colindexes[colind][1]
				colind++
			} else {

				// get the rune at the current position
				r, size := utf8.DecodeRuneInString(s[idx:])

				// now, verify whether this is a vertical separator
				if isVerticalSeparator(r) {

					// consider adding a splitter above this location in the
					// physical location (i-1, idx) which maps to the logical
					// location (i-1, j)
					if i > 0 {

						pi, tab[i-1] = logicalToPhysical(tab[i-1], j, true)
						addSplitter(tab, i-1, pi, j)
					}

					// there will be a lot of times when the following statement is
					// just repetitive (i.e., it will be anticipating the work that
					// can be done with the previous if statement when i increases).
					// However, it is necessary for handling some special cases
					// where there is no vertical bar beneath the location of the
					// rune to modify:
					//                      │<---- A
					//                    ━━X━━
					//                      .<---- B
					//
					// in this case, only when being located at A it is possible to
					// substitute the rune at X, whether when being located at B, X
					// will not be invoked if . is any rune other than a vertical
					// separator
					if i <= len(tab)-2 {

						pi, tab[i+1] = logicalToPhysical(tab[i+1], j, true)
						addSplitter(tab, i+1, pi, j)
					}
				}

				// and move forward
				idx += size
				j += 1
			}
		}
	}
}
