package connector

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

// A connector for emulator, by default, this connector will output a sine wave

import (
	"context"
	"math"
	"time"
)

type Emulator struct {
	ctx    context.Context
	cancel context.CancelFunc
	Connector
	buf     chan []float32
	waveBuf []float32
}

func NewEmulator(ctxf context.Context) *Emulator {
	ctx, cancel := context.WithCancel(ctxf)
	return &Emulator{
		ctx:    ctx,
		cancel: cancel,
		buf:    make(chan []float32),
	}
}

func (em *Emulator) GetBufferChannel() chan []float32 {
	return em.buf
}

func (em *Emulator) Open() error {
	go em.sinWave()
	go func() {
		tk := time.NewTicker(30 * time.Millisecond)
		for {
			select {
			case <-tk.C:
				em.buf <- em.waveBuf
				em.waveBuf = make([]float32, 0)
			case <-em.ctx.Done():
				tk.Stop()
				return
			}
		}
	}()
	return nil
}

func (em *Emulator) ReadBytes() ([]byte, error) {
	return nil, nil
}

func (em *Emulator) Close() {
	em.cancel()
}

func (em *Emulator) sinWave() {
	count := int64(0)
	for {
		s := math.Sin(float64(count) * math.Pi / 73)
		em.waveBuf = append(em.waveBuf, float32(s))
		//fmt.Println(s)
		time.Sleep(10 * time.Millisecond)
		count++
	}
}

func (em *Emulator) Name() string {
	return "Emulator"
}
