package main

import (
	"fmt"

	"github.com/lukeorth/battleship-solver"
)

func main() {
    solver := battleshipsolver.NewSolver()
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
    fmt.Println(solver.HuntBoard.String())
    fmt.Println(solver.TargetBoard.String())
    fmt.Println(solver.Probabilities.String())
}
