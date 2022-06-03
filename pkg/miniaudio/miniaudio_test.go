package miniaudio_test

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
	"io/ioutil"
	"testing"
	"time"

	engine "github.com/bhojpur/speech/pkg/miniaudio"
)

func TestCapturePlayback(t *testing.T) {
	onLog := func(message string) {
		fmt.Fprintf(ioutil.Discard, message)
	}

	ctx, err := engine.InitContext(nil, engine.ContextConfig{}, onLog)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := engine.DefaultDeviceConfig(engine.Capture)
	deviceConfig.Capture.Format = engine.FormatS16
	deviceConfig.Capture.Channels = 2
	deviceConfig.Playback.Format = engine.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	var playbackSampleCount uint32
	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(engine.SampleSizeInBytes(deviceConfig.Playback.Format))
	onRecvFrames := func(outpuSamples, inputSamples []byte, framecount uint32) {
		sampleCount := framecount * deviceConfig.Playback.Channels * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount

		pCapturedSamples = append(pCapturedSamples, inputSamples...)

		capturedSampleCount = newCapturedSampleCount
	}

	captureCallbacks := engine.DeviceCallbacks{
		Data: onRecvFrames,
	}
	captureDeviceConfig := deviceConfig
	captureDeviceConfig.DeviceType = engine.Capture
	device, err := engine.InitDevice(ctx.Context, captureDeviceConfig, captureCallbacks)
	if err != nil {
		t.Fatal(err)
	}

	if device.Type() != engine.Capture {
		t.Errorf("wrong device type")
	}

	if device.PlaybackFormat() != engine.FormatS16 {
		t.Errorf("wrong format")
	}

	if device.PlaybackChannels() != 2 {
		t.Errorf("wrong number of channels")
	}

	if device.SampleRate() != 44100 {
		t.Errorf("wrong samplerate")
	}

	err = device.Start()
	if err != nil {
		t.Fatal(err)
	}

	if !device.IsStarted() {
		t.Fatalf("device not started")
	}

	time.Sleep(1 * time.Second)

	device.Uninit()

	onSendFrames := func(outputSamples, inputSamples []byte, framecount uint32) {
		samplesToRead := framecount * deviceConfig.Playback.Channels * sizeInBytes
		if samplesToRead > capturedSampleCount-playbackSampleCount {
			samplesToRead = capturedSampleCount - playbackSampleCount
		}

		copy(outputSamples, pCapturedSamples[playbackSampleCount:playbackSampleCount+samplesToRead])

		playbackSampleCount += samplesToRead
	}

	playbackCallbacks := engine.DeviceCallbacks{
		Data: onSendFrames,
	}

	playbackDeviceConfig := deviceConfig
	playbackDeviceConfig.DeviceType = engine.Playback
	device, err = engine.InitDevice(ctx.Context, playbackDeviceConfig, playbackCallbacks)
	if err != nil {
		t.Fatal(err)
	}

	if device.Type() != engine.Playback {
		t.Errorf("wrong device type")
	}

	err = device.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	device.Uninit()
}

func TestErrors(t *testing.T) {
	_, err := engine.InitContext([]engine.Backend{engine.Backend(99)}, engine.ContextConfig{}, nil)
	if err == nil {
		t.Fatalf("context init with invalid backend")
	}

	ctx, err := engine.InitContext(nil, engine.ContextConfig{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	onSendFrames := func(outputSamples, inputSamples []byte, framecount uint32) {
	}

	deviceConfig := engine.DefaultDeviceConfig(engine.Playback)
	deviceConfig.Playback.Format = engine.FormatType(99)
	deviceConfig.Playback.Channels = 99
	deviceConfig.SampleRate = 44100

	_, err = engine.InitDevice(ctx.Context, deviceConfig, engine.DeviceCallbacks{})
	if err == nil {
		t.Fatalf("device init with invalid config")
	}

	deviceConfig.Playback.Format = engine.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = 44100

	dev, err := engine.InitDevice(ctx.Context, deviceConfig, engine.DeviceCallbacks{
		Data: onSendFrames,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = dev.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = dev.Start()
	if err == nil {
		t.Fatalf("device start but already started")
	}

	time.Sleep(1 * time.Second)

	err = dev.Stop()
	if err != nil {
		t.Fatal(err)
	}

	err = dev.Stop()
	if err == nil {
		t.Fatalf("device stop but already stopped")
	}

	dev.Uninit()
}