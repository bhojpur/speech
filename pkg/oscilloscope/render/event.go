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
	"fmt"
	"image"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/size"
)

func (rd *Render) windowEnventHandler(src screen.Screen) {
	fmt.Println("window event handler open")
	for !rd.quit {
		switch e := rd.window.NextEvent().(type) {
		case size.Event:
			if e.WidthPx == 0 && e.HeightPx == 0 {
				fmt.Println("exit progress")
				rd.quit = true
				rd.cancel()
			} else {
				rd.ResizeWindow(e.WidthPx, e.HeightPx)
				//fmt.Println("fix size: [", rd.height, rd.width, "]")
				rd.windowBuffer, _ = src.NewBuffer(image.Point{X: rd.width, Y: rd.height})
			}
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				fmt.Println("cancel done")
				rd.quit = true
				rd.cancel()
			}
		}
	}
}

func (rd *Render) keyboardEventHandler() {

	// keyboard event
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	fmt.Println("keyboard event handler open")
	for !rd.quit {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}
		fmt.Println(key, err, char)
		switch key {
		case keyboard.KeyArrowUp:
			rd.ZoomIn()
		case keyboard.KeyArrowDown:
			rd.ZoomOut()
		case keyboard.KeyEsc:
			rd.quit = true
			rd.cancel()
		case keyboard.KeyCtrlC:
			rd.quit = true
			rd.cancel()
		case 32:
			rd.Pause()
		}
	}

}

func (rd *Render) dataEventHandler() {
	// read & draw process
	tk := time.NewTicker(1 * time.Second)
	fmt.Println("data event handler open")
	min, max := rd.normalizeValue()
	rd.draw(min, max)
	ch := rd.connector.GetBufferChannel()
	for !rd.quit {
		// pause
		if rd.pause {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		select {
		case b := <-ch:
			for _, v := range b {
				rd.addBuffer(v)
			}
		case <-tk.C:
		}

		min, max := rd.normalizeValue()
		rd.draw(min, max)
	}

}
