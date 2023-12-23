package rasterizer

import (
	"assignment1/common"
	"assignment1/triangle"
	"errors"
	"fmt"
	"math"
	"sync"

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
	color    common.Vec4i // (r, g, b, a) => 0 <= r, g, b, a <= 255
	// TODO: Add more elements
	// for example:
	/*
	  normalLine    common.Vec3f // (nx, ny, nz)
	  texture       common.Vec2f // (u, v)
	*/
}

func NewVertexBufferElement() VertexBufferElement {
	return VertexBufferElement{
		position: common.NewVec3f(),
		color:    common.NewVec4i(),
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
	color common.Vec4i // (r, g, b, a) => 0 <= r, g, b, a <= 255
	depth float64
	// TODO: Add more elements
	// for example: tencil, accumulation...
}

func NewFrameBufferElement() FrameBufferElement {
	return FrameBufferElement{
		color: common.Vec4i{0, 0, 0, 255},
		depth: math.Inf(1),
	}
}

func (fe *FrameBufferElement) GetColor() common.Vec4i {
	return fe.color
}

func (fe *FrameBufferElement) GetDepth() float64 {
	return fe.depth
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
		frameBuf:      make([]FrameBufferElement, (w*h)<<2),
		primitive:     primitive,
		IndBufInd:     0,
		VerBufInd:     0,
	}
}

func (r *Rasterizer) Resize(w, h int) {
	r.width = w
	r.height = h
	r.frameBuf = make([]FrameBufferElement, (w*h)<<2)
}

func (r *Rasterizer) GetSize() (int, int) {
	return r.width, r.height
}

func (r *Rasterizer) SetPrimitive(primitive PrimitiveType) {
	r.primitive = primitive
}

// Only load the position of the vertexs.
func (r *Rasterizer) LoadVer(positions []common.Vec3f, colors []common.Vec4i) {
	if r.VerBufInd != len(positions) {
		r.VerBufInd = len(positions)
		r.vertexBuf = make([]VertexBufferElement, r.VerBufInd)
	}
	for i, pos := range positions {
		r.vertexBuf[i].position = pos
		r.vertexBuf[i].color = colors[i]
	}

}
func (r *Rasterizer) LoadInd(indices []common.Vec3i) {
	if r.IndBufInd != len(indices) {
		r.IndBufInd = len(indices)
		r.indexBuf = make([]common.Vec3i, r.IndBufInd)
	}
	for i, ind := range indices {
		r.indexBuf[i] = ind
	}
}

func (r *Rasterizer) ClearFrameBuf(signal ClearSignal) {
	if signal == 0 {
		return
	}
	var wg sync.WaitGroup
	/*for i := 0; i < r.width*r.height; {
		wg.Add(1)
		r.clearPerRow(i, i+r.width, signal, &wg)
		i += r.width
	}
	wg.Wait()*/
	if (signal & COLOR) == COLOR {
		wg.Add(1)
		go r.clearColor(&wg)
	}
	if (signal & DEPTH) == DEPTH {
		wg.Add(1)
		go r.clearDepth(&wg)
	}
	if ((signal & COLOR) == COLOR) || ((signal & DEPTH) == DEPTH) {
		wg.Wait()
	}
}

func (r *Rasterizer) clearColor(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(r.frameBuf); i++ {
		r.frameBuf[i].color = common.Vec4i{0, 0, 0, 255}
	}
}

func (r *Rasterizer) clearDepth(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < len(r.frameBuf); i++ {
		r.frameBuf[i].depth = math.Inf(1)
	}
}

/*func (r *Rasterizer) clearPerRow(begin, end int, signal ClearSignal, wg *sync.WaitGroup) {
	defer wg.Done()
	if (signal & COLOR) == COLOR {
		for i := begin; i < end && i < r.width*r.height; i++ {
			r.frameBuf[i].color = common.Vec4i{0, 0, 0, 255}
		}
	}
	if (signal & DEPTH) == DEPTH {
		for i := begin; i < end && i < r.width*r.height; i++ {
			r.frameBuf[i].depth = math.Inf(1)
		}
	}
}*/

// map (x, y) of point to (x', y') of screen
func (r *Rasterizer) GetFrameInd(x, y int) int {
	return x*r.width + y
}

func (r *Rasterizer) setPixel(ind int, depth float64, color common.Vec4i) {
	if ind < 0 || ind >= len(r.frameBuf) {
		return
	}
	if depth < r.frameBuf[ind].depth {
		r.frameBuf[ind].color = color
		r.frameBuf[ind].depth = depth
	}
}

/*
func (r *Rasterizer) setPixel(point common.Vec3f, color common.Vec4i) {
	ind := r.GetFrameInd(int(point[0]), int(point[1]))
	if ind < 0 || ind >= len(r.frameBuf) {
		return
	}
	if point[2] < r.frameBuf[ind].depth {
		r.frameBuf[ind].color = color
		r.frameBuf[ind].depth = point[2]
	}
}
*/

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

/*
func (r *Rasterizer) drawLine(begin, end common.Vec3f, lineColor common.Vec4i) {
	x1 := int(begin[0])
	y1 := int(begin[1])
	x2 := int(end[0])
	y2 := int(end[1])
	dx := x2 - x1
	dy := y2 - y1
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

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
		point := common.Vec3f{float64(x1), float64(y1), 1}
		r.setPixel(point, lineColor)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := err << 1
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
*/

// Create bounding box
// for each pixel p in bounding box:
//
//	if p is inside the triangle:
//	  interpolated depth value of p
//	  if depth of p < depth in frame buffer:
//	    draw pixel
//	    update depth in buffer
func (r *Rasterizer) rasterizeTriangle(t *triangle.Triangle) {
	// find BoundingBox
	v := t.GetVertxs()
	minX := int(math.Min(v[0][0], math.Min(v[1][0], v[2][0])))
	minY := int(math.Min(v[0][1], math.Min(v[1][1], v[2][1])))
	maxX := int(math.Max(v[0][0], math.Max(v[1][0], v[2][0])))
	maxY := int(math.Max(v[0][1], math.Max(v[1][1], v[2][1])))

	var wg sync.WaitGroup
	for x := minX; x <= maxX; x++ {
		wg.Add(1)
		go r.rasterizeTriangleCol(x, minY, maxY, t, &wg)
	}
	wg.Wait()
}

func (r *Rasterizer) rasterizeTriangleCol(x, minY, maxY int, t *triangle.Triangle, wg *sync.WaitGroup) {
	defer wg.Done()
	for y := minY; y < maxY; y++ {
		v := t.GetVertxs()
		/*if insideTriangle(float64(x)+0.5, float64(y)+0.5, v) {
			alpha, beta, gamma := computeBarycentric2D(float64(x), float64(y), v)
			z_interpolated := (alpha*v[0][2] + beta*v[1][2] + gamma*v[2][2]) / (alpha + beta + gamma)
			r.setPixel(common.Vec3f{float64(x), float64(y), -z_interpolated}, t.GetColor(0))
		}*/
		// MSAA4
		alpha, beta, gamma := computeBarycentric2D(float64(x), float64(y), v)
		z_interpolated := (alpha*v[0][2] + beta*v[1][2] + gamma*v[2][2]) / (alpha + beta + gamma)
		ind := (r.GetFrameInd(x, y)) << 2
		color := t.GetColor(0)
		if insideTriangle(float64(x)+0.25, float64(y)+0.25, v) {
			r.setPixel(ind, z_interpolated, color)
		}
		if insideTriangle(float64(x)+0.25, float64(y)+0.75, v) {
			r.setPixel(ind+1, z_interpolated, color)
		}
		if insideTriangle(float64(x)+0.75, float64(y)+0.25, v) {
			r.setPixel(ind+2, z_interpolated, color)
		}
		if insideTriangle(float64(x)+0.75, float64(y)+0.75, v) {
			r.setPixel(ind+3, z_interpolated, color)
		}
	}
}

func (r *Rasterizer) Draw() error {
	if r.primitive != TriangleList {
		return errors.New("rasterizer: drawing primitives other than triangle is not implemented yet!")
	}
	vers := &r.vertexBuf
	indices := &r.indexBuf

	mvp := mat.NewDense(4, 4, nil)
	mvp.Mul(r.projectionMat, r.viewMat)
	mvp.Mul(mvp, r.modelMat)

	for _, ind := range *indices {
		t := triangle.NewTriangle()
		v1 := (*vers)[ind[0]].position.ToHomoVec(1.)
		v2 := (*vers)[ind[1]].position.ToHomoVec(1.)
		v3 := (*vers)[ind[2]].position.ToHomoVec(1.)

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

		//f1 := (100 - 0.1) / 2.0
		//f2 := (100 + 0.1) / 2.0
		for i := 0; i < 3; i++ {
			//Homogeneous division
			w_backwards := 1 / v[i][3]
			v[i][0] *= w_backwards
			v[i][1] *= w_backwards
			v[i][2] *= w_backwards
			v[i][3] = 1

			// Viewport Transformation
			v[i][0] = float64(r.width>>1) * (v[i][0] + 1)
			v[i][1] = float64(r.height>>1) * (v[i][1] + 1)
			//v[i][2] = v[i][2]*f1 + f2
			//Depth Range Mapping:
			v[i][2] = v[i][2]*49.95 + 50.05
			t.SetVertex(i, common.Vec3f(v[i][:3]))
		}

		t.SetColor(0, (*vers)[ind[0]].color[0], (*vers)[ind[0]].color[1], (*vers)[ind[0]].color[2], (*vers)[ind[0]].color[3])
		t.SetColor(1, (*vers)[ind[1]].color[0], (*vers)[ind[1]].color[1], (*vers)[ind[1]].color[2], (*vers)[ind[1]].color[3])
		t.SetColor(2, (*vers)[ind[2]].color[0], (*vers)[ind[2]].color[1], (*vers)[ind[2]].color[2], (*vers)[ind[2]].color[3])

		r.rasterizeTriangle(&t)
	}
	return nil
}

// OTHER FUNCTIONS
func computeBarycentric2D(x, y float64, v [3]common.Vec3f) (float64, float64, float64) {
	c1 := (x*(v[1][1]-v[2][1]) + (v[2][0]-v[1][0])*y + v[1][0]*v[2][1] - v[2][0]*v[1][1]) / (v[0][0]*(v[1][1]-v[2][1]) + (v[2][0]-v[1][0])*v[0][1] + v[1][0]*v[2][1] - v[2][0]*v[1][1])
	c2 := (x*(v[2][1]-v[0][1]) + (v[0][0]-v[2][0])*y + v[2][0]*v[0][1] - v[0][0]*v[2][1]) / (v[1][0]*(v[2][1]-v[0][1]) + (v[0][0]-v[2][0])*v[1][1] + v[2][0]*v[0][1] - v[0][0]*v[2][1])
	c3 := (x*(v[0][1]-v[1][1]) + (v[1][0]-v[0][0])*y + v[0][0]*v[1][1] - v[1][0]*v[0][1]) / (v[2][0]*(v[0][1]-v[1][1]) + (v[1][0]-v[0][0])*v[2][1] + v[0][0]*v[1][1] - v[1][0]*v[0][1])
	return c1, c2, c3
}

func insideTriangle(x, y float64, v [3]common.Vec3f) bool {
	cross1 := (v[1][0]-v[0][0])*(y-v[0][1]) - (v[1][1]-v[0][1])*(x-v[0][0])
	cross2 := (v[2][0]-v[1][0])*(y-v[1][1]) - (v[2][1]-v[1][1])*(x-v[1][0])
	cross3 := (v[0][0]-v[2][0])*(y-v[2][1]) - (v[0][1]-v[2][1])*(x-v[2][0])

	return (cross1 >= 0 && cross2 >= 0 && cross3 >= 0) || (cross1 <= 0 && cross2 <= 0 && cross3 <= 0)
}
