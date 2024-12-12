# table

This package implements means for drawing data in tabular form and it is
intended as a substitution of the Go standard package `tabwriter`. Its design is
based on the functionality of tables in LaTeX but extends its functionality in
various ways through a very simple interface

It honours [UTF-8 characters](https://www.utf8-chartable.de/), [ANSI color escape sequences](https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences), fixed- and variable-width columns, full/partial
horizontal rules, various vertical and horizontal alignment options, and
multicolumns.

Remarkably, it prints any *stringer* and as tables are stringers, tables can be
nested to any degree.

While I do not want to delve into the [controversy between tabs and
spaces](https://www.youtube.com/watch?v=SsoOG6ZeyUI) this is my position: "*tabs
should be used for indentation while spaces should be preferred for alignment*".
Thus, this package uses spaces for aligning contents and it should then be used
with fixed-size fonts.


# Installation 

Clone and install the `table` package with the following command:

    $ go get github.com/clinaresl/table
    
To try the different examples given in the package change dir to
`$GOPATH/github.com/clinaresl/table` and type:

    $ go test
   
# Usage #

This section provides various examples of usage with the hope of providing a
flavour of the different capabilities of the package. For a full description of
the package check out the technical documentation.

## First step: Create a table ##

Before inserting data to a new table it is necessary to create it first:

```Go
	t, err := NewTable("l   p{25}")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
```

This snippet creates a table with two columns. The first one displays its
contents ragged left, whereas the second one takes a fixed width of 25
characters to display the contents of each cell and, in case a cell exceeds the
available width, its contents are shown left ragged in as many lines as needed
---i.e., `table` supports multi-line cells. In case it was not possible to
successfully process the *column specification*, an error is immediately
returned.

The following table shows all the different options for specifying the format of
a single column:

| Syntax | Purpose |
|:------:|:-------:|
|  `l`   | the contents of the column are ragged left |
|  `c`   | the contents of the column are horizontally aligned |
|  `r`   | the contents of the column are ragged right |
|  `p{NUMBER}` | the cell takes a fixed with equal to *NUMBER* characters and the contents are split across various lines if needed |
|  `L{NUMBER}` | the width of the column does not exceed *NUMBER* characters and the contents are ragged left |
|  `C{NUMBER}` | the width of the column does not exceed *NUMBER* characters and the contents are centered |
|  `R{NUMBER}` | the width of the column does not exceed *NUMBER* characters and the contents are ragged right |

The *column specification* allows the usage of `|`, e.g.:

``` Go
	t, _ := NewTable("|c|c|c|c|c|")
```

creates a table with five different columns all separated by a single vertical
separator. It is possible to create also *double* and *thick* vertical
separators using `||` and `|||` respectively. As a matter of fact, these are
just shortcuts and the UTF-8 characters `│`, `║` and `┃` can be used
respectively. It is also possible to provide any other character (e.g., blank
spaces) either before or after any column. These are then copied either before
or after the contents of each cell in each row.

In case a second string is given to `NewTable` it is interpreted as the *row
specification*:

```Go
	t, _ := NewTable("| c | c  c |", "cct")
```

This line (where no error checking is performed!) creates three different
columns whose contents are horizontally centered surrounded by a single space
and with vertical single separators between adjacent columns and before and
after the first and last column. In addition, it sets the *vertical alignment*
of each cell as follows: the contents of the first and second columns are
vertically centered (`c`), whereas the contents of the last column are pushed to
the top of the cell ---`t`. The modifiers available to be used in the *row
specification* are shown next:

| Syntax | Purpose |
|:------:|:-------:|
|  `t`   | the contents of the column are aligned to the top |
|  `c`   | the contents of the column are vertically aligned |
|  `b`   | the contents of the column are aligned to the bottom |

By default, all columns are vertically aligned to the top. In case a *row
specification* is given it must refer to as many columns as there are in the
*column specification* given first or less. In contraposition to the *column
specification*, the *row specification* can only consist of any of the modifiers
shown above.
   
`NewTable` returns a pointer to `Table` which can be used next for adding data
to it and, in the end, printing it.

## Second step: Adding rows ##

`table` acknowledges two different types of rows either horizontal rules or
lines of data.

### Adding horizontal rules ###

There are three different services for adding horizontal rules anywhere in a
table:

```Go
   func (t *Table) AddSingleRule(cols ...int) error
   func (t *Table) AddDoubleRule(cols ...int) error
   func (t *Table) AddThickRule(cols ...int) error
```

When invoked with no arguments they just show a full horizontal rule spanning
over all columns of the table. Single rules are shown with the UTF-8 character
`─`; double rules are drawn using `═`, and thick rules use `━`.

If they are invoked with arguments, then these are taken in pairs, each pair
standing for a *starting* and *ending* column numbered from 0, so that
horizontal rules are drawn only over those columns in the given range.

In case it is not possible to process the given arguments then an informative
error is returned.

A couple of examples follow:

``` Go
    t, _ := NewTable("|c|c|c|c|c|")
    t.AddThickRule ()
    t.AddSingleRule(0, 1, 2, 3, 4, 5)
```

### Adding data ###

Data is added to the bottom of a table with:

``` Go
    func (t *Table) AddRow(cells ...any) error
```

It accepts an arbitrary number of *any* arguments and adds the result of the
`Sprintf` operation of each argument to each cell of the last row of the table.
If the number of arguments is strictly less than the number of columns given in
the *column specification* then the remaining cells are left empty. Thus, if no
argument is provided, an empty line of data is generated. However, if the number
of arguments given is strictly larger than the number of columns of the table an
error is returned.

The following example adds data to a table with three columns: 

``` Go
	t, err := NewTable("| c || c ||| c |")
	err = t.AddRow("Year\n1979", "Year\n2013", "Year\n2018")
	if err != nil {
		log.Fatalln(" AddRow: Fatal error!")
	}
	err = t.AddRow("Ariane", "Gaia\nProba Series\nSwarm", "Aeolus\nBepicolombo\nMetop Series")
	if err != nil {
		log.Fatalln(" AddRow: Fatal error!")
	}
```

Note that the contents of any cell can contain any newline characters `\n`. If
so, the text is split in as many lines as needed, i.e., `table` supports
multi-line cells.

## Third step: Printing tables ##

The last step consists of printing the contents of any table. By definition,
tables are stringers and thus, all that is required is just to print the
contents with a `Print`-like function:

``` Go
	t, _ := NewTable("l | r ")
	t.AddThickRule()
	t.AddRow("Country", "Population")
	t.AddSingleRule()
	t.AddRow("China", "1,394,015,977")
	t.AddRow("India", "1,326,093,247")
	t.AddRow("United States", "329,877,505")
	t.AddRow("Indonesia", "267,026,366")
	t.AddRow("Pakistan", "233,500,636")
	t.AddRow("Nigeria", "214,028,302")
	t.AddThickRule()
	fmt.Printf("%v", t)
```

Which produces the result shown next (all examples are shown as images to avoid
your browser to show unrealistic renderings as a result of your preferences):

![example-1](figs/example-1.png "example-1")


# Gotchas #

Beyond the basic usage of tables, `table` provides other features which are
described next

## ANSI color codes ##

Of course, `table` fully supports UTF-8 encoded characters, but it also manages
[ANSI color escape
sequences](https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences)
provided that they are supported by your terminal. 

There are several Go packages that can actually produce the additional
characters required to show the output in various forms, but this implementation
is not tied to any, so that the following examples explicitly show the specific
ANSI color codes required to render each fragment:

``` Go
	t, err := NewTable("r l c c c l c")
	if err != nil {
		log.Fatalf(" NewTable: Fatal error (%v)", err)
	}
	t.AddRow("\033[36;3;4mID", "Age", "Project", "Tags", "Due", "Description", "Urg \033[0m")
	t.AddRow(1, "8mo", "personal.programming.go", "program", "\033[33;1m2022-10-21\033[0m", "Document table", 15.3)
	t.AddRow(2, "3mo", "gii.cag", "video", "\033[33;1m2022-04-11\033[0m", "Create a video promoting UC3M", 14.4)
	t.AddRow(3, "6w", "research.editorial.review", "aicomm", "\033[33;1m2023-05-30\033[0m", "Review the latest papers", 14.1)

	fmt.Printf("%v", t)
```

which produces:

![example-2](figs/example-2.png "example-2")

Mind the trick! The table contains no horizontal rule and the same effect is
created by underlining the header as in [task warrior](https://taskwarrior.org/)
---as a matter of fact, the example resembles the output of the command `task
list` of task warrior. 

Many other combinations and tricks are possible for improving presentations in
tabular form. For example, the following snippet shows how to colour the
vertical separators and horizontal rules also:

``` Go
	t, err := NewTable("\033[38;2;160;10;10m| c \033[38;2;10;160;10m| c \033[38;2;80;80;160m| c \033[38;2;160;80;40m|\033[0m", "cb")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
	t.AddRow("\033[38;2;206;10;0mPlayer\033[0m", "\033[38;2;10;206;0mYear\033[0m", "\033[38;2;100;0;206mTournament\033[0m")
	t.AddSingleRule()
	t.AddRow("\033[38;5;206mRafa Nadal\033[0m", "2010", "French Open\nWimbledon\nUS Open")
	t.AddSingleRule()
	t.AddRow("Roger Federer", "2007", "Australian Open\nWimbledon\nUS Open")
	t.AddSingleRule()

	fmt.Printf("%v", t)	
```

which is rendered as follows:

![example-3](figs/example-3.png "example-3")


Mind (again) the trick! The ANSI color codes of each line including the headers
are automatically ended with `\033[0m` just simply by adding it to the *column
specification* of the table. Of course, one could end each line manually but as
the example shows this is not necessary at all.

## Multicolumns ##

Multicolumns are defined as ordinary cells which span over several columns in
the same row. Because `table` creates column-orientated tables, it is also
possible to substitute an arbitrary number of columns in the table by a
different number of columns with a different format or, in other words,
multicolumns can be used both for *merging* or *splitting* columns. For example:


``` Go
	t, _ := NewTable("l c c || c c")
	t.AddRow(Multicolumn(5, "c", "Table 2: Overall Results"))
	t.AddThickRule()
	t.AddRow("", Multicolumn(2, "c", "Females"), Multicolumn(2, "c", "Males"))
	t.AddSingleRule(1, 5)
	t.AddRow("Treatment", "Mortality", "Mean\nPressure", "Mortality", "Mean\nPressure")
	t.AddSingleRule()
	t.AddRow("Placebo", 0.21, 163, 0.22, 164)
	t.AddRow("ACE Inhibitor", 0.13, 142, 0.15, 144)
	t.AddRow("Hydralazine", 0.17, 143, 0.16, 140)
	t.AddThickRule()
	t.AddRow(Multicolumn(5, "c", "Adapted from\nhttps://tex.stackexchange.com/questions/314025/making-stats-table-with-multicolumn-and-cline"))
	t.AddSingleRule()
	fmt.Printf("%v", t)
```

results in the following table:

![example-4](figs/example-4.png "example-4")

Note that multicolumns are created with the function `Multicolumn` which expects
first, the number of columns it has to take; their format which has to be given
according to the rules discussed in [First step: Create a table](#usage); and
finally, the contents to be shown in the multicolumn. Because `Multicolumn`
accepts any valid column specification in its second argument, `Multicolumn`
serves to various purposes:

1. *Merging* columns: the previous example shows how to *merge* different
   columns into one. The table is originally defined with 5 different columns so
   that the first, second and last row merge a number of columns into one whose
   specification is given in the definition of the multicolumn.

2. *Splitting* columns: the following example shows how to use `Multicolumn` to
   create a single multicolumn which, however, consists of various columns,
   effectively splitting the original column into others:
   
``` Go
	t, _ := NewTable("r|c|")
	t.AddThickRule(1, 5)
	t.AddRow("", Multicolumn(1, "|c|c|c|c|c|", "Mon", "Tue", "Wed", "Thu", "Fri"))
	t.AddSingleRule()
	t.AddRow("9:00", "Enter school")
	t.AddSingleRule()
	t.AddRow("13:00", "Lunch")
	t.AddSingleRule()
	t.AddRow("14:30", "More classes")
	t.AddSingleRule()
	t.AddRow("17:00",
		Multicolumn(1, "|c|c|", "Basketball", "Guitar"))
	t.AddSingleRule()
	t.AddRow("18:30", "Go home!")
	t.AddThickRule()
	fmt.Printf("%v", t)
```
  
  which produces:
  
![example-5](figs/example-5.png "example-5")  

3. Also to selectively *modify the appearance* of the table at selected points.
   In the following example, the boxes shown in the middle and the bottom are
   created using multicolumns of width 1. In fact, these multicolumns are used
   just for modifying the vertical separators so that the boxes are correctly
   drawn. Much the same happens with the multicolumn created for showing the
   description of our planet below the thick rule: This line is created with a
   multicolumn of width 3 which also modifies the column specification to
   `C{30}` so that it actually takes several lines ---note in passing that ANSI
   color escape sequences are used here to show the text slanted.

``` Go
	t, _ := NewTable("    r   l c")
	t.AddRow(Multicolumn(3, "    c", "♁ Earth"))
	t.AddThickRule()
	t.AddRow(Multicolumn(3, "    C{30}", "\033[37;3mEarth is the third planet from the Sun and the only astronomical object known to harbor life\033[0m"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(1, "   |c", "Feature"),
		Multicolumn(1, "   c", "Measure"),
		Multicolumn(1, "c|", "Unit"))
	t.AddSingleRule()
	t.AddRow("Aphelion", 152100000, "km")
	t.AddRow("Perihelion", 147095000, "km")
	t.AddRow("Eccentricity", 0.0167086)
	t.AddRow("Orbital period", 365.256363004)
	t.AddRow("Semi-major axis", 149598023, "km")
	t.AddSingleRule()
	t.AddRow(Multicolumn(3, "   │c│", "\033[37;3mData provided by Wikipedia\033[0m"))
	t.AddSingleRule()
	fmt.Printf("%v", t)
```

  which yields the following results:

![example-6](figs/example-6.png "example-6")

 
  Other than this, this example shows also that tables can be indented by adding
  the same text (e.g., blanks) to the beginning of each row.

## Multirows ##

Multirows are defined analogously to multicolumns, i.e., as ordinary cells which
span over an arbitrary number of rows in the same column. An important
difference with multicolumns though is that multirows *merge* several lines into
one, but they do not provide means for *splitting* a specific line into others.
A couple of usages follow:



## Nested tables ##

`table` prints *stringers* and because tables as created by this package are
also *stringers*, they can then be nested to any degree:

``` Go
	board1, _ := NewTable("||cccccccc||")
	board1.AddDoubleRule()
	board1.AddRow("\u265c", "\u265e", "\u265d", "\u265b", "\u265a", "\u265d", "", "\u265c")
	board1.AddRow("\u265f", "\u265f", "\u265f", "\u265f", "\u2592", "\u265f", "\u265f", "\u265f")
	board1.AddRow("", "\u2592", "", "\u2592", "", "\u265e", "", "\u2592")
	board1.AddRow("\u2592", "", "\u2592", "", "\u265f", "", "\u2592", "")
	board1.AddRow("", "\u2592", "", "\u2592", "\u2659", "\u2659", "", "\u2592")
	board1.AddRow("\u2592", "", "\u2658", "", "\u2592", "", "\u2592", "")
	board1.AddRow("\u2659", "\u2659", "\u2659", "\u2659", "", "\u2592", "\u2659", "\u2659")
	board1.AddRow("\u2656", "", "\u2657", "\u2655", "\u2654", "\u2657", "\u2658", "\u2656")
	board1.AddDoubleRule()

	board2, _ := NewTable("||cccccccc||")
	board2.AddDoubleRule()
	board2.AddRow("\u265c", "\u265e", "\u265d", "\u265b", "\u265a", "\u265d", "\u265e", "\u265c")
	board2.AddRow("\u265f", "\u265f", "\u265f", "", "\u265f", "\u265f", "\u265f", "\u265f")
	board2.AddRow("", "\u2592", "", "\u2592", "", "\u2592", "", "\u2592")
	board2.AddRow("\u2592", "", "\u2592", "", "\u2592", "", "\u2592", "")
	board2.AddRow("", "\u2592", "", "\u2659", "\u265f", "\u2592", "", "\u2592")
	board2.AddRow("\u2592", "", "\u2592", "", "\u2592", "\u2659", "\u2592", "")
	board2.AddRow("\u2659", "\u2659", "\u2659", "\u2592", "", "\u2592", "\u2659", "\u2659")
	board2.AddRow("\u2656", "\u2658", "\u2657", "\u2655", "\u2654", "\u2657", "\u2658", "\u2656")
	board2.AddDoubleRule()

	t, _ := NewTable("| c | c  c |", "cct")
	t.AddSingleRule()
	t.AddRow("ECO Code", "Moves", "Board")
	t.AddSingleRule()
	t.AddRow("C26 Vienna Game: Vienna Gambit", "1.e4 e5 2.♘c3 ♞6 3.f4", board1)
	t.AddRow("D00 Blackmar-Diemer Gambit: Gedult Gambit", "1.e4 d5 2.d4 exd4 3.f3", board2)
	t.AddSingleRule()

	fmt.Printf("%v", t)	
```

Both chess boards are tables so that the last table, named `t` just simply adds
them to each row:

![example-7](figs/example-7.png "example-7")


# License #

MIT License

Copyright (c) 2023, Carlos Linares López

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>  
Computer Science Department <https://www.inf.uc3m.es/en>  
Universidad Carlos III de Madrid <https://www.uc3m.es/home>
