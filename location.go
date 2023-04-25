/*
    Copied from Dan Vittegleo's work at:
    https://github.com/SeaRbSg/battleship-golang/blob/master/battleship/location.go

    Retrieved on April 8, 2023
*/

package battleshipsolver

import (
	"strconv"
	"strings"
)

type (
    Position string
    cell []int
)

type Location struct {
    Position Position
    Row int
    Col int
}

type Locator interface {
    Locate() Location
}

func (c Position) Locate() Location {
	row := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(c[0:1]))
	col, _ := strconv.Atoi(string(c[1:]))
    return Location{
        Position: c,
        Row: row,
        Col: col,
    }
}

func (c cell) Locate() Location {
    const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    row := c[0]
    col := c[1]
    Position := Position(abc[row:row+1] + strconv.Itoa(col+1))
    
    return Location{
        Position: Position,
        Row: row,
        Col: col,
    }
}

func Cell (row int, col int) cell {
    return cell([]int{row, col})
}
