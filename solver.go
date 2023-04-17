package battleshipsolver

import (
	"strconv"
)

const (
    huntMode = "hunt"
    targetMode = "target"
)

type Solver struct {
    Probabilities *probabilities
    BestMove Location
    mode string
    Fleet *fleet
    HuntBoard *board
    TargetBoard *board
}

func NewSolver() *Solver {
    solver := &Solver{
        Probabilities: newProbabilities(),
        Fleet: buildFleet(),
        HuntBoard: newBoard(),
        TargetBoard: newBoard(),
        mode: huntMode,
    }
    return solver
}

func (s *Solver) Hit(location Location) {
    s.mode = targetMode
    s.TargetBoard.mark(location)
}

func (s *Solver) Miss(location Location) {
    s.mode = huntMode
    s.HuntBoard.mark(location)
}

func (s *Solver) HitAndSunk(location Location, ship string) {
    s.mode = huntMode
    s.TargetBoard.mark(location)
    s.HuntBoard.merge(s.TargetBoard)
    s.Fleet.sinkShip(ship)
}

func (s *Solver) Evaluate() {
    s.Probabilities = newProbabilities()
    for _, ship := range s.Fleet.ships {
        // horizontal
        for x, row := range s.HuntBoard.xy {
            for y := 0; y < boardSize; y++ {
                if row & (ship.mask>>y) == ship.mask>>y {
                    for i := 0; i < ship.length; i++ {
                        if y + i < boardSize {
                            s.Probabilities[x][y+i] += 1
                        }
                    }
                }
            }
        }
        for x, row := range s.HuntBoard.yx {
            for y := 0; y < boardSize; y++ {
                if row & (ship.mask>>y) == ship.mask>>y {
                    for i := 0; i < ship.length; i++ {
                        if y + i < boardSize {
                            s.Probabilities[y+i][x] += 1
                        }
                    }
                }
            }
        }
    }
}

func (s *Solver) EvaluateTarget() {
    s.Probabilities = newProbabilities()
    outOfBoundsMask, _ := strconv.ParseUint("0000000000111111", 2, 64)
    for _, ship := range s.Fleet.ships {
        for x, row := range s.TargetBoard.xy {
            row = row | uint(outOfBoundsMask)
            for y := 0; y < boardSize; y++ {
                if row | (ship.mask>>y) > row && s.HuntBoard.xy[x] & (ship.mask>>y) == ship.mask>>y {
                    for i := 0; i < ship.length; i++ {
                        if y + i < boardSize {
                            s.Probabilities[x][y+i] += 1
                        }
                    }
                }
                if row & (uint(pegMask)>>y) == 0 {
                    s.Probabilities[x][y] = 0
                }
            }
        }
        for x, row := range s.TargetBoard.yx {
            row = row | uint(outOfBoundsMask)
            for y := 0; y < boardSize; y++ {
                if row | (ship.mask>>y) > row && s.HuntBoard.yx[x] & (ship.mask>>y) == ship.mask>>y {
                    for i := 0; i < ship.length; i++ {
                        if y + i < boardSize {
                            s.Probabilities[y+i][x] += 1
                        }
                    }
                }
                if row & (uint(pegMask)>>y) == 0 {
                    s.Probabilities[y][x] = 0
                }
            }
        }
    }
}

func bin(n uint) string {
    return strconv.FormatUint(uint64(n), 2)
}
