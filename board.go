package battleshipsolver

import (
	"fmt"
)

const (
    boardSize = 10
    rowMask uint = 65472    // 1111111111000000
    pegMask uint = 32768    // 1000000000000000
)

type board struct {
    xy [boardSize]uint
    yx [boardSize]uint
}

type probabilities [boardSize][boardSize]int

func newBoard() *board {
    board := &board{
        xy: [boardSize]uint{
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
        },
        yx: [boardSize]uint{
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
        },
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
    for _, i := range b.xy {
        out += fmt.Sprintf("%016b\n", i)
    }
    /*
    out += "\n" 
    for _, i := range b.yx {
        out += fmt.Sprintf("%016b\n", i)
    }
    */

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
    b.xy[location.Row()] = b.xy[location.Row()] ^ uint(mask)
    mask = pegMask>>(location.Row())
    b.yx[location.Column()] = b.yx[location.Column()] ^ uint(mask)
}

func (b *board) merge(b2 *board) {
    for i := range b.xy {
        b.xy[i] = b.xy[i] & b2.xy[i]
    }
    for i := range b.yx {
        b.yx[i] = b.yx[i] & b2.yx[i]
    }
}
