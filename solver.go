package battleshipsolver

import (
	"encoding/json"
)

type Solver struct {
    Probabilities *probabilities
    fleet *fleet
    huntBoard *board
    targetBoard *board
}

type SolverJSON struct {
    Board [boardSize][boardSize]int `json:"board"`
    Fleet []string `json:"fleet"`
}

func (s *Solver) UnmarshalJSON(data []byte) error {
    solver := &SolverJSON{}

    if err := json.Unmarshal(data, &solver); err != nil {
        return err
    }

    tempHuntBoard := &board{}
    tempTargetBoard := &board{}
    tempHitCount := 0

    for row := 0; row < boardSize; row++ {
        huntRow := rowMask
        targetRow := rowMask
        for col := 0; col < boardSize; col++ {
            switch solver.Board[row][col] {
            case 0, 3:
                huntRow = huntRow ^ (pegMask>>col)
            case 2:
                tempHitCount++
                targetRow = targetRow ^ (pegMask>>col)
            }
        }
        tempHuntBoard[row] = huntRow
        tempTargetBoard[row] = targetRow
    }
    s.huntBoard = tempHuntBoard
    s.targetBoard = tempTargetBoard

    tempShips := make(map[string]*ship)
    for _, tempShip := range solver.Fleet {
        switch tempShip {
        case Carrier:
            tempShips[Carrier] = &ship{Carrier, carrierMask, carrierLength}
        case Battleship:
            tempShips[Battleship] = &ship{Battleship, battleshipMask, battleshipLength}
        case Submarine:
            tempShips[Submarine] = &ship{Submarine, submarineMask, submarineLength}
        case Cruiser:
            tempShips[Cruiser] = &ship{Cruiser, cruiserMask, cruiserLength}
        case Destroyer:
            tempShips[Destroyer] = &ship{Destroyer, destroyerMask, destroyerLength}
        }
    }

    s.fleet = &fleet{
        ships: tempShips,
        hitCount: tempHitCount,
    }

    return nil
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
    s.fleet.sunk(ship)
    if !s.isTargetMode() {
        s.huntBoard.merge(s.targetBoard)
    }
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
    switch s.isTargetMode() {
    case true:
        if s.isTargetableRow(row, col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[row][col+i] += 1
            }
        }
        if isAlreadyMarked(s.targetBoard[row], pegMask>>col) {
            s.Probabilities[row][col] = 0
        }
    case false:
        if s.isPlayableRow(row, col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[row][col+i] += 1
            }
        }
    }
}

func (s *Solver) evaluateCol(row int, col int, ship *ship) {
    switch s.isTargetMode() {
    case true:
        if s.isTargetableCol(row, col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[row+i][col] += 1
            }
        }
        if isAlreadyMarked(s.targetBoard[row], pegMask>>col) {
            s.Probabilities[row][col] = 0
        }
    case false:
        if s.isPlayableCol(row, col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[row+i][col] += 1
            }
        }
    }
}

func (s *Solver) isTargetableRow(row int, col int, ship *ship) bool {
    if s.isPlayableRow(row, col, ship) {
        return isTargetable(s.targetBoard[row], ship.mask>>col) 
    }
    return false
}

func (s *Solver) isTargetableCol(row int, col int, ship *ship) bool {
    if s.isPlayableCol(row, col, ship) {
        rowCopy := s.targetBoard.condenseRows(row, ship.length)

        return isTargetable(rowCopy, pegMask>>col)
    }
    return false
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

func (s *Solver) isTargetMode() bool {
    return s.fleet.hitCount > 0
}

func isAlreadyMarked(rowMask uint, evalMask uint) bool {
    return rowMask & evalMask == 0
}

func isTargetable(rowMask uint, evalMask uint) bool {
    return rowMask | evalMask > rowMask
}

func isPlayable(rowMask uint, evalMask uint) bool {
    return rowMask & evalMask == evalMask
}

func isInBounds(position int, shipLen int) bool {
    return position <= boardSize - shipLen
}
