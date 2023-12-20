package rasterizer

import (
	"assignment1/common"
	"assignment1/triangle"
	"errors"
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

type ClearSignal int64

const (
	COLOR ClearSignal = iota
	DEPTH ClearSignal = iota
)

// STRUCT OF VertexBufferElement AND ITS CONSTRUCTOR
type VertexBufferElement struct {
	position common.Vec3f // (x, y, z)
	color    common.Vec4f // (r, g, b, a) => 0 <= r, g, b, a <= 1
	// TODO: Add more buffers
	// for example:
	/*
	  normalLine    []common.Vec3f // (nx, ny, nz)
	  texture       []common.Vec2f // (u, v)
	*/
}

func NewVertexBufferElement() VertexBufferElement {
	return VertexBufferElement{
		position: common.NewVec3f(),
		color:    common.NewVec4f(),
	}
}

// STRUCT OF IndexBufferElement AND ITS CONSTRUCTOR
// Type of Primitive
type PrimitiveType int

const (
	LineList     PrimitiveType = iota
	TriangleList PrimitiveType = iota
	// TODO: Add more primitive types
	// for example:
	/*
	  LineStrip PrimitiveType = iota
	  TriangleStrip PrimitiveType = iota
	  TriangleFan PrimitiveType = iota
	*/
)

// STRUCT OF FrameBufferElement AND ITS CONSTROCTOR
type FrameBufferElement struct {
	color common.Vec4f // (r, g, b, a) => 0 <= r, g, b, a <= 1
	depth float64
	// TODO: Add more buffers
	// for example: tencil, accumulation...
}

func NewFrameBufferElement() FrameBufferElement {
	return FrameBufferElement{
		color: common.NewVec4f(),
		depth: math.Inf(1),
	}
}

func (fe *FrameBufferElement) GetColor() common.Vec4f {
	return fe.color
}

// Rasterizer type
type Rasterizer struct {
	// Screen size: width x height pixels
	width  int // num of pixels per row
	height int // num of pixels per col
	// Three Matrix: MVP
	modelMat      *mat.Dense // Mat(m), size: 4x4
	viewMat       *mat.Dense // Mat(v), size: 4x4
	projectionMat *mat.Dense // Mat(p), size: 4x4

	// Some buffers
	// TODO: Add more buffers
	// for example: constBuffer, texture Buffer, pixel Buffer...
	vertexBuf []VertexBufferElement
	indexBuf  []common.Vec3i // NOTE THAT Vec3i is wrong in a more general situation. Beacuse each element in indexBuf (a vector of n dimention) refers to n vertexs in vertexBuf.
	frameBuf  []FrameBufferElement
	primitive PrimitiveType

	// index is used for vertexBuf and indexBuf
	IndBufInd int
	VerBufInd int
}

func NewRasterizer(w, h int, primitive PrimitiveType) Rasterizer {
	return Rasterizer{
		width:         w,
		height:        h,
		modelMat:      mat.NewDense(4, 4, nil),
		viewMat:       mat.NewDense(4, 4, nil),
		projectionMat: mat.NewDense(4, 4, nil),
		vertexBuf:     make([]VertexBufferElement, 1),
		indexBuf:      make([]common.Vec3i, 1),
		frameBuf:      make([]FrameBufferElement, w*h),
		primitive:     primitive,
		IndBufInd:     0,
		VerBufInd:     0,
	}
}

func (r *Rasterizer) GetSize() (int, int) {
	return r.width, r.height
}

func (r *Rasterizer) SetPrimitive(primitive PrimitiveType) {
	r.primitive = primitive
}

// Only load the position of the vertexs.
func (r *Rasterizer) LoadVerPosAndInd(positions []common.Vec3f, indices []common.Vec3i) {
	r.VerBufInd = len(positions)
	r.IndBufInd = len(indices)
	r.vertexBuf = make([]VertexBufferElement, r.VerBufInd)
	r.indexBuf = make([]common.Vec3i, r.IndBufInd)
	for i := 0; i < r.VerBufInd; i++ {
		r.vertexBuf[i].position = positions[i]
	}
	for i := 0; i < r.IndBufInd; i++ {
		r.indexBuf[i] = indices[i]
	}
}

func (r *Rasterizer) ClearFrameBuf(signal ClearSignal) {
	if signal == 0 {
		return
	}
	if (signal & COLOR) == COLOR {
		for i := 0; i < len(r.frameBuf); i++ {
			r.frameBuf[i].color = common.NewVec4f()
		}
	}
	if (signal & DEPTH) == DEPTH {
		for i := 0; i < len(r.frameBuf); i++ {
			r.frameBuf[i].depth = math.Inf(1)
		}
	}
}

// map (x, y) of point to (x', y') of screen
func (r *Rasterizer) GetFrameInd(x, y int) int {
	return x* r.width+ y
}

func (r *Rasterizer) setPixel(point common.Vec3f, color common.Vec4f) error {
	x := point[0]
	y := point[1]
	if x < 0 || x >= float64(r.width) || y < 0 || y >= float64(r.height) {
		return errors.New(fmt.Sprintf("rasterizer: wrong point. Got (%.3f, %.3f), expected range: (0, 0) to (%d, %d)", x, y, r.width, r.height))
	}
	x = math.Floor(x)
	y = math.Floor(y)
	ind := r.GetFrameInd(int(x), int(y))
	if ind >= len(r.frameBuf) {
		return errors.New(fmt.Sprintf("rasterizer: index out of range. Got: %d, max: %d", ind, r.width*r.height))
	}
	r.frameBuf[ind].color = color
	return nil
}

func (r *Rasterizer) GetFrameBuf() []FrameBufferElement {
	return r.frameBuf
}

// SETTERs OF MATS
func (r *Rasterizer) SetModelMat(m *mat.Dense) error {
	if r, c := m.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong dimension of model matrix. Got: %dx%d, expected: 4x4", r, c))
	}
	r.modelMat = m
	return nil
}

func (r *Rasterizer) SetViewMat(v *mat.Dense) error {
	if r, c := v.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong dimension of view matrix. Got: %dx%d, expected: 4x4", r, c))
	}
	r.viewMat = v
	return nil
}

func (r *Rasterizer) SetProjectionMat(p *mat.Dense) error {
	if r, c := p.Dims(); r != 4 || c != 4 {
		return errors.New(fmt.Sprintf("rasterizer: wrong dimension of projection matrix. Got: %dx%d, expected: 4x4", r, c))
	}
	r.projectionMat = p
	return nil
}

// METHODS ABOUT DRAWING
func (r *Rasterizer) drawLine(begin, end common.Vec3f, lineColor common.Vec4f) {
	x1 := int(begin[0])
	y1 := int(begin[1])
	x2 := int(end[0])
	y2 := int(end[1])

	dx := int(math.Abs(float64(x2 - x1)))
	dy := int(math.Abs(float64(y2 - y1)))
	var sx, sy int

	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		point := common.Vec3f{float64(x1), float64(y1), 1.}
		r.setPixel(point, lineColor)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func (r *Rasterizer) rasterizeWireframe(t *triangle.Triangle) {
	color, _ := t.GetColor(0)
	r.drawLine(t.GetC(), t.GetA(), color)
	color, _ = t.GetColor(1)
	r.drawLine(t.GetC(), t.GetB(), color)
	color, _ = t.GetColor(2)
	r.drawLine(t.GetB(), t.GetA(), color)
}

func (r *Rasterizer) Draw() error {
	if r.primitive != TriangleList {
		return errors.New("rasterizer: drawing primitives other than triangle is not implemented yet!")
	}
	vers := r.vertexBuf
	indices := r.indexBuf

	mvp := mat.NewDense(4, 4, nil)
	mvp.Mul(r.projectionMat, r.viewMat)
	mvp.Mul(mvp, r.modelMat)

	for _, ind := range indices {
		t := triangle.NewTriangle()
		v1 := vers[ind[0]].position.ToHomoVec(1.)
		v2 := vers[ind[1]].position.ToHomoVec(1.)
		v3 := vers[ind[2]].position.ToHomoVec(1.)

		v1.MulVec(mvp, &v1)
		v2.MulVec(mvp, &v2)
		v3.MulVec(mvp, &v3)
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

		f1 := (100 - 0.1) / 2.0
		f2 := (100 + 0.1) / 2.0
		for i := 0; i < 3; i++ {
			v[i][0] /= v[i][3]
			v[i][1] /= v[i][3]
			v[i][2] /= v[i][3]
			v[i][3] = 1.

			v[i][0] = 0.5 * float64(r.width) * (v[i][0] + 1)
			v[i][1] = 0.5 * float64(r.height) * (v[i][1] + 1)
			v[i][2] = v[i][2]*f1 + f2
			t.SetVertex(i, common.Vec3f(v[i][:3]))
		}

		t.SetColor(0, 255, 0, 0, 255)
		t.SetColor(1, 0, 255, 0, 255)
		t.SetColor(2, 0, 0, 255, 255)

		r.rasterizeWireframe(&t)
	}
	return nil
}
