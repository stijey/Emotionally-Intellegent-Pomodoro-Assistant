package tests

import "github.com/the-friyia/go-affect/AffectControlLib"
import "testing"
import "fmt"

// func TestAverage(t *testing.T) {
//
//     a := affect.CalculateDeflection([3]float64{2.0, 2.0, 2.0}, [3]float64{1.0, 1.0, 1.0})
//
//     if a != 1.0 {
//         t.Error(a)
//     }
// }

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
