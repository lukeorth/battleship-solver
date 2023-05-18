package battleshipsolver

import (
	"encoding/json"
	"fmt"
)

type Solver struct {
    Probabilities *probabilities
    BestCell Location 
    Errors []error
    fleet *fleet
    huntBoard *board
    targetBoard *board
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
    }
    return solver
}

func (s *Solver) Hit(location Locator) {
    s.fleet.hit()
    s.targetBoard.mark(location)
}

func (s *Solver) Miss(location Locator) {
    s.huntBoard.mark(location)
}

func (s *Solver) HitAndSunk(location Locator, shipName string) {
    s.fleet.hit()
    ship := s.fleet.ships[shipName]
    s.fleet.hitCount -= ship.length
    s.fleet.sink(location, shipName)
}

func (s *Solver) Evaluate() {
    s.Probabilities = newProbabilities()
    s.sinkShips(s.fleet.sunkShips())
    for _, ship := range s.fleet.floatingShips() {
        for row := 0; row < boardSize; row++ {
            for col := 0; col < boardSize; col++ {
                s.evaluateRow(Cell(row, col).Locate(),  ship)
                s.evaluateCol(Cell(row, col).Locate(), ship)
            }
        }
    }
    s.updateBestCell()
}

func (s *Solver) evaluateRow(loc Location, ship *ship) {
    switch s.isTargetMode() {
    case true:
        if s.isTargetableRow(loc, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[loc.Row][loc.Col+i] += 1
            }
        }
        if isAlreadyMarked(s.targetBoard[loc.Row], pegMask>>loc.Col) {
            s.Probabilities[loc.Row][loc.Col] = 0
        }
    case false:
        if s.isPlayableRow(loc.Row, loc.Col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[loc.Row][loc.Col+i] += 1
            }
        }
    }
}

func (s *Solver) evaluateCol(loc Location, ship *ship) {
    switch s.isTargetMode() {
    case true:
        if s.isTargetableCol(loc, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[loc.Row+i][loc.Col] += 1
            }
        }
        if isAlreadyMarked(s.targetBoard[loc.Row], pegMask>>loc.Col) {
            s.Probabilities[loc.Row][loc.Col] = 0
        }
    case false:
        if s.isPlayableCol(loc.Row, loc.Col, ship) {
            for i := 0; i < ship.length; i++ {
                s.Probabilities[loc.Row+i][loc.Col] += 1
            }
        }
    }
}

func (s *Solver) sinkShips(ships []*ship) {
    if len(ships) == 0 {
        return
    }
    ship := ships[0]
    row := ship.sunkAt.Locate().Row
    col := ship.sunkAt.Locate().Col

    rowPositions := make([][]int, 0, ship.length)
    colPositions := make([][]int, 0, ship.length)
    rowShift := row - ship.length + 1
    colShift := col - ship.length + 1
    hitRow := (^s.targetBoard[row]) | (pegMask>>col)

    for i := 0; i < ship.length; i++ {
        if s.isPlayableRow(row, colShift + i, ship) {
            if hitRow == ship.mask>>(colShift + i) | hitRow {
                rowPositions = append(rowPositions, []int{row, colShift + i})
            }
        }
        if s.isPlayableCol(rowShift + i, col, ship) {
            var rowCopy uint
            for j := 0; j < ship.length; j++ {
                if rowShift + i + j != row {
                    rowCopy |= s.targetBoard[rowShift + i + j]
                }
            }
            if rowCopy & (pegMask>>col) != pegMask>>col {
                colPositions = append(colPositions, []int{rowShift + i, col})
            }
        }
    }
    positions := append(rowPositions, colPositions...)
    if len(positions) == 0 {
        s.huntBoard.mark(Cell(row, col))
        s.fleet.remove(ship.name)
        s.Errors = append(s.Errors, fmt.Errorf("Illegal sink: %s at %s", ship.name, Cell(row, col).Locate().Position))
        s.sinkShips(ships[1:])
    } else if len(positions) > 1 {
        s.targetBoard.mark(Cell(row, col))
        s.sinkShips(ships[1:])
    } else if len(rowPositions) == 1 {
        s.fleet.hitCount -= ship.length
        s.fleet.remove(ship.name)
        for i := 0; i < ship.length; i++ {
            s.huntBoard.mark(Cell(row, rowPositions[0][1] + i))
        }
        s.sinkShips(s.fleet.sunkShips())
    } else {
        s.fleet.hitCount -= ship.length
        s.fleet.remove(ship.name)
        for i := 0; i < ship.length; i++ {
            s.huntBoard.mark(Cell(colPositions[0][0] + i, col))
        }
        s.sinkShips(s.fleet.sunkShips())
    }
}

func (s *Solver) isTargetableRow(loc Location, ship *ship) bool {
    if s.isPlayableRow(loc.Row, loc.Col, ship) {
        return isTargetable(s.targetBoard[loc.Row], ship.mask>>loc.Col) 
    }
    return false
}

func (s *Solver) isTargetableCol(loc Location, ship *ship) bool {
    if s.isPlayableCol(loc.Row, loc.Col, ship) {
        rowCopy := s.targetBoard[loc.Row]
        for i := 0; i < ship.length; i++ {
            rowCopy &= s.targetBoard[loc.Row + i]
        }
        return isTargetable(rowCopy, pegMask>>loc.Col)
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
    var errorMsgs []string
    for _, err := range s.Errors {
        errorMsgs = append(errorMsgs, err.Error())
    }

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
        "errors": errorMsgs,
    })
}

func (s *Solver) UnmarshalJSON(data []byte) error {
    tempSolver := &SolverJSON{}

    if err := json.Unmarshal(data, &tempSolver); err != nil {
        return err
    }

    s.huntBoard = &board{}
    s.targetBoard = &board{}
    s.fleet = buildFleet()

    for row := 0; row < boardSize; row++ {
        huntRow := rowMask
        targetRow := rowMask
        for col := 0; col < boardSize; col++ {
            switch tempSolver.Board[row][col] {
            case 0:
                huntRow = huntRow ^ (pegMask>>col)
            case 1:
                s.fleet.hitCount++
                s.fleet.sink(Cell(row, col), Carrier)
            case 2:
                s.fleet.hitCount++
                s.fleet.sink(Cell(row, col), Battleship)
            case 3:
                s.fleet.hitCount++
                s.fleet.sink(Cell(row, col), Submarine)
            case 4:
                s.fleet.hitCount++
                s.fleet.sink(Cell(row, col), Cruiser)
            case 5:
                s.fleet.hitCount++
                s.fleet.sink(Cell(row, col), Destroyer)
            case 6:
                s.fleet.hitCount++
                targetRow = targetRow ^ (pegMask>>col)
            }
        }
        s.huntBoard[row] = huntRow
        s.targetBoard[row] = targetRow
    }

    for _, ship := range s.fleet.ships {
        sunk := true
        for _, tempShip := range tempSolver.Fleet {
            if ship.name == tempShip {
                sunk = false
            }
        }
        if ship.sunkAt != nil {
            s.fleet.hitCount -= 1
            sunk = false
        }
        if sunk {
            s.fleet.remove(ship.name)
        }
    }

    fmt.Println(s.fleet.hitCount)

    return nil
}
