package main

import (
	"bufio"
	"fmt"
	"os"

	bs "github.com/lukeorth/battleship-solver"
)

func main() {
    /*
    var boardBlob = []byte(`{"board": [
    [2, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 2, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1]
    ],
    "fleet": ["carrier","battleship","cruiser","submarine","destroyer"]}`)
    */

    //var solver bs.Solver
    solver := bs.NewSolver();
    /*
    if err := json.Unmarshal(boardBlob, &solver); err != nil {
        fmt.Printf("ERROR: %s", err)
    }
    solver.Evaluate()
    fmt.Println(solver.Probabilities.String())

    b, err := json.Marshal(solver)
    if err != nil {
        fmt.Printf("ERROR: %s", err)
    }
    os.Stdout.Write(b)
    fmt.Println()
    */
    run(solver)

    //solver.TestBoardUnmarshal(boardBlob)
    /*
    var solver bs.Solver
    solver.Miss(bs.Cell(5, 7))
    solver.Hit(bs.Cell(5, 6))
    solver.Miss(bs.Cell(6, 6))
    //solver.Miss("A4")
    solver.Hit(bs.Cell(3, 6))
    solver.HitAndSunk(bs.Cell(4, 6), bs.Destroyer)
    solver.Hit(bs.Cell(0, 6))
    //solver.EvaluateTarget()
    solver.Miss(bs.Cell(0, 7))
    solver.Evaluate()

    //fmt.Println(solver.HuntBoard.String())
    fmt.Println(solver.Probabilities.String())
    */
}

func run(s *bs.Solver) {
    in := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println(s.Probabilities.String())
        move := getMoveType(in)
        switch move {
        case "HIT":
            location := getLocation(in)
            s.Hit(bs.Position(location))
            s.Evaluate()
        case "MISS":
            location := getLocation(in)
            s.Miss(bs.Position(location))
            s.Evaluate()
        case "SUNK":
            location := getLocation(in)
            ship := getShip(in)
            s.HitAndSunk(bs.Position(location), ship)
            s.Evaluate()
        default:
            break
        }
    }
}

func getMoveType(s *bufio.Scanner) string {
    fmt.Printf("Move: ")
    s.Scan()
    move := s.Text()
    return move
}

func getLocation(s *bufio.Scanner) string {
    fmt.Printf("Location: ")
    s.Scan()
    return s.Text()
}

func getShip(s *bufio.Scanner) string {
    fmt.Printf("Ship: ")
    s.Scan()
    ship := s.Text()
    return ship
}
