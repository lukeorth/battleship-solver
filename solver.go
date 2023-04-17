package battleshipsolver

const (
    huntMode = "hunt"
    targetMode = "target"
)

type Solver struct {
    Probabilities *probabilities
    mode string
    fleet *fleet
    huntBoard *board
    targetBoard *board
}

func NewSolver() *Solver {
    solver := &Solver{
        Probabilities: newProbabilities(),
        fleet: buildFleet(),
        huntBoard: newBoard(),
        targetBoard: newBoard(),
        mode: huntMode,
    }
    return solver
}

func (s *Solver) Hit(location Location) {
    s.mode = targetMode
    s.targetBoard.mark(location)
}

func (s *Solver) Miss(location Location) {
    s.mode = huntMode
    s.huntBoard.mark(location)
}

func (s *Solver) HitAndSunk(location Location, ship string) {
    s.mode = huntMode
    s.targetBoard.mark(location)
    s.huntBoard.merge(s.targetBoard)
    s.fleet.sinkShip(ship)
}

func (s *Solver) Evaluate() {
    s.Probabilities = newProbabilities()
    for _, ship := range s.fleet.ships {
        // horizontal
        for x, row := range s.huntBoard.xy {
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
        for x, row := range s.huntBoard.yx {
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
    outOfBoundsMask := ^rowMask
    for _, ship := range s.fleet.ships {
        for x, row := range s.targetBoard.xy {
            row = row | uint(outOfBoundsMask)
            for y := 0; y < boardSize; y++ {
                if row | (ship.mask>>y) > row && s.huntBoard.xy[x] & (ship.mask>>y) == ship.mask>>y {
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
        for x, row := range s.targetBoard.yx {
            row = row | uint(outOfBoundsMask)
            for y := 0; y < boardSize; y++ {
                if row | (ship.mask>>y) > row && s.huntBoard.yx[x] & (ship.mask>>y) == ship.mask>>y {
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
