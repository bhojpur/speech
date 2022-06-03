package miniaudio

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

// A mini audio library (miniaudio cgo bindings).

/*
#cgo CFLAGS: -std=gnu99 -Wno-unused-result
#cgo ma_debug CFLAGS: -DMA_DEBUG_OUTPUT=1

#cgo linux,!android LDFLAGS: -ldl -lpthread -lm
#cgo openbsd LDFLAGS: -ldl -lpthread -lm
#cgo netbsd LDFLAGS: -ldl -lpthread -lm
#cgo freebsd LDFLAGS: -ldl -lpthread -lm
#cgo android LDFLAGS: -lm

#cgo !noasm,!arm,!arm64 CFLAGS: -msse2
#cgo !noasm,arm,arm64 CFLAGS: -mfpu=neon -mfloat-abi=hard
#cgo noasm CFLAGS: -DMA_NO_SSE2 -DMA_NO_AVX2 -DMA_NO_AVX512 -DMA_NO_NEON

#include "minialgo.h"
*/
import "C"

// SampleSizeInBytes retrieves the size of a sample in bytes for the given format.
func SampleSizeInBytes(format FormatType) int {
	cformat := (C.ma_format)(format)
	ret := C.ma_get_bytes_per_sample(cformat)
	return int(ret)
}

const (
	rawContextConfigSize = C.sizeof_ma_context_config
	rawDeviceInfoSize    = C.sizeof_ma_device_info
	rawDeviceConfigSize  = C.sizeof_ma_device_config
)