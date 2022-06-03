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
	"flag"
	"fmt"
	"testing"

	engine "github.com/bhojpur/speech/pkg/miniaudio"
)

var testWithHardware = flag.Bool("engine.hardware", false, "run tests with expecting hardware")

func TestContextLifecycle(t *testing.T) {
	config := engine.ContextConfig{ThreadPriority: engine.ThreadPriorityNormal}

	ctx, err := engine.InitContext(nil, config, nil)
	assertNil(t, err, "No error expected initializing context")
	assertNotNil(t, ctx, "Context instance expected")
	assertNotEqual(t, engine.Context(0), ctx.Context, "Context value expected")

	err = ctx.Uninit()
	assertNil(t, err, "No error expected uninitializing")

	ctx.Free()
	assertEqual(t, engine.Context(0), ctx.Context, "Expected context value to be reset")
}

func TestContextDeviceEnumeration(t *testing.T) {
	if *testWithHardware {
		t.Log("Running test expecting devices\n")
	}

	ctx, err := engine.InitContext(nil, engine.ContextConfig{}, nil)
	assertNil(t, err, "No error expected initializing context")
	defer func() {
		err := ctx.Uninit()
		assertNil(t, err, "No error expected uninitializing")
		ctx.Free()
	}()

	playbackDevices, err := ctx.Devices(engine.Playback)
	assertNil(t, err, "No error expected querying playback devices")
	if *testWithHardware {
		assertTrue(t, len(playbackDevices) > 0, "No playback devices found")
	}

	captureDevices, err := ctx.Devices(engine.Capture)
	assertNil(t, err, "No error expected querying capture devices")
	if *testWithHardware {
		assertTrue(t, len(captureDevices) > 0, "No capture devices found")
	}
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}
	t.Fatal(message)
}

func assertNil(t *testing.T, v interface{}, message string) {
	if v == nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expected nil, got %#v", v)
	}
	t.Fatal(message)
}

func assertNotNil(t *testing.T, v interface{}, message string) {
	if v != nil {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expected value not to be nil")
	}
	t.Fatal(message)
}

func assertTrue(t *testing.T, v bool, message string) {
	if v {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("should be true")
	}
	t.Fatal(message)
}
