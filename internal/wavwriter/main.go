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

import (
	"math"
	"os"

	"github.com/bhojpur/speech/go/wave"
)

func main() {
	f, err := os.Create("./internal/test.wav")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	param := wave.WriterParam{
		Out:           f,
		Channel:       1,
		SampleRate:    44100,
		BitsPerSample: 16,
	}

	w, err := wave.NewWriter(param)

	amplitude := 0.1
	hz := 440.0
	length := param.SampleRate * 1

	for i := 0; i < length; i++ {
		_data := amplitude * math.Sin(2.0*math.Pi*hz*float64(i)/float64(param.SampleRate))
		_data = (_data + 1.0) / 2.0 * 65536.0
		if _data > 65535.0 {
			_data = 65535.0
		} else if _data < 0.0 {
			_data = 0.0
		}
		data := uint16(_data+0.5) - 32768 //
		var td []int16
		td = []int16{int16(data)}
		_, err = w.WriteSample16(td)
		if err != nil {
			panic(err)
		}
	}

	defer w.Close()
	if err != nil {
		panic(err)
	}
}
