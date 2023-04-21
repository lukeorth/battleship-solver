package battleshipsolver

type Solver struct {
    Probabilities *probabilities
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
    }
    return solver
}

func (s *Solver) Hit(location Location) {
    s.fleet.hit()
    s.targetBoard.mark(location)
}

func (s *Solver) Miss(location Location) {
    s.huntBoard.mark(location)
}

func (s *Solver) HitAndSunk(location Location, ship string) {
    s.targetBoard.mark(location)
    s.huntBoard.merge(s.targetBoard)
    s.fleet.sunk(ship)
}

func (s *Solver) Evaluate() {
    s.Probabilities = newProbabilities()
    for _, ship := range s.fleet.ships {
        for row := 0; row < boardSize; row++ {
            for col := 0; col < boardSize; col++ {
                s.evaluateRow(row, col, ship)
                s.evaluateCol(row, col, ship)
            }
        }
    }
}

func (s *Solver) evaluateRow(row int, col int, ship *ship) {
    if s.isPlayableRow(row, col, ship) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row][col+i] += 1
        }
    }
}

func (s *Solver) evaluateCol(row int, col int, ship *ship) {
    if s.isPlayableCol(row, col, ship) {
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
    targetRow := s.targetBoard.condenseRows(row, ship.length)
    huntRow := s.huntBoard.condenseRows(row, ship.length)

    if isHitIntersect(targetRow, pegMask>>col) && isPlayable(huntRow, pegMask>>col) {
        for i := 0; i < ship.length; i++ {
            s.Probabilities[row+i][col] += 1
        }
    }
}

func (s *Solver) isPlayableRow(row int, col int, ship *ship) bool {
    if isInBounds(col, ship.length) {
        return isPlayable(s.huntBoard[row], ship.mask>>col)
    }
    return false
}

func (s *Solver) isPlayableCol(row int, col int, ship *ship) bool {
    if isInBounds(row, ship.length) {
        rowCopy := s.huntBoard.condenseRows(row, ship.length)
        
        return isPlayable(rowCopy, pegMask>>col)
    }
    return false
}

func isHitIntersect(rowMask uint, evalMask uint) bool {
    return rowMask | evalMask > rowMask
}

func isPlayable(rowMask uint, evalMask uint) bool {
    return rowMask & evalMask == evalMask
}

func isInBounds(position int, shipLen int) bool {
    return position <= boardSize - shipLen
}
