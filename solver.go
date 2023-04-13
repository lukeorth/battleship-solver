package battleshipsolver

import (
	"fmt"
	"strconv"
)


const (
    HUNT_MODE = "hunt"
    TARGET_MODE = "target"
)

type Solver struct {
    Probabilities *Probabilities
    BestMove Location
    mode string
    Fleet *Fleet
    HuntBoard *Board
    TargetBoard *Board
}

func NewSolver() *Solver {
    solver := &Solver{
        Probabilities: NewProbabilities(),
        Fleet: BuildFleet(),
        HuntBoard: NewBoard(),
        TargetBoard: NewBoard(),
        mode: HUNT_MODE,
    }
    return solver
}

func (s *Solver) Hit(location Location) {
    s.mode = TARGET_MODE
    s.TargetBoard.Mark(location)
}

func (s *Solver) Miss(location Location) {
    s.mode = HUNT_MODE
    s.HuntBoard.Mark(location)
}

func (s *Solver) HitAndSunk(location Location, ship string) {
    s.mode = HUNT_MODE
    s.TargetBoard.Mark(location)
    s.HuntBoard.merge(s.TargetBoard)
    s.Fleet.SinkShip(ship)
}

func (s *Solver) Evaluate() {
    s.Probabilities = NewProbabilities()
    for _, ship := range s.Fleet.ships {
        // horizontal
        for x, row := range s.HuntBoard.xy {
            for y := 0; y < BOARD_SIZE; y++ {
                if row & (ship.Mask>>y) == ship.Mask>>y {
                    for i := 0; i < ship.Length; i++ {
                        if y + i < BOARD_SIZE {
                            s.Probabilities[x][y+i] += 1
                        }
                    }
                }
            }
        }
        for x, row := range s.HuntBoard.yx {
            for y := 0; y < BOARD_SIZE; y++ {
                if row & (ship.Mask>>y) == ship.Mask>>y {
                    for i := 0; i < ship.Length; i++ {
                        if y + i < BOARD_SIZE {
                            s.Probabilities[y+i][x] += 1
                        }
                    }
                }
            }
        }
    }
}

func (s *Solver) EvaluateTarget() {
    s.Probabilities = NewProbabilities()
    //outOfBoundsMask, _ := strconv.ParseUint("0000000000111111", 2, 64)
    //for _, ship := range s.Fleet.ships {
    ship := s.Fleet.ships["submarine"]
        for x, row := range s.TargetBoard.xy {
            if x < 1 {
                for y := 0; y < BOARD_SIZE; y++ {
                    fmt.Printf("X: %d\nY: %d\nrow1: %s\nrow2: %s\nship: %s\ncalc: %s\n\n", x, y, bin(s.HuntBoard.xy[x]), bin(row), bin(ship.Mask>>y), bin(row | (ship.Mask>>y)))
                    if row | (ship.Mask>>y) > row {
                        for i := 0; i < ship.Length; i++ {
                            if y + i < BOARD_SIZE {
                                s.Probabilities[x][y+i] += 1
                            }
                        }
                    }
                }
                // try marching with 1000000000 to see if there's a zero
            }
        }
    //}
}

func bin(n uint) string {
    return strconv.FormatUint(uint64(n), 2)
}
