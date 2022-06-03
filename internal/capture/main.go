package main

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

// It captures data from your default microphone until you press Enter, after which
// it plays back the captured audio.

import (
	"fmt"
	"os"

	engine "github.com/bhojpur/speech/pkg/miniaudio"
)

func main() {
	fmt.Println("Bhojpur Speech capture utility")
	fmt.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	fmt.Printf("All rights reserved.\n")

	ctx, err := engine.InitContext(nil, engine.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := engine.DefaultDeviceConfig(engine.Duplex)
	deviceConfig.Capture.Format = engine.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = engine.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	var playbackSampleCount uint32
	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(engine.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {

		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount

		pCapturedSamples = append(pCapturedSamples, pSample...)

		capturedSampleCount = newCapturedSampleCount

	}

	fmt.Println("Please speak now, I am recording your voice...")
	captureCallbacks := engine.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := engine.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press ENTER key to stop recording...")
	fmt.Scanln()

	device.Uninit()

	onSendFrames := func(pSample, nil []byte, framecount uint32) {
		samplesToRead := framecount * deviceConfig.Playback.Channels * sizeInBytes
		if samplesToRead > capturedSampleCount-playbackSampleCount {
			samplesToRead = capturedSampleCount - playbackSampleCount
		}

		copy(pSample, pCapturedSamples[playbackSampleCount:playbackSampleCount+samplesToRead])

		playbackSampleCount += samplesToRead

		if playbackSampleCount == uint32(len(pCapturedSamples)) {
			playbackSampleCount = 0
		}
	}

	fmt.Println("Now, I am playing your recorded voice...")
	playbackCallbacks := engine.DeviceCallbacks{
		Data: onSendFrames,
	}

	device, err = engine.InitDevice(ctx.Context, deviceConfig, playbackCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Press ENTER key to quit this program...")
	fmt.Scanln()

	device.Uninit()
}
