package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lukeorth/battleship-solver"
)

func main() {
    var boardBlob = []byte(`{"board": [
    [1, 1, 1, 1, 1, 1, 1, 1, 0, 1],
    [1, 1, 1, 1, 3, 1, 1, 1, 0, 1],
    [1, 1, 1, 1, 2, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 2, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 2, 0, 1, 1, 1, 1],
    [1, 1, 1, 1, 0, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 0, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1],
    [1, 1, 1, 1, 1, 1, 1, 1, 1, 1]
    ],
    "fleet": ["carrier","battleship","cruiser","submarine","destroyer"]}`)


    var solver battleshipsolver.Solver
    if err := json.Unmarshal(boardBlob, &solver); err != nil {
        fmt.Printf("ERROR: %s", err)
    }
    solver.Evaluate()
    fmt.Println(solver.Probabilities.String())
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

func run(s *battleshipsolver.Solver) {
    in := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println(s.Probabilities.String())
        move := getMoveType(in)
        switch move {
        case "HIT":
            location := getLocation(in)
            s.Hit(location)
            s.Evaluate()
        case "MISS":
            location := getLocation(in)
            s.Miss(location)
            s.Evaluate()
        case "SUNK":
            location := getLocation(in)
            ship := getShip(in)
            s.HitAndSunk(location, ship)
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

func getLocation(s *bufio.Scanner) battleshipsolver.Location {
    fmt.Printf("Location: ")
    s.Scan()
    location := s.Text()
    return battleshipsolver.Location(location)
}

func getShip(s *bufio.Scanner) string {
    fmt.Printf("Ship: ")
    s.Scan()
    ship := s.Text()
    return ship
}
