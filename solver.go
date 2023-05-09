package battleshipsolver

import (
	"encoding/json"
	"fmt"
)

type Solver struct {
    Probabilities *probabilities
    BestCell Location 
    fleet *fleet
    huntBoard *board
    targetBoard *board
    sunkBoard *board
}

type SolverJSON struct {
    Fleet []string `json:"fleet"`
    Board [boardSize][boardSize]int `json:"board"`
}

func NewSolver() *Solver {
    solver := &Solver{
        Probabilities: newProbabilities(),
        fleet: buildFleet(),
        huntBoard: newBoard(),
        targetBoard: newBoard(),
        sunkBoard: newBoard(),
    }
    return solver
}

func (s *Solver) Hit(location Locator) {
    // remove
    // fmt.Println(s.isSinkableRow(location.Locate().Row, location.Locate().Col, s.fleet.ships[Battleship]))
    fmt.Println(s.isSinkableCol(location.Locate().Row, location.Locate().Col, s.fleet.ships[Battleship]))
    s.fleet.hit()
    s.targetBoard.mark(location)
}

func (s *Solver) Miss(location Locator) {
    s.huntBoard.mark(location)
}

func (s *Solver) HitAndSunk(location Locator, ship string) {
    s.targetBoard.mark(location)
    s.fleet.sink(ship)
    if !s.isTargetMode() {
        s.huntBoard.merge(s.targetBoard)
    }
}

func (s *Solver) Evaluate() {
    s.Probabilities = newProbabilities()
    for _, ship := range s.fleet.floatingShips() {
        for row := 0; row < boardSize; row++ {
            for col := 0; col < boardSize; col++ {
                s.evaluateRow(row, col, ship)
                s.evaluateCol(row, col, ship)
            }
        }
    }
    s.updateBestCell()
    // remove these
    fmt.Println(s.huntBoard.String())
    fmt.Println(s.targetBoard.String())
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

func (s *Solver) isSinkableRow(row int, col int, ship *ship) bool {
    hitRow := (^s.targetBoard[row] & rowMask) | (pegMask>>col)
    colShift := col - ship.length + 1
    for i := 0; i < ship.length; i++ {
        if s.isPlayableRow(row, colShift + i, ship) {
            if hitRow == ship.mask>>(colShift + i) | hitRow {
                return true
            }
        }
    }
    return false
}

func (s *Solver) isSinkableCol(row int, col int, ship *ship) bool {
    rowShift := row - ship.length + 1
    if rowShift < 0 {
        rowShift = 0
    }
    for i := 0; i < ship.length; i++ {
        if s.isPlayableCol(rowShift + i, col, ship) {
            var rowCopy uint
            for j := 0; j < ship.length; j++ {
                fmt.Println(rowShift + i + j)
                if rowShift + i + j == row {
                    continue
                }
                fmt.Printf("%010b\n", s.targetBoard[rowShift + i + j])
                rowCopy |= s.targetBoard[rowShift + i + j]
            }
            fmt.Printf("%010b\n", rowCopy & (pegMask>>col))
            if rowCopy & (pegMask>>col) != (pegMask>>col) {
                return true
            }
        }
    }
    return false
}

func (s *Solver) isTargetableRow(row int, col int, ship *ship) bool {
    if s.isPlayableRow(row, col, ship) {
        return isTargetable(s.targetBoard[row], ship.mask>>col) 
    }
    return false
}

func (s *Solver) isTargetableCol(row int, col int, ship *ship) bool {
    if s.isPlayableCol(row, col, ship) {
        rowCopy := s.targetBoard[row]
        for i := 0; i < ship.length; i++ {
            rowCopy &= s.targetBoard[row + i]
        }
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
        rowCopy := s.huntBoard[row]
        for i := 0; i < ship.length; i++ {
            rowCopy &= s.huntBoard[row+i]
        }
        return isPlayable(rowCopy, pegMask>>col)
    }
    return false
}

func (s *Solver) isTargetMode() bool {
    return s.fleet.hitCount > 0
}

func (s *Solver) updateBestCell() {
    for row := 0; row < boardSize; row++ {
        for col := 0; col < boardSize; col++ {
            bestCell := s.Probabilities[s.BestCell.Row][s.BestCell.Col]
            if s.Probabilities[row][col] > bestCell {
                s.BestCell = Cell(row, col).Locate()
            }
        }
    }
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
    return position >= 0 && position <= boardSize - shipLen
}

func (s Solver) MarshalJSON() ([]byte, error) {
    bestCell := struct{
        Coordinates []int `json:"coordinates"`
        Position Position `json:"position"`
    }{
        []int{s.BestCell.Row, s.BestCell.Col},
        s.BestCell.Position,
    }

    return json.Marshal(map[string]interface{}{
        "probabilities": s.Probabilities,
        "bestCell": bestCell,
    })
}

func (s *Solver) UnmarshalJSON(data []byte) error {
    tempSolver := &SolverJSON{}

    if err := json.Unmarshal(data, &tempSolver); err != nil {
        return err
    }

    s.huntBoard = &board{}
    s.targetBoard = &board{}
    s.fleet = &fleet{
        ships: make(map[string]*ship),
    }

    for row := 0; row < boardSize; row++ {
        huntRow := rowMask
        targetRow := rowMask
        for col := 0; col < boardSize; col++ {
            switch tempSolver.Board[row][col] {
            case 0, 3:
                huntRow = huntRow ^ (pegMask>>col)
            case 2:
                s.fleet.hitCount++
                targetRow = targetRow ^ (pegMask>>col)
            }
        }
        s.huntBoard[row] = huntRow
        s.targetBoard[row] = targetRow
    }

    for _, tempShip := range tempSolver.Fleet {
        switch tempShip {
        case Carrier:
            s.fleet.ships[Carrier] = &ship{Carrier, carrierMask, carrierLength, false}
        case Battleship:
            s.fleet.ships[Battleship] = &ship{Battleship, battleshipMask, battleshipLength, false}
        case Submarine:
            s.fleet.ships[Submarine] = &ship{Submarine, submarineMask, submarineLength, false}
        case Cruiser:
            s.fleet.ships[Cruiser] = &ship{Cruiser, cruiserMask, cruiserLength, false}
        case Destroyer:
            s.fleet.ships[Destroyer] = &ship{Destroyer, destroyerMask, destroyerLength, false}
        }
    }

    return nil
}
