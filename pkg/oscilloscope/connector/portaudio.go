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

// A connector for PortAudio

import (
	"context"
	"fmt"

	"github.com/bhojpur/speech/pkg/portaudio"
)

type PortAudio struct {
	Connector
	ctx    context.Context
	inBuf  chan []float32
	outBuf chan []float32
	host   *portaudio.HostApiInfo
	device *portaudio.StreamParameters
	*portaudio.Stream
	cancel context.CancelFunc
}

func NewPortAudio(ctx context.Context) *PortAudio {
	myCtx, cancel := context.WithCancel(ctx)
	e := &PortAudio{
		inBuf:  make(chan []float32, 1024),
		outBuf: make(chan []float32, 1024),
		ctx:    myCtx,
		cancel: cancel,
	}
	portaudio.Initialize()
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		panic(err)
	}
	d := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	fmt.Println(d.Input.Device.Name)
	d.Input.Channels = 1
	d.Output.Channels = 1
	e.host = h
	e.device = &d
	return e
}

func (al *PortAudio) Open() error {
	go func() {
		s, err := portaudio.OpenStream(*al.device, al.processAudio)
		if err != nil {
			panic(err)
		}
		s.Start()
		al.Stream = s
		<-al.ctx.Done()
		al.Stream.Close()
		err = al.Stream.Stop()
		if err != nil {
			fmt.Println(err)
		}
		err = portaudio.Terminate()
		if err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("port audio open")
	return nil
}

func (al *PortAudio) Close() {

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println("port audio close stream close")
	fmt.Println("port audio close terminate")
	al.cancel()
}

func (al *PortAudio) ReadBytes() ([]byte, error) {
	return nil, nil
}

func (al *PortAudio) GetBufferChannel() chan []float32 {
	return al.inBuf
}

func (al *PortAudio) GetOutPutBufferChannel() chan []float32 {
	return al.outBuf
}

func (al *PortAudio) Info() {

}

func (al *PortAudio) processAudio(in, out []float32) {
	//	fmt.Println("signal come")
	al.inBuf <- in
	//al.outBuf <- out
}

func (al *PortAudio) Name() string {
	return al.device.Input.Device.Name
}
