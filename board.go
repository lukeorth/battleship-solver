package battleshipsolver

import (
	"fmt"
)

const (
    boardSize = 10
    rowMask uint = 1023 // 1111111111
    pegMask uint = 512  // 1000000000
)

type board [boardSize]uint

type probabilities [boardSize][boardSize]int

func newBoard() *board {
    board := &board{
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
        rowMask,
    }
    return board
}

func newProbabilities() *probabilities {
    probabilities := &probabilities{
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
    }

    return probabilities
}

func (b *board) String() string {
    var out string
    for _, i := range b {
        out += fmt.Sprintf("%010b\n", i)
    }
    return out
}

func (p *probabilities) String() string {
    var out string
    for _, x := range p {
        for _, y := range x {
            out += fmt.Sprintf("%2d ", y)
        }
        out += "\n"
    }
    return out
}

func (b *board) mark(location Location) {
    mask := pegMask>>(location.Column())
    b[location.Row()] = b[location.Row()] ^ uint(mask)
}

func (b *board) condenseRows(rowStart int, count int) uint {
    rowCopy := b[rowStart]
    for i := 0; i < count; i++ {
        rowCopy &= b[rowStart+i]
    }
    return rowCopy
}

func (b *board) merge(b2 *board) {
    for i := range b {
        b[i] &= b2[i]
    }
}
