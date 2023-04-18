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
        for row := range s.huntBoard {
            for col := 0; col < boardSize; col++ {
                if isInBounds(col, ship.length) {
                    s.evaluateRow(row, col, ship)
                }
                if isInBounds(row, ship.length) {
                    s.evaluateCol(row, col, ship)
                }
            }
        }
    }
}

func (s *Solver) evaluateRow(row int, col int, ship *ship) {
    if isPlayable(s.huntBoard[row], ship.mask>>col) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row][col+i] += 1
        }
    }
}

func (s *Solver) evaluateCol(row int, col int, ship *ship) {
    rowCopy := s.huntBoard[row]
    for i := 0; i < ship.length; i++ {
        rowCopy &= s.huntBoard[row+i]
    }
    if isPlayable(rowCopy, pegMask>>col) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row+i][col] += 1
        }
    }
}

func (s *Solver) EvaluateTarget() {
    s.Probabilities = newProbabilities()
    for _, ship := range s.fleet.ships {
        for row := range s.targetBoard {
            for col := 0; col < boardSize; col++ {
                if isInBounds(col, ship.length) {
                    s.evaluateTargetRow(row, col, ship)
                }
                if isInBounds(row, ship.length) {
                    s.evaluateTargetCol(row, col, ship)
                }
                if s.targetBoard[row] & (uint(pegMask)>>col) == 0 {
                    s.Probabilities[row][col] = 0
                }
            }
        }
    }
}

func (s *Solver) evaluateTargetRow(row int, col int, ship *ship) {
    if isHitIntersect(s.targetBoard[row], ship.mask>>col) && isPlayable(s.huntBoard[row], ship.mask>>col) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row][col+i] += 1
        }
    }
}

func (s *Solver) evaluateTargetCol(row int, col int, ship *ship) {
    rowCopy := s.targetBoard[row]
    for i := 0; i < ship.length; i++ {
        rowCopy &= s.targetBoard[row+i]
    }
    if isHitIntersect(rowCopy, pegMask>>col) && isPlayable(s.huntBoard[row], pegMask>>col) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row+i][col] += 1
        }
    }
}

func isHitIntersect(rowMask uint, shipMask uint) bool {
    return rowMask | shipMask > rowMask
}

func isPlayable(rowMask uint, shipMask uint) bool {
    return rowMask & shipMask == shipMask
}

func isInBounds(position int, shipLen int) bool {
    return position <= boardSize - shipLen
}
