// defs.go
// -*- coding: utf-8 -*-
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const none = 0
const horizontal_empty = ""

const horizontal_blank = ' '
const horizontal_single = '\u2500' // ─
const horizontal_double = '\u2550' // ═
const horizontal_thick = '\u2501'  // ━

const vertical_single = '\u2502' // │
const vertical_double = '\u2551' // ║
const vertical_thick = '\u2503'  // ┃

// Regexps

// the following regexp is used to mach an entire column specification string
const colSpecRegexAll = `^([^clrCLRp]*(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\}))+`

// and the following regexp is used to match the specification of a single
// column
const colSpecRegex = `^[^clrCLRp]*(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})`

// and the following regexp is used to match the specification of a single
// row
const rowSpecRegex = `^[^cbt]*(c|b|t)`

// to extract the format of a single column the following regexp is used
const columnSpecRegex = `(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})`

// in case a paragraph style is used, the following regexp serves to extract the
// numerical argument
const pRegex = `^(C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})$`

// to split strings using the newline as a separator
const newlineRegex = `\n`

// the following regexp is used to start and end ANSI color escape sequences
const ansiColorRegex = `\033([\[;]\d+)+m`

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Table is the main type provided by this package. In order to draw data in
// tabular form it is necessary first to create a table with NewTable. Once a
// Table has been created, it is then possible to use all services provided for
// them
type Table struct {

	// A table consists of a slice of columns, each one with its own
	// specification, and a number of rows where the height of each row is
	// stored. Note that the last separator (if any) is stored as a column
	// without content. They also store the cells of the table as a
	// bidimensional matrix that can be both processed and formatted, i.e., as
	// formatters.
	columns []column
	rows    []row
	cells   [][]formatter
}

// columns do not store contents. A column consists then of a vertical separator
// (to be inserted before its text), their width (number of physical columns),
// and the corresponding styles for showing its contents both horizontally and
// vertically.
type column struct {
	sep              string
	width            int
	hformat, vformat style
}

// rows do not store contents. A row consists then of a number of physical lines
// for displaying its contents
type row struct {
	height int
}

// The style of a cell specifies how to draw it and it is represented typically
// with a string and, additionally, with a numerical value in case a specific
// style (such as 'p') requires it
type style struct {
	alignment byte
	arg       int
}

// Contents are simply strings to be shown on each cell
type content string

// Both splitters (between horizontal and vertical rules) along with other
// surrounding characters, and the rune used as a separator above/below other
// contents are defined as horizontal rules
type hrule string

// Tables can draw cells provided that they can be both processed and formatted:
// cells are first formatted to generate the physical lines required to display
// its contents in the form of formatters, which are then formatted one by one
// to generate a single string which is shown on the table.
//
// The procedure is always the same: for any formatter, it is first "Process"ed
// and each resulting formatter is then "Format"ted. As a result:
//
//		a. All implementation of formatters X shall guarantee that each item in
//		   the output slice []formatter can be casted back into its corresponding
//		   type X, so that they can then be formatted accordingly.
//
//	    b. Tables print directly the result of formatting each item in the result
//	       of the processing step
type formatter interface {

	// Processing a cell means transforming logical rows into physical ones by
	// splitting its contents across several (physical) rows, and also adding
	// blank lines so that the result satisfies the vertical format of the
	// column where it has to be shown, if and only if the height of the
	// corresponding row is larger than the number of physical rows necessary to
	// display the contents of the cell. To properly process a cell it is
	// necessary to get a pointer to the table, and also the integer indices to
	// the row and column of the cell
	Process(t *Table, irow, jcol int) []formatter

	// Cells are also formatted (physical) line by line where each physical line
	// is the result of processing cell (irow, jcol) and should be given in the
	// receiver of this method. Each invocation returns a string where each
	// (physical) line is forrmatted according to the horizontal format
	// specification of the j-th column.
	Format(t *Table, irow, jcol int) string
}

// Multicells, defined next, are a general case of multicolumns and multirows.
// Muluticolumns are defined as multicells with just one logical row, whereas
// the latter are defined as multicells with just a single logical column. The
// following enum is created to distinguish the different types of multicells
// ---note the user can specifically create multicells (as such) with one column
// and/or one row just to selectively update the horizontal/vertical format of a
// cell!
type multicellType int

const (
	multicolumn_t multicellType = iota
	multirow_t
	multicell_t
)

// Multicells are used to merge several cells (along rows and/or columns) into
// one single cell. They are essentially tables with an arbitrary number of
// columns and rows whose specification is given by the user.
//
// They are formatters, i.e., they can be inserted into a table to merge
// nbcolumns columns and nbrows rows from an initial row/column under a
// different format explicitly given by the user as a row and column
// specifications that are processed to produce a table which is filled in with
// data from a number of arguments.
//
// Multicells allow the row and column specification to contain a last
// separator. If a last separator is given in the column specification, it is
// used as the separator of the next cell; if one is given in the row
// specification, it is then used as the horizontal rule of the next row.
type multicell struct {
	mtype              multicellType // type of multicell
	jinit, nbcolumns   int           // initial column and # columns
	iinit, nbrows      int           // initial row and # rows
	cspec, rspec       string        // column and row specification
	clastsep, rlastsep string        // column and row last separator
	table              Table         // table used to render its contents
	args               []any         // arguments given to the multicell
	output             string        // rendered contents of the multicell
}

// The splitter is defined as an association of four different runes: the west,
// east, north and south runes of the splitter:
//
//		    north
//	          |
//
// west - splitter - east
//
//	          |
//		     south
//
// this creates an association of ASCII/UTF-8 characters defined by hand below.
// For example, with the four runes shown in the diagram above, the splitter
// would be a single cross ┼ (\u253c). Importantly, some of these four runes
// might be empty and the associations given in the map below are then used to
// draw different types of corners.
//
// note that some combinations below are commented out. This is important as
// those combinations which are not recognized are then properly substituted by
// the algorithm.
//
// 'none' has to be recognized at every level of the tree. The reason is that if
// a query is performed to this nested map with a character which is not found,
// then none is used, so that it must exist as a key in all entries
var splitterUTF8 = map[rune]map[rune]map[rune]map[rune]rune{

	none: {
		none: {
			none: {
				none: none,
				// vertical_single: '\u2577', // ╷
				// vertical_double: '\u257b', // ╻: double half down not supported!
				// vertical_thick:  '\u257b', // ╻
			},
			vertical_single: {
				// none:            '\u2575', // ╵,
				vertical_single: vertical_single,
				vertical_double: none,
				vertical_thick:  none,
			},
			vertical_double: {
				// none:            '\u2579', // ╹: double half down not supported!
				vertical_single: none,
				vertical_double: vertical_double,
				vertical_thick:  none,
			},
			vertical_thick: {
				// none:            '\u2579', // ╹
				vertical_single: none,
				vertical_double: none,
				vertical_thick:  vertical_thick,
			},
		},
		horizontal_single: {
			none: {
				none:            none,
				vertical_single: '\u250c', // ┌
				vertical_double: '\u2553', // ╓
				vertical_thick:  '\u250e', // ┎
			},
			vertical_single: {
				none:            '\u2514', // └
				vertical_single: '\u251c', // ├
				vertical_double: '\u251f', // ┟: south double not supported!
				vertical_thick:  '\u251f', // ┟: south double not supported!
			},
			vertical_double: {
				none:            '\u2559', // ╙
				vertical_single: '\u251e', // ┞: north double not supported!
				vertical_double: '\u255f', // ╟
				vertical_thick:  '\u2520', // ┠: north double not supported!
			},
			vertical_thick: {
				none:            '\u2516', // ┖
				vertical_single: '\u251e', // ┞
				vertical_double: '\u2520', // ┠: south double not supported!
				vertical_thick:  '\u2520', // ┠
			},
		},
		horizontal_double: {
			none: {
				none:            none,
				vertical_single: '\u2552', // ╒
				vertical_double: '\u2554', // ╔
				vertical_thick:  '\u250f', // ┏: east double not supported!
			},
			vertical_single: {
				none:            '\u2558', // ╘
				vertical_single: '\u255e', // ╞
				vertical_double: '\u2522', // ┢: east/south double not supported!
				vertical_thick:  '\u2522', // ┢: east double not supported!
			},
			vertical_double: {
				none:            '\u255a', // ╚
				vertical_single: '\u2521', // ┡: east/north double not supported!
				vertical_double: '\u2560', // ╠
				vertical_thick:  '\u2523', // ┣: east/north double not supported!
			},
			vertical_thick: {
				none:            '\u2517', // ┗: easth double not supported!
				vertical_single: '\u2521', // ┡: east double not supported!
				vertical_double: '\u2523', // ┣: east/south double not supported!
				vertical_thick:  '\u2523', // ┣: east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            none,
				vertical_single: '\u250d', // ┍
				vertical_double: '\u250f', // ┏: south double not supported!
				vertical_thick:  '\u250f', // ┏: south double not supported!
			},
			vertical_single: {
				none:            '\u2515', // ┕
				vertical_single: '\u251d', // ┝
				vertical_double: '\u2522', // ┢: south double not supported!
				vertical_thick:  '\u2522', // ┢
			},
			vertical_double: {
				none:            '\u2517', // ┗: north double not supported
				vertical_single: '\u2521', // ┡: north double not supported!
				vertical_double: '\u2523', // ┣: north/south double not supported!
				vertical_thick:  '\u2523', // ┣: north double not supported!
			},
			vertical_thick: {
				none:            '\u2517', // ┗
				vertical_single: '\u2521', // ┡
				vertical_double: '\u2523', // ┣: south double not supported
				vertical_thick:  '\u2523', // ┣
			},
		},
	},

	horizontal_single: {
		none: {
			none: {
				none:            none,
				vertical_single: '\u2510', // ┐
				vertical_double: '\u2556', // ╖
				vertical_thick:  '\u2512', // ┒
			},
			vertical_single: {
				none:            '\u2518', // ┘
				vertical_single: '\u2524', // ┤
				vertical_double: '\u2527', // ┧: south double not supported!
				vertical_thick:  '\u2527', // ┧: south double not supported!
			},
			vertical_double: {
				none:            '\u255c', // ╜
				vertical_single: '\u2526', // ┦: north double not supported!
				vertical_double: '\u2562', // ╢:
				vertical_thick:  '\u2528', // ┨: north double not supported!
			},
			vertical_thick: {
				none:            '\u251a', // ┚
				vertical_single: '\u2526', // ┦
				vertical_double: '\u2528', // ┨: south double not supported!
				vertical_thick:  '\u2528', // ┨
			},
		},
		horizontal_single: {
			none: {
				none:            none,
				vertical_single: '\u252c', // ┬
				vertical_double: '\u2565', // ╥
				vertical_thick:  '\u2530', // ┰
			},
			vertical_single: {
				none:            '\u2534', // ┴
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_double: {
				none:            '\u2568', // (cannot be shown on Emacs :( )
				vertical_single: '\u2540', // ╀: north double not supported!
				vertical_double: '\u256b', // ╫
				vertical_thick:  '\u2542', // ╂: north double not supported!
			},
			vertical_thick: {
				none:            '\u2538', // ┸
				vertical_single: '\u2540', // ╀
				vertical_double: '\u2542', // ╂: south double not supported!
				vertical_thick:  '\u2542', // ╂
			},
		},
		horizontal_double: {
			none: {
				none:            none,
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_single: {
				none:            '\u2536', // ┶: east double not supported!
				vertical_single: '\u253e', // ┾: east dobule not supported!
				vertical_double: '\u2546', // ╆: east/south double not supported!
				vertical_thick:  '\u2546', // ╆: east double not supported!
			},
			vertical_double: {
				none:            '\u253a', // ┺: easth double not supported!
				vertical_single: '\u2544', // ╄: east/north double not supported!
				vertical_double: '\u254a', // ╊: east/north/south double not supported!
				vertical_thick:  '\u254a', // ╊: east/north double not supported!
			},
			vertical_thick: {
				none:            '\u253a', // ┺: easth double not supported!
				vertical_single: '\u2544', // ╄: east double not supported!
				vertical_double: '\u254a', // ╊: east/south double not supported!
				vertical_thick:  '\u254a', // ╊: east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            none,
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_single: {
				none:            '\u2536', // ┶
				vertical_single: '\u253e', // ┾
				vertical_double: '\u2546', // ╆: south double not supported!
				vertical_thick:  '\u2546', // ╆
			},
			vertical_double: {
				none:            '\u2536', // ┶: north double not supported
				vertical_single: '\u2544', // ╄: north double not supported!
				vertical_double: '\u254a', // ╊: north/south double not supported!
				vertical_thick:  '\u254a', // ╊: north double not supported!
			},
			vertical_thick: {
				none:            '\u253a', // ┺
				vertical_single: '\u2544', // ╄
				vertical_double: '\u254a', // ╊: south double not supported
				vertical_thick:  '\u254a', // ╊
			},
		},
	},

	horizontal_double: {
		none: {
			none: {
				none:            none,
				vertical_single: '\u2555', // ╕
				vertical_double: '\u2557', // ╗
				vertical_thick:  '\u2513', // ┓: west double not supported!
			},
			vertical_single: {
				none:            '\u255b', // ╛
				vertical_single: '\u2561', // ╡
				vertical_double: '\u252a', // ┪: west/south double not supported!
				vertical_thick:  '\u252a', // ┪: west double not supported!
			},
			vertical_double: {
				none:            '\u255d', // ╝
				vertical_single: '\u2529', // ┩: west/north double not supported!
				vertical_double: '\u2563', // ╣
				vertical_thick:  '\u252b', // ┫: west/north double not supported!
			},
			vertical_thick: {
				none:            '\u251b', // ┛: west double not supported!
				vertical_single: '\u2529', // ┩: west double not supported!
				vertical_double: '\u252a', // ┪: west/south double not supported!
				vertical_thick:  '\u252b', // ┫: west double not supported!
			},
		},
		horizontal_single: {
			none: {
				none:            none,
				vertical_single: '\u252d', // ┭: west double not supported!
				vertical_double: '\u2531', // ┱: west double not supported!
				vertical_thick:  '\u2531', // ┱: west double not supported!
			},
			vertical_single: {
				none:            '\u2535', // ┵: west double not supported!
				vertical_single: '\u253d', // ┽: west double not supported!
				vertical_double: '\u2545', // ╅: west/south double not supported!
				vertical_thick:  '\u2545', // ╅: west double not supported!
			},
			vertical_double: {
				none:            '\u2539', // ┹: west/north double not supported!
				vertical_single: '\u2543', // ╃: west/north double not supported!
				vertical_double: '\u2549', // ╉: west/north/south double not supported!
				vertical_thick:  '\u2549', // ╉: west/north double not supported!
			},
			vertical_thick: {
				none:            '\u2539', // ┹: west double not supported!
				vertical_single: '\u2543', // ╃: west double not supported!
				vertical_double: '\u2549', // ╉: west/south double not supported!
				vertical_thick:  '\u2549', // ╉: west double not supported!
			},
		},
		horizontal_double: {
			none: {
				none:            none,
				vertical_single: '\u2564', // ╤
				vertical_double: '\u2566', // ╦
				vertical_thick:  '\u2533', // ┳: west/east double not supported!
			},
			vertical_single: {
				none:            '\u2567', // (cannot be shown on Emacs :( )
				vertical_single: '\u256a', // ╪
				vertical_double: '\u2548', // ╈: west/east/south double not supported!
				vertical_thick:  '\u2548', // ╆: west/east double not supported!
			},
			vertical_double: {
				none:            '\u2569', // ╩
				vertical_single: '\u2547', // ╇: west/east/north double not supported!
				vertical_double: '\u256c', // ╬
				vertical_thick:  '\u254b', // ╋: west/east/north double not supported!
			},
			vertical_thick: {
				none:            '\u253b', // ┻: west/east double not supported!
				vertical_single: '\u2547', // ╇: west/east double not supported!
				vertical_double: '\u254b', // ╋: west/east/south double not supported!
				vertical_thick:  '\u254b', // ╋: west/east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            none,
				vertical_single: '\u252f', // ┯: west double not supported!
				vertical_double: '\u2533', // ┳: west/south double not supported
				vertical_thick:  '\u2533', // ┳: west double not supported!
			},
			vertical_single: {
				none:            '\u2537', // ┷: west double not supported!
				vertical_single: '\u253f', // ┿: west double not supported!
				vertical_double: '\u2548', // ╈: west/south double not supported!
				vertical_thick:  '\u2548', // ╆: west double not supported!
			},
			vertical_double: {
				none:            '\u253b', // ┻: west/north double not supported!
				vertical_single: '\u2547', // ╇: west/north double not supported!
				vertical_double: '\u254b', // ╋: west/north/south double not supported
				vertical_thick:  '\u254b', // ╋: west/north double not supported!
			},
			vertical_thick: {
				none:            '\u253b', // ┻: west double not supported!
				vertical_single: '\u2547', // ╇: west double not supported!
				vertical_double: '\u254b', // ╋: west/south double not supported!
				vertical_thick:  '\u254b', // ╋: west double not supported!
			},
		},
	},

	horizontal_thick: {
		none: {
			none: {
				none:            none,
				vertical_single: '\u2511', // ┑
				vertical_double: '\u2513', // ┓: south double not supported!
				vertical_thick:  '\u2513', // ┓: west double not supported!
			},
			vertical_single: {
				none:            '\u2519', // ┙
				vertical_single: '\u2525', // ┥
				vertical_double: '\u252a', // ┪: south double not supported!
				vertical_thick:  '\u252a', // ┪
			},
			vertical_double: {
				none:            '\u251b', // ┛: north double not supported!
				vertical_single: '\u2529', // ┩: north double not supported!
				vertical_double: '\u252b', // ┫: north/south double not supported!
				vertical_thick:  '\u252b', // ┫: west/north double not supported!
			},
			vertical_thick: {
				none:            '\u251b', // ┛
				vertical_single: '\u2529', // ┩
				vertical_double: '\u252b', // ┫: south double not supported!
				vertical_thick:  '\u252b', // ┫
			},
		},
		horizontal_single: {
			none: {
				none:            none,
				vertical_single: '\u252d', // ┭
				vertical_double: '\u2531', // ┱: south double not supported!
				vertical_thick:  '\u2531', // ┱
			},
			vertical_single: {
				none:            '\u2535', // ┵
				vertical_single: '\u253d', // ┽
				vertical_double: '\u2545', // ╅
				vertical_thick:  '\u2545', // ╅
			},
			vertical_double: {
				none:            '\u2539', // ┹: north double not supported!
				vertical_single: '\u2543', // ╃: north double not supported!
				vertical_double: '\u2549', // ╉: north/south double not supported!
				vertical_thick:  '\u2549', // ╉: north double not supported!
			},
			vertical_thick: {
				none:            '\u2539', // ┹
				vertical_single: '\u2543', // ╃
				vertical_double: '\u2549', // ╉: south double not supported!
				vertical_thick:  '\u2549', // ╉
			},
		},
		horizontal_double: {
			none: {
				none:            none,
				vertical_single: '\u252f', // ┯: east double not supported!
				vertical_double: '\u2533', // ┳: east/south double not supported!
				vertical_thick:  '\u2533', // ┳: east double not supported!
			},
			vertical_single: {
				none:            '\u2537', // ┷: east double not supported!
				vertical_single: '\u253f', // ┿: east double not supported!
				vertical_double: '\u2548', // ╈: east/south double not supported!
				vertical_thick:  '\u2548', // ╆: east double not supported!
			},
			vertical_double: {
				none:            '\u253b', // ┻: east/north double not supported!
				vertical_single: '\u2547', // ╇: east/north double not supported!
				vertical_double: '\u254b', // ╋: east/north/south double not supported!
				vertical_thick:  '\u254b', // ╋: east/north double not supported!
			},
			vertical_thick: {
				none:            '\u253b', // ┻: east double not supported!
				vertical_single: '\u2547', // ╇: east double not supported!
				vertical_double: '\u254b', // ╋: east/south double not supported!
				vertical_thick:  '\u254b', // ╋: east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            horizontal_thick,
				vertical_single: '\u252f', // ┯
				vertical_double: '\u2533', // ┳: south double not supported
				vertical_thick:  '\u2533', // ┳
			},
			vertical_single: {
				none:            '\u2537', // ┷
				vertical_single: '\u253f', // ┿
				vertical_double: '\u2548', // ╈: south double not supported!
				vertical_thick:  '\u2548', // ╆
			},
			vertical_double: {
				none:            '\u253b', // ┻: north double not supported!
				vertical_single: '\u2547', // ╇: north double not supported!
				vertical_double: '\u254b', // ╋: north/south double not supported
				vertical_thick:  '\u254b', // ╋: north double not supported!
			},
			vertical_thick: {
				none:            '\u253b', // ┻
				vertical_single: '\u2547', // ╇
				vertical_double: '\u254b', // ╋: south double not supported!
				vertical_thick:  '\u254b', // ╋
			},
		},
	},
}
