package triangle

import (
	"assignment1/common"
	"errors"
	"fmt"
)

type Triangle struct {
	/* the original coordinates of the triangle, v0, v1, v2 in
	counter clockwise order */
	v [3]common.Vec3f
	/* Per vertex values */
	color      [3]common.Vec4i // color at each vertex;
	tex_coords [3]common.Vec2f // texture u,v
	normal     [3]common.Vec3f // normal common.vector for each vertex
}

func NewTriangle() Triangle {
	return Triangle{
		v:          [3]common.Vec3f{{0., 0., 0.}, {0., 0., 0.}, {0., 0., 0.}},
		color:      [3]common.Vec4i{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		tex_coords: [3]common.Vec2f{{0., 0.}, {0., 0.}, {0., 0.}},
		normal:     [3]common.Vec3f{{0., 0., 0.}, {0., 0., 0.}, {0., 0., 0.}},
	}
}

func (triangle Triangle) GetA() common.Vec3f {
	return triangle.v[0]
}

func (triangle Triangle) GetB() common.Vec3f {
	return triangle.v[1]
}

func (triangle Triangle) GetC() common.Vec3f {
	return triangle.v[2]
}

/* set i-th vertex coordinates */
func (triangle *Triangle) SetVertex(ind int, ver common.Vec3f) error {
	if ind < 0 || ind > 2 {
		return errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
	}

	triangle.v[ind][0] = ver[0]
	triangle.v[ind][1] = ver[1]
	triangle.v[ind][2] = ver[2]

	return nil
}

/* set i-th vertex normal common.vector */
func (triangle *Triangle) SetNormal(ind int, ver common.Vec3f) error {
	if ind < 0 || ind > 2 {
		return errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
	}

	triangle.normal[ind][0] = ver[0]
	triangle.normal[ind][1] = ver[1]
	triangle.normal[ind][2] = ver[2]

	return nil
}

/* set i-th vertex color */
func (triangle *Triangle) SetColor(ind, r, g, b, a int64) error {
	if ind < 0 || ind > 2 {
		return errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0. || b > 255 || a < 0 || a > 255 {
		return errors.New(fmt.Sprintf("triangle: invalid color value (%d, %d, %d, %d)\n", r, g, b, a))
	}
	triangle.color[ind][0] = r
	triangle.color[ind][1] = g
	triangle.color[ind][2] = b
    triangle.color[ind][3] = a
	return nil
}
func (triangle *Triangle) GetColor(ind int) (common.Vec4i, error) {
  if ind < 0 || ind >= 3 {
    return common.NewVec4i(), errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
  }
  return triangle.color[ind], nil
}

/* set i-th vertex texture coordinate */
func (triangle *Triangle) SetTexCoord(ind int, s float64, t float64) error {
	if ind < 0 || ind > 2 {
		return errors.New(fmt.Sprintf("triangle: ind %d out of range\n", ind))
	}

	triangle.tex_coords[ind][0] = s
	triangle.tex_coords[ind][1] = t

	return nil
}

func (triangle *Triangle) ToVec4() [3]common.Vec4f {
	res := [3]common.Vec4f{
		{triangle.v[0][0], triangle.v[0][1], triangle.v[0][2], 1.},
		{triangle.v[1][0], triangle.v[1][1], triangle.v[1][2], 1.},
		{triangle.v[2][0], triangle.v[2][1], triangle.v[2][2], 1.},
	}
	return res
}
