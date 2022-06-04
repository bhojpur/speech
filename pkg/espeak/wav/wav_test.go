package wav

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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkInt32ToBytes(b *testing.B) {
	w := newWavHeader()
	for i := 0; i < b.N; i++ {
		w.littleEndianInt32ToBytes(40, 1<<16-1)
	}
}

func BenchmarkInt32ToBytesBinary(b *testing.B) {
	w := newWavHeader()
	for i := 0; i < b.N; i++ {
		w.littleEndianInt32ToBytesBinary(40, 1<<16-1)
	}
}

func TestWriter_WriteSamples(t *testing.T) {
	tmp, err := ioutil.TempDir("", "bhojpur-speech-wav-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	fh, err := os.Create(filepath.Join(tmp, "test.wav"))
	if err != nil {
		t.Fatal(err)
	}
	w := NewWriter(fh, 44100)
	in := []int16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	w.WriteSamples(in)
	fh.Close()

	fh, err = os.Open(filepath.Join(tmp, "test.wav"))
	if err != nil {
		t.Fatal(err)
	}
	got := make([]byte, 44, 44)
	io.ReadFull(fh, got)
	head := newWavHeader()
	for i := 0; i < 4; i++ {
		if head[i] != got[i] {
			t.Errorf("expected byte %d to be %v, instead got %v", i, head[i], got[i])
		}
	}
	// data size in bytes is len(in)*2
	if got[40] != byte(len(in)*2) {
		t.Errorf("expected byte 40 to be %v, instead got %v", len(in), got[40])
	}
}
