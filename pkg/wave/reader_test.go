package wave

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
	"io"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestRead(t *testing.T) {
	blockAlign := 4

	file, err := fixtureFile("a.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatPCM, int(fmt.AudioFormat))
	assert.Equal(t, 2, int(fmt.NumChannels))
	assert.Equal(t, 44100, int(fmt.SampleRate))
	assert.Equal(t, 44100*4, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 16, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1.381496598s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	sample := samples[0]

	assert.Equal(t, 318, reader.IntValue(sample, 0))
	assert.Equal(t, 289, reader.IntValue(sample, 1))
	assert.Assert(t, math.Abs(reader.FloatValue(sample, 0)-0.009705) <= 0.0001)
	assert.Assert(t, math.Abs(reader.FloatValue(sample, 1)-0.008820) <= 0.0001)

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}

func TestReadMulaw(t *testing.T) {
	blockAlign := 1

	file, err := fixtureFile("mulaw.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatMULaw, int(fmt.AudioFormat))
	assert.Equal(t, 1, int(fmt.NumChannels))
	assert.Equal(t, 8000, int(fmt.SampleRate))
	assert.Equal(t, 8000, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 8, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "4.59125s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}

func TestReadAlaw(t *testing.T) {
	blockAlign := 1

	file, err := fixtureFile("alaw.wav")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	fmt, err := reader.Format()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, AudioFormatALaw, int(fmt.AudioFormat))
	assert.Equal(t, 1, int(fmt.NumChannels))
	assert.Equal(t, 8000, int(fmt.SampleRate))
	assert.Equal(t, 8000, int(fmt.ByteRate))
	assert.Equal(t, blockAlign, int(fmt.BlockAlign))
	assert.Equal(t, 8, int(fmt.BitsPerSample))

	duration, err := reader.Duration()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "4.59125s", duration.String())

	samples, err := reader.ReadSamples(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(samples))

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(bytes), int(reader.WavData.Size)-(1*blockAlign))

	t.Logf("Data size: %d", len(bytes))
}

func BenchmarkReadSamples(b *testing.B) {
	n := []uint32{1, 10, 100, 1000, 2000, 3000, 5000, 8000, 10000, 20000, 40000}

	var t int

	for _, numSamples := range n {
		b.Run(fmt.Sprintf("%d", numSamples), func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				file, _ := os.Open("./files/a.wav")
				reader := NewReader(file)

				for {
					samples, err := reader.ReadSamples(numSamples)
					if err == io.EOF {
						break
					}
					for _, sample := range samples {
						t += reader.IntValue(sample, 0)
						t += reader.IntValue(sample, 1)
					}
				}
			}
		})
	}
}
