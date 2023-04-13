package main

import (
	"fmt"

	"github.com/lukeorth/battleship-solver"
)

func main() {
    solver := battleshipsolver.NewSolver()
    solver.Hit("A2")
    solver.Hit("A3")
    solver.Miss("A4")
    //solver.HitAndSunk("C4", battleshipsolver.CARRIER_NAME)
    solver.EvaluateTarget()

    //fmt.Println(solver.HuntBoard.String())
    fmt.Println(solver.TargetBoard.String())
    fmt.Println(solver.Probabilities.String())
}
