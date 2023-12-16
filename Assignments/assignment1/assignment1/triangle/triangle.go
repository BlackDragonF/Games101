package triangle

import (
  "errors"
  "fmt"
  "gonum.org/v1/gonum/mat"
)

type Triangle struct {
  /* the original coordinates of the triangle, v0, v1, v2 in
  counter clockwise order */
  v          [3]mat.VecDense
  /* Per vertex values */
  color      [3]mat.VecDense  // color at each vertex;
  tex_coords [3]mat.VecDense  // texture u,v
  normal     [3]mat.VecDense  // normal vector for each vertex
}

func NewTriangle() Triangle {
  return Triangle {
    v: [3]mat.VecDense{
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.})}, 

    color: [3]mat.VecDense{
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.})},

    tex_coords: [3]mat.VecDense{
      *mat.NewVecDense(2, []float64{0., 0.}),
      *mat.NewVecDense(2, []float64{0., 0.}),
      *mat.NewVecDense(2, []float64{0., 0.})}, 

    normal: [3]mat.VecDense{
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.}),
      *mat.NewVecDense(3, []float64{0., 0., 0.})}, 
  }
}

func (triangle Triangle) GetA() mat.VecDense {
  return triangle.v[0]
}

func (triangle Triangle) GetB() mat.VecDense {
  return triangle.v[1]
}

func (triangle Triangle) GetC() mat.VecDense {
  return triangle.v[2]
}

/* set i-th vertex coordinates */
func (triangle Triangle) SetVertex(ind int, ver [3]float64) (bool, error) {
  if ind < 0 || ind > 2 {
    return false, errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
  }

  triangle.v[ind].SetVec(0, ver[0])
  triangle.v[ind].SetVec(1, ver[1])
  triangle.v[ind].SetVec(2, ver[2])

  return true, nil
}

/* set i-th vertex normal vector */
func (triangle Triangle) SetNormal(ind int, ver [3]float64) (bool, error) {
  if ind < 0 || ind > 2 {
    return false, errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
  }

  triangle.normal[ind].SetVec(0, ver[0])
  triangle.normal[ind].SetVec(1, ver[1])
  triangle.normal[ind].SetVec(2, ver[2])

  return true, nil
}

/* set i-th vertex normal vector */
func (triangle Triangle) SetColor(ind int, r float64, g float64, b float64) (bool, error) {
  if ind < 0 || ind > 2 {
    return false, errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
  }
  if r < 0. || r > 255. || g < 0. || g > 255. || b < 0. || b > 255 {
    return false, errors.New(fmt.Sprintf("triangle: invalid color value (%.2f, %.2f, %.2f)\n", r, g, b))
  }

  triangle.color[ind].SetVec(0, r/ 255.)
  triangle.color[ind].SetVec(1, g/ 255.)
  triangle.color[ind].SetVec(2, b/ 255.)

  return true, nil
}

/* set i-th vertex texture coordinate */
func (triangle Triangle) SetTexCoord(ind int, s float64, t float64) (bool, error) {
  if ind < 0 || ind > 2 {
    return false, errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
  }

  triangle.tex_coords[ind].SetVec(0, s)
  triangle.tex_coords[ind].SetVec(1, t)

  return true, nil
}

func (triangle Triangle) ToVector4() [3]mat.VecDense {
  res := [3]mat.VecDense {
    *mat.NewVecDense(4, []float64{triangle.v[0].AtVec(0), triangle.v[0].AtVec(1), triangle.v[0].AtVec(2), 1.}),
    *mat.NewVecDense(4, []float64{triangle.v[1].AtVec(0), triangle.v[1].AtVec(1), triangle.v[1].AtVec(2), 1.}),
    *mat.NewVecDense(4, []float64{triangle.v[2].AtVec(0), triangle.v[2].AtVec(1), triangle.v[2].AtVec(2), 1.}),
  }
  return res
}

