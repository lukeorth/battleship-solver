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

type Location string

func (l Location) Row() int {
	row := strings.Index("ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(l[0:1]))
	return row
}

func (l Location) Column() int {
	column, _ := strconv.Atoi(string(l[1:]))
	return (column - 1)
}
