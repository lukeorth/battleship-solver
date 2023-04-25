package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	bs "github.com/lukeorth/battleship-solver"
)

func main() {
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


    var solver bs.Solver
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
    //run(solver)

    //solver.TestBoardUnmarshal(boardBlob)

    /*
    solver.Miss("F7")
    solver.Hit("F6")
    solver.Miss("G6")
    //solver.Miss("A4")
    solver.Hit("D6")
    solver.HitAndSunk("E6", battleshipsolver.Destroyer)
    solver.Hit("A6")
    //solver.EvaluateTarget()
    solver.Miss("A7")
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
