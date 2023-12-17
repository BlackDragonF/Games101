package common

import (
	"errors"
    "fmt"  

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/mat"
)

type Vec interface {
	Equal() bool
    ToHomoVec(w float64) mat.VecDense
}

type Vec2f [2]float64
type Vec3f [3]float64
type Vec4f [4]float64

type Vec2i [2]int64
type Vec3i [3]int64
type Vec4i [4]int64

func NewVec2f() Vec2f {
	return Vec2f{0., 0.}
}

func NewVec3f() Vec3f {
	return Vec3f{0., 0., 0.}
}

func NewVec4f() Vec4f {
	return Vec4f{0., 0., 0., 0.}
}

func NewVec2i() Vec2i {
	return Vec2i{0, 0}
}

func NewVec3i() Vec3i {
	return Vec3i{0, 0, 0}
}

func NewVec4i() Vec4i {
	return Vec4i{0, 0, 0, 0}
}

func (v1 Vec2f) Equal(v2 Vec2f) bool {
	return (scalar.EqualWithinAbsOrRel(v1[0], v2[0], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[1], v2[1], 1e-6, 1e-6))
}

func (v1 Vec3f) Equal(v2 Vec3f) bool {
	return (scalar.EqualWithinAbsOrRel(v1[0], v2[0], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[1], v2[1], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[2], v2[2], 1e-6, 1e-6))
}

func (v1 Vec4f) Equal(v2 Vec4f) bool {
	return (scalar.EqualWithinAbsOrRel(v1[0], v2[0], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[1], v2[1], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[2], v2[2], 1e-6, 1e-6) &&
		scalar.EqualWithinAbsOrRel(v1[3], v2[3], 1e-6, 1e-6))
}

func (v1 Vec2i) Equal(v2 Vec2i) bool {
	return v1[0] == v2[0] && v1[1] == v2[1]
}

func (v1 Vec3i) Equal(v2 Vec3i) bool {
	return v1[0] == v2[0] && v1[1] == v2[1] && v1[2] == v2[2]
}

func (v1 Vec4i) Equal(v2 Vec4i) bool {
	return v1[0] == v2[0] && v1[1] == v2[1] && v1[2] == v2[2] && v1[3] == v2[3]
}

func (v Vec2f) ToHomoVec(w float64) mat.VecDense {
  return *mat.NewVecDense(3, []float64 {v[0], v[1], w})
}

func (v Vec3f) ToHomoVec(w float64) mat.VecDense {
  return *mat.NewVecDense(4, []float64 {v[0], v[1], v[2], w})
}

func (v Vec2i) ToHomoVec(w float64) mat.VecDense {
  return *mat.NewVecDense(3, []float64 {float64(v[0]), float64(v[1]), w})
}

func (v Vec3i) ToHomoVec(w float64) mat.VecDense {
  return *mat.NewVecDense(3, []float64 {float64(v[0]), float64(v[1]), float64(v[2]), w})
}

func DenseToVec2f(v *mat.VecDense) (Vec2f, error) {
  if v.Len() != 2 {
    return NewVec2f(), errors.New(fmt.Sprintf("common: wrong size of vec. Got: %d, expected: 2", v.Len()))
  }
  return Vec2f{v.AtVec(0), v.AtVec(1)}, nil
}

func DenseToVec3f(v *mat.VecDense) (Vec3f, error) {
  if v.Len() != 3 {
    return NewVec3f(), errors.New(fmt.Sprintf("common: wrong size of vec. Got: %d, expected: 3", v.Len()))
  }
  return Vec3f{v.AtVec(0), v.AtVec(1), v.AtVec(2)}, nil
}

func DenseToVec4f(v *mat.VecDense) (Vec4f, error) {
  if v.Len() != 4 {
    return NewVec4f(), errors.New(fmt.Sprintf("common: wrong size of vec. Got: %d, expected: 4", v.Len()))
  }
  return Vec4f{v.AtVec(0), v.AtVec(1), v.AtVec(2), v.AtVec(3)}, nil
}

