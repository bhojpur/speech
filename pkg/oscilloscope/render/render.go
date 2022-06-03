package render

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/bhojpur/speech/pkg/oscilloscope/connector"
	"github.com/eiannone/keyboard"
	"github.com/fogleman/gg"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

var (
	Red   = image.NewUniform(color.RGBA{0, 255, 0, 128})
	Green = image.NewUniform(color.RGBA{0, 255, 0, 128})
	Blue  = image.NewUniform(color.RGBA{0, 0, 255, 128})
)

type Render struct {
	ctx             context.Context
	width           int
	height          int
	activeWidth     int
	activeHeight    int
	buffer          []float32
	normalBuffer    []float32
	window          screen.Window
	windowBuffer    screen.Buffer
	connector       connector.Connector
	zoomOffsetStart int
	zoomParam       float32
	pause           bool
	quit            bool
	cancel          context.CancelFunc
}

func NewRender(ctx context.Context, w, h int, conn connector.Connector) *Render {
	err := conn.Open()
	if err != nil {
		panic(err)
	}
	myCtx, cancel := context.WithCancel(ctx)
	return &Render{
		ctx:          myCtx,
		width:        w,
		height:       h,
		activeWidth:  w,
		activeHeight: h,
		buffer:       make([]float32, 0),
		normalBuffer: make([]float32, 0),
		connector:    conn,
		zoomParam:    1,
		cancel:       cancel,
	}
}

func (rd *Render) Close() {
	fmt.Println("start render close")
	rd.connector.Close()
	keyboard.Close()
	rd.quit = true
	fmt.Println("render close")
}

func (rd *Render) ResizeWindow(w, h int) {
	rd.width = w
	rd.height = h
	rd.activeWidth = w
	rd.activeHeight = h * 7 / 8
}

func (rd *Render) Start() {
	// create a winow
	driver.Main(func(src screen.Screen) {
		win, _ := src.NewWindow(&screen.NewWindowOptions{
			Width:  rd.width,
			Height: rd.height,
			Title:  "goosilloscope",
		})
		rd.window = win
		rd.windowBuffer, _ = src.NewBuffer(image.Point{X: rd.width, Y: rd.height})
		// winow events
		go func() {
			rd.windowEnventHandler(src)
		}()

		go func() {
			rd.keyboardEventHandler()
		}()
		go func() {
			rd.dataEventHandler()
		}()

		<-rd.ctx.Done()
		fmt.Println("exit progress")
		rd.Close()
	})
}

func (rd *Render) draw(min, max float32) {
	dc := gg.NewContext(rd.width, rd.height)

	dc.SetLineWidth(5)
	rd.MakeBoarder(dc)
	rd.MakeBaseLines(dc, min, max)
	rd.MakeStatusBar(dc, min, max)
	rd.MakeInfoBar(dc)
	lastX := float64(0)
	lastY := float64(0)

	dc.SetRGBA(128.0, 64.0, 255.0, 220)
	for i := 0; i < len(rd.normalBuffer); i++ {
		dc.DrawLine(lastX, lastY, float64(i), float64(rd.normalBuffer[i]))
		dc.Stroke()
		lastX = float64(i)
		lastY = float64(rd.normalBuffer[i])
	}

	img := dc.Image()
	// dc.SavePNG("./result.png")
	draw.Draw(rd.windowBuffer.RGBA(), img.Bounds(), img, image.Point{}, draw.Src)
	rd.window.Upload(img.Bounds().Min, rd.windowBuffer, img.Bounds())
	rd.window.Publish()
}
