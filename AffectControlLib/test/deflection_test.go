package tests

import "github.com/the-friyia/go-affect/AffectControlLib"
import "testing"
import "fmt"

// func TestAverage(t *testing.T) {
//
//     a := affect.CalculateDeflection([3]float32{2.0, 2.0, 2.0}, [3]float32{1.0, 1.0, 1.0})
//
//     if a != 1.0 {
//         t.Error(a)
//     }
// }

func Test_calculateActorTransient(t *testing.T) {
    fmt.Println(affect.CalculateActorTransient([3]float32{-2.0, -3.0, 1.0}))
}