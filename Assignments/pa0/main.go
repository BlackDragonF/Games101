package main

import (
    "fmt"
    "gonum.org/v1/gonum/mat"
    "math"
)

func main() {
    p := mat.NewVecDense(3, []float64{2, 1, 1})
    fmt.Printf("point P in homogeneous coordinates: \n%v\n", mat.Formatted(p))
    a := math.Pi/4
    R := mat.NewDense(3, 3, []float64{math.Cos(a), -math.Sin(a), 0, math.Sin(a), math.Cos(a), 0, 0, 0, 1})
    fmt.Printf("Rotation Matrix: \n%v\n", mat.Formatted(R))
    T := mat.NewDense(3, 3, []float64{1, 0, 1, 0, 1, 2, 0, 0, 1})
    fmt.Printf("Translation Matrix: \n%v\n", mat.Formatted(T))
    p.MulVec(R, p)
    p.MulVec(T, p)
    fmt.Printf("point P' in homogeneous coordinates: \n%v\n", mat.Formatted(p))
    pCartesian := mat.NewVecDense(2, []float64{p.At(0, 0), p.At(1, 0)})
    fmt.Printf("point P' in cartesian coordinates: \n%v\n", mat.Formatted(pCartesian))
}
