package rasterizer

import (
	"assignment1/common"
	"assignment1/triangle"
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

/* declare basic types */
type Buffers int64

const (
	Color Buffers = 1
	Depth Buffers = 2
)

type Primitive int64

const (
	Line Primitive = iota
	Triangle
)

type posBufId struct {
	posId int64
}

func NewPosBufId(id int64) posBufId {
	return posBufId{
		posId: id,
	}
}

type indBufId struct {
	indId int64
}

func NewIndBufId(id int64) indBufId {
	return indBufId{
		indId: id,
	}
}

/* declare Rasterizer */
type Rasterizer struct {
	model      *mat.Dense
	view       *mat.Dense
	projection *mat.Dense

	posBuf map[int64][]common.Vec3f
	indBuf map[int64][]common.Vec3i

	frameBuf []common.Vec3f
	depthBuf []float64

	width  int64
	height int64
	nextId int64
}

func NewRasterizer(w, h int64) Rasterizer {
	return Rasterizer{
		model:      mat.NewDense(4, 4, []float64{}),
		view:       mat.NewDense(4, 4, []float64{}),
		projection: mat.NewDense(4, 4, []float64{}),

		posBuf: make(map[int64][]common.Vec3f),
		indBuf: make(map[int64][]common.Vec3i),

		frameBuf: []common.Vec3f{},
		depthBuf: []float64{},

		width:  w,
		height: h,

		nextId: 0,
	}
}

/* methods of Rasterizer  */
func (rasterizer *Rasterizer) Clear(buff Buffers) {
	if (buff & Color) == Color {
		for i := 0; i < len(rasterizer.frameBuf); i++ {
			rasterizer.frameBuf[i] = common.NewVec3f()
		}
	}
	if (buff & Depth) == Depth {
		for i := 0; i < len(rasterizer.depthBuf); i++ {
			rasterizer.depthBuf[i] = math.Inf(1)
		}
	}
}

/* some setter methods */
func (rasterizer *Rasterizer) SetModel(m *mat.Dense) error {
	if r, c := m.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong size of matrix. Got: (%d, %d), expected: (4, 4)", r, c))
	}
	rasterizer.model.CloneFrom(m)
	return nil
}

func (rasterizer *Rasterizer) SetView(v *mat.Dense) error {
	if r, c := v.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong size of matrix. Got: (%d, %d), expected: (4, 4)", r, c))
	}
	rasterizer.view.CloneFrom(v)
	return nil
}

func (rasterizer *Rasterizer) SetProjection(p *mat.Dense) error {
	if r, c := p.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong size of matrix. Got: (%d, %d), expected: (4, 4)", r, c))
	}
	rasterizer.projection.CloneFrom(p)
	return nil
}

/* some getter methods */
func (rasterizer *Rasterizer) GetFrameBuf() []common.Vec3f {
	return rasterizer.frameBuf
}

func (rasterizer *Rasterizer) GetIndex(x, y int64) int64 {
	return (rasterizer.height-y)*rasterizer.width + x
}

func (rasterizer *Rasterizer) GetNextId() int64 {
	t := rasterizer.nextId
	rasterizer.nextId++
	return t
}

func (rasterizer *Rasterizer) LoadPositions(positions []common.Vec3f) posBufId {
	id := rasterizer.GetNextId()
	rasterizer.posBuf[id] = positions
	return NewPosBufId(id)
}

func (rasterizer *Rasterizer) LoadIndices(indices []common.Vec3i) indBufId {
	id := rasterizer.GetNextId()
	rasterizer.indBuf[id] = indices
	return NewIndBufId(id)
}

/* some methods about drawing */
func (rasterizer *Rasterizer) SetPixel(point common.Vec3f, color common.Vec3f) error {
	if point[0] < 0 || point[0] >= float64(rasterizer.width) || point[1] < 0 || point[1] >= float64(rasterizer.height) {
		return errors.New(fmt.Sprintf("rasterizer: wrong range of point. Got: (%f, %f), expected: from (0, 0) to (%d, %d)", point[0], point[1], rasterizer.width, rasterizer.height))
	}
	ind := int64((float64(rasterizer.height)-point[1])*float64(rasterizer.width) + point[0])
	rasterizer.frameBuf[ind] = color
	return nil
}

func (rasterizer *Rasterizer) drawLine(begin, end common.Vec3f) {
	x1 := begin[0]
	y1 := begin[1]
	x2 := end[0]
	y2 := end[1]

	lineColor := common.Vec3f{255., 255., 255.}

	var x, y, xe, ye, dx, dy, dx1, dy1, px, py float64

	dx = x2 - x1
	dy = y2 - y1
	dx1 = math.Abs(dx)
	dy1 = math.Abs(dy)
	px = 2*dy1 - dx1
	py = 2*dx1 - dy1

	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
		}
		point := common.Vec3f{x, y, 1.}
		rasterizer.SetPixel(point, lineColor)
		for i := 0; x < xe; i++ {
			x = x + 1
			if px < 0 {
				px = px + 2*dy1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y = y + 1
				} else {
					y = y - 1
				}
				px = px + 2*(dy1-dx1)
			}
			point := common.Vec3f{x, y, 1.}
			rasterizer.SetPixel(point, lineColor)
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
		}
		point := common.Vec3f{x, y, 1.}
		rasterizer.SetPixel(point, lineColor)
		for i := 0; y < ye; i++ {
			y = y + 1
			if py <= 0 {
				py = py + 2*dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x = x + 1
				} else {
					x = x - 1
				}
				py = py + 2*(dx1-dy1)
			}
			point := common.Vec3f{x, y, 1.}
			rasterizer.SetPixel(point, lineColor)
		}
	}
}

func (rasterizer *Rasterizer) rasterizeWireframe(t *triangle.Triangle) {
	rasterizer.drawLine(t.GetC(), t.GetA())
	rasterizer.drawLine(t.GetC(), t.GetB())
	rasterizer.drawLine(t.GetB(), t.GetA())
}

func (rasterizer Rasterizer) Draw(posBuffer posBufId, indBuffer indBufId, _type Primitive) error {
	if _type != Triangle {
		return errors.New("Drawing primitives other than triangle is not implemented yet!")
	}
	buf := rasterizer.posBuf[posBuffer.posId]
	ind := rasterizer.indBuf[indBuffer.indId]

	var f1 float64 = (100. - 0.1) / 2.0
	var f2 float64 = (100. + 0.1) / 2.0

	var mvp mat.Dense
	mvp.Mul(rasterizer.projection, rasterizer.view)
	mvp.Mul(&mvp, rasterizer.model)

	for _, i := range ind {
		t := triangle.NewTriangle()

		v1 := buf[i[0]].ToHomoVec(1.)
		v2 := buf[i[1]].ToHomoVec(1.)
		v3 := buf[i[2]].ToHomoVec(1.)

		v1.MulVec(&mvp, &v1)
		v2.MulVec(&mvp, &v2)
		v3.MulVec(&mvp, &v3)

		v1f, err := common.DenseToVec4f(&v1)
		if err != nil {
			return err
		}
		v2f, err := common.DenseToVec4f(&v2)
		if err != nil {
			return err
		}
		v3f, err := common.DenseToVec4f(&v3)
		if err != nil {
			return err
		}
		v := [3]common.Vec4f{v1f, v2f, v3f}

		for j := 0; j < 3; j++ {
			v[j][0] /= v[j][3]
			v[j][1] /= v[j][3]
			v[j][2] /= v[j][3]
			v[j][3] = 1.

			v[j][0] = 0.5*float64(rasterizer.width)*v[j][0] + 1.
			v[j][1] = 0.5*float64(rasterizer.height)*v[j][1] + 1.
			v[j][2] = v[j][2] * f1 * f2

			t.SetVertex(j, common.Vec3f(v[j][:3]))
		}

		t.SetColor(0, 255.0, 0.0, 0.0)
		t.SetColor(1, 0.0, 255.0, 0.0)
		t.SetColor(2, 0.0, 0.0, 255.0)

		rasterizer.rasterizeWireframe(&t)
	}
	return nil
}
