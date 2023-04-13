package battleshipsolver

import (
	"fmt"
	"strconv"
)

const (
    BOARD_SIZE = 10
    INIT_ROW = "1111111111000000"
    MARK_MASK = "1000000000000000"
)

type Board struct {
    xy [BOARD_SIZE]uint
    yx [BOARD_SIZE]uint
}

type Probabilities [BOARD_SIZE][BOARD_SIZE]int

func NewBoard() *Board {
    initRow, _ := strconv.ParseUint(INIT_ROW, 2, 64)
    board := &Board{
        xy: [BOARD_SIZE]uint{
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
        },
        yx: [BOARD_SIZE]uint{
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
            uint(initRow),
        },
    }

    return board
}

func NewProbabilities() *Probabilities {
    probabilities := &Probabilities{
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


func (b *Board) String() string {
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

func (p *Probabilities) String() string {
    var out string
    for _, x := range p {
        for _, y := range x {
            out += fmt.Sprintf("%d ", y)
        }
        out += "\n"
    }

    return out
}

func (b *Board) Mark(location Location) {
    markMask, _ := strconv.ParseUint(MARK_MASK, 2, 64)
    mask := markMask>>(location.Column())
    b.xy[location.Row()] = b.xy[location.Row()] ^ uint(mask)
    mask = markMask>>(location.Row())
    b.yx[location.Column()] = b.yx[location.Column()] ^ uint(mask)
}

func (b *Board) merge(b2 *Board) {
    for i := range b.xy {
        b.xy[i] = b.xy[i] & b2.xy[i]
    }
    for i := range b.yx {
        b.yx[i] = b.yx[i] & b2.yx[i]
    }
}
