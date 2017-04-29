package tests

import "github.com/the-friyia/go-affect/AffectControlLib"
import "testing"
import "fmt"

func Test_CalculateTransient(t *testing.T) {
	newFundamental := affect.CalculateTransient([3]float64{1, 1, 1},
		[3]float64{1, 1, 1},
		[3]float64{1, 1, 1})
	fmt.Println(newFundamental)
}

func Test_DeflectionCalculation(t *testing.T) {
	d := affect.CalculateDeflection([3]float64{2, 1, 1},
		[3]float64{1, 3, 1},
		[3]float64{1, 1, 4},
		[3]float64{4, 1, 1},
		[3]float64{3, 1, 1},
		[3]float64{1, 1, 1})
	fmt.Println(d)
}

func Test_make(t *testing.T) {
	affect.MakeAffectiveState()
}

func Test_LoadFunction(t *testing.T) {
	affect.LoadBehaviours()
}

func Test_simulatedInteraction(t *testing.T) {
    state := affect.MakeAffectiveState()
    fmt.Println("#### Simulation start")

    // Actor says something to the client
    state.PropegateForward(state.Behaviours["implore"])
    fmt.Println(state.Deflection)

    // Object Speaks
    state.PropegateForward(state.Behaviours["lecture"])
    fmt.Println(state.Deflection)

    // Actor says something to the client
    state.PropegateForward(state.Behaviours["discourage"])
    fmt.Println(state.Deflection)

    // Object Speaks
    state.PropegateForward(state.Behaviours["compliment"])
    fmt.Println(state.Deflection)

    // Actor says something to the client
    state.PropegateForward(state.Behaviours["counsel"])
    fmt.Println(state.Deflection)
}
