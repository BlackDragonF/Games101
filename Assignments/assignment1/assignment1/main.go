package main

import (
  "fmt"
  "assignment1/triangle"
)

func main() {
  triangle := triangle.NewTriangle()
  // initial value
  fmt.Println(triangle.GetA())
  fmt.Println(triangle.GetB())
  fmt.Println(triangle.GetC())
  fmt.Println(triangle.ToVector4())
  // set value
  if ok, err := triangle.SetVertex(0, [3]float64{1, 2, 3}); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetVertex(1, [3]float64{4, 5, 6}); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetVertex(2, [3]float64{7, 8, 9}); !ok {
    fmt.Println(err)
  }
  // changed value
  fmt.Println(triangle.GetA())
  fmt.Println(triangle.GetB())
  fmt.Println(triangle.GetC())
  fmt.Println(triangle.ToVector4())
  // err
  if ok, err := triangle.SetVertex(-3, [3]float64{1, 2, 3}); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetVertex(3, [3]float64{1, 2, 3}); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetColor(-3, 1, 2, 3); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetColor(3, 1, 2, 3); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetColor(1, -1, 2, 3); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetColor(1, 1, 1234, 3); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetTexCoord(-3, 1, 2); !ok {
    fmt.Println(err)
  }
  if ok, err := triangle.SetTexCoord(3, 1, 2); !ok {
    fmt.Println(err)
  }
}
