package main

import (
	"assignment1/common"
	"assignment1/rasterizer"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gonum.org/v1/gonum/mat"
)

// table of trigonometric function
// tableSin[i] = sin(i/ 100)
// tableCos[i] = cos(i/ 100)
// tableCot[i] = cot(i/ 100)
var tableSin, tableCos, tableCot [31416]float64

func newIdentify4() *mat.Dense {
	return mat.NewDense(4, 4, []float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
}

func getViewMatrix(eyePos common.Vec3f) *mat.Dense {
	return mat.NewDense(4, 4, []float64{1, 0, 0, -eyePos[0], 0, 1, 0, -eyePos[1], 0, 0, 1, -eyePos[2], 0, 0, 0, 1})
}

func getModelMatrix(rotationAngle int) *mat.Dense {
	model := newIdentify4()
	model.Set(0, 0, tableCos[rotationAngle])
	model.Set(0, 1, -tableSin[rotationAngle])
	model.Set(1, 0, tableSin[rotationAngle])
	model.Set(1, 1, tableCos[rotationAngle])
    model.Mul(model, mat.NewDense(4, 4, []float64{-1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}))

	return model
}

func getProjectionMatrix(eyeFov, aspectRatio, zNear, zFar float64) *mat.Dense {
	/*
			projection := newIdentify4()
			t := math.Tan(eyeFov/2) * math.Abs(zNear)
			// b = -t
			r := aspectRatio * t
			// l = -r
			projection.Set(0, 0, zNear/r)
			projection.Set(1, 1, zNear/t)
			projection.Set(2, 2, (zNear+zFar)/(zNear-zFar))
			projection.Set(2, 3, -2*zNear*zFar/(zNear-zFar))
		    projection.Set(3, 2, 1) // I have to say: I missed this '1'. :(
	*/
	//cot_fov := 1 / math.Tan(0.5*eyeFov)
	//eyeFov50 := int(eyeFov)<<5 + int(eyeFov)<<4 + int(eyeFov)<<1
	cot_fov := tableCot[int(eyeFov*50)]
	n_f_bw := 1 / (zNear - zFar)
	projection := mat.NewDense(4, 4, []float64{cot_fov / aspectRatio, 0, 0, 0, 0, cot_fov, 0, 0, 0, 0, (zNear + zFar) * n_f_bw, 2 * zNear * zFar * n_f_bw, 0, 0, -1, 0})
	return projection
}

func run() {
	for i := 0; i < 31416; i++ {
		tableSin[i] = math.Sin(float64(i) / 100)
		tableCos[i] = math.Cos(float64(i) / 100)
		tableCot[i] = 1 / math.Tan(float64(i)/100)
	}
	var winWidth, winHeight int = 700, 700
	r := rasterizer.NewRasterizer(winWidth, winHeight, rasterizer.TriangleList)
	eyePos := common.Vec3f{0, 0, 5}

	pos := []common.Vec3f{
		{2, 0, -2},
		{0, 2, -2},
		{-2, 0, -2},
		{3.5, -1, -5},
		{2.5, 1.5, -5},
		{-1, 0.5, -5}}
	ind := []common.Vec3i{{0, 1, 2}, {3, 4, 5}}
	//pos := []common.Vec3f{{2, 0, -2}, {0, 2, -2}, {-2, 0, -2}}
	//ind := []common.Vec3i{{0, 1, 2}}
	//pos := []common.Vec3f{{1, 6, 0}, {3, 6, -1}, {6, 1, -2}, {2, 1, 0}, {1, 2, -1}, {3, 6, -2}, {1, 3, -2}, {6, 3, -1}, {7, 2, 0}}
	//ind := []common.Vec3i{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}}
	angle := 0
	r.SetPrimitive(rasterizer.TriangleList)
	r.LoadVerPosAndInd(pos, ind)

	cfg := pixelgl.WindowConfig{
		Title:  "Rotation",
		Bounds: pixel.R(0, 0, float64(winWidth), float64(winHeight)),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	maxInd := winWidth * winHeight
	var frameCount int64 = 0
	startTime := time.Now().UnixMilli()
	zNear, zFar := 0.1, 50.
	const n = 8016.31794 //float64(winHeight) / (2 * math.Tan(2.5*3.1416/180))
	fmt.Println(2 * math.Atan(float64(winHeight)/(2*n)*180/3.1416))
	eyeFov, aspectRatio := 2*math.Atan(float64(winHeight)/(2*n))*180/3.1416, float64(winWidth)/float64(winHeight)

	for !win.Closed() {
		h := int(win.Bounds().Size().Y)
		w := int(win.Bounds().Size().X)
		if h != winHeight || w != winWidth {
			winWidth = w
			winHeight = h
			r.Resize(w, h)
			eyeFov, aspectRatio = 2*math.Atan(float64(winHeight)/(2*n))*180/3.1416, float64(winWidth)/float64(winHeight)
			maxInd = winWidth * winHeight
		}
		frameCount++
		r.ClearFrameBuf(rasterizer.COLOR | rasterizer.DEPTH)
		win.SetTitle(fmt.Sprintf("Rotation - FrameCount: %d", frameCount))

		r.SetModelMat(getModelMatrix(angle))
		r.SetViewMat(getViewMatrix(eyePos))
		r.SetProjectionMat(getProjectionMatrix(eyeFov, aspectRatio, zNear, zFar))

		go r.Draw()

		frame := r.GetFrameBuf()
		img := image.NewRGBA(image.Rect(0, 0, winWidth, winHeight))

		for i := 0; i < winWidth; i++ {
			for j := 0; j < winHeight; j++ {
				ind := r.GetFrameInd(i, j)
				if ind >= maxInd {
					continue
				}
				frameColor := frame[ind].GetColor()
				img.Set(i, j, color.RGBA{uint8(frameColor[0]), uint8(frameColor[1]), uint8(frameColor[2]), uint8(frameColor[3])})
				//img.Set(i, j, color.RGBA{uint8(frame[ind].GetDepth()), uint8(frame[ind].GetDepth()), uint8(frame[ind].GetDepth()), 255})
			}
		}

		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		go sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		if win.Pressed(pixelgl.KeyA) {
			angle--
		} else if win.Pressed(pixelgl.KeyD) {
			angle++
		}
		if angle == 31416 {
			angle = 0
		}
		if angle == -1 {
			angle = 31415
		}
		frameDuration := time.Now().UnixMilli() - startTime
		if frameDuration >= 1000 {
			fmt.Println("FPS: ", 1000*float64(frameCount)/float64(frameDuration))
			frameCount = 0
			startTime = time.Now().UnixMilli()
		}
		//time.Sleep(5 * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}
