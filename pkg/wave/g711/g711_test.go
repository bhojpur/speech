package g711

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
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"testing/iotest"
)

var EncoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 12},
}

var DecoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00}, 4},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 24},
	{[]byte{0x01, 0x00, 0xdc, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 26},
}

var TranscoderTest = []struct {
	data     []byte
	expected int
}{
	{[]byte{}, 0},
	{[]byte{0x01, 0x00}, 2},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde}, 12},
	{[]byte{0x01, 0x00, 0x7c, 0x7f, 0xd1, 0xd0, 0xd3, 0xd2, 0xdd, 0xdc, 0xdf, 0xde, 0xd9}, 13},
}

// Test Encoding
func TestEncode(t *testing.T) {
	aenc, _ := NewAlawEncoder(ioutil.Discard, Lpcm)
	for _, tc := range EncoderTest {
		i, _ := aenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("Alaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	uenc, _ := NewUlawEncoder(ioutil.Discard, Lpcm)
	for _, tc := range EncoderTest {
		i, _ := uenc.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Encode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	utrans, _ := NewUlawEncoder(ioutil.Discard, Alaw)
	for _, tc := range TranscoderTest {
		i, _ := utrans.Write(tc.data)
		if i != tc.expected {
			t.Errorf("ulaw Transcode: expected: %d , actual: %d", tc.expected, i)
		}
	}
}

// Test Decoding
func TestDecode(t *testing.T) {
	b := new(bytes.Buffer)
	adec, _ := NewAlawDecoder(b)
	for _, tc := range DecoderTest {
		b.Write(tc.data)
		p := make([]byte, 16)
		var err error
		var i, n int
		for err == nil {
			n, err = adec.Read(p)
			i += n
		}
		if i != tc.expected {
			t.Errorf("Alaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	b.Reset()
	udec, _ := NewUlawDecoder(b)
	for _, tc := range DecoderTest {
		b.Write(tc.data)
		p := make([]byte, 16)
		var err error
		var i, n int
		for err == nil {
			n, err = udec.Read(p)
			i += n
		}
		if i != tc.expected {
			t.Errorf("ulaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
	b.Reset()
	// Edge Case
	udec, _ = NewUlawDecoder(iotest.TimeoutReader(b))
	for _, tc := range DecoderTest {
		b.Write(tc.data)
		p := make([]byte, 16)
		var err error
		var i, n int
		for err == nil || err.Error() == "timeout" {
			n, err = udec.Read(p)
			i += n
		}
		if i != tc.expected {
			t.Errorf("ulaw Decode: expected: %d , actual: %d", tc.expected, i)
		}
	}
}

// Benchmark Encoding data to Alaw
func BenchmarkAEncode(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)
	}
	b.SetBytes(int64(len(rawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewAlawEncoder(ioutil.Discard, Lpcm)
		if err != nil {
			b.Fatalf("Failed to create Writer: %s\n", err)
		}
		_, err = encoder.Write(rawData)
		if err != nil {
			b.Fatalf("Encoding failed: %s\n", err)
		}
	}
}

// Benchmark Encoding data to Ulaw
func BenchmarkUEncode(b *testing.B) {
	rawData, err := ioutil.ReadFile("testing/speech.raw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(rawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder, err := NewUlawEncoder(ioutil.Discard, Lpcm)
		if err != nil {
			b.Fatalf("Failed to create Writer: %s\n", err)

		}
		_, err = encoder.Write(rawData)
		if err != nil {
			b.Fatalf("Encoding failed: %s\n", err)

		}
	}
}

// Benchmark transcoding g711 data
func BenchmarkTranscode(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(alawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transcoder, err := NewAlawEncoder(ioutil.Discard, Ulaw)
		if err != nil {
			b.Fatalf("Failed to create Writer: %s\n", err)

		}
		_, err = transcoder.Write(alawData)
		if err != nil {
			b.Fatalf("Transcoding failed: %s\n", err)

		}
	}
}

// Benchmark Decoding Ulaw data
func BenchmarkUDecode(b *testing.B) {
	ulawData, err := ioutil.ReadFile("testing/speech.ulaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(ulawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewUlawDecoder(bytes.NewReader(ulawData))
		if err != nil {
			b.Fatalf("Failed to create Reader: %s\n", err)

		}
		_, err = io.Copy(ioutil.Discard, decoder)
		if err != nil {
			b.Fatalf("Decoding failed: %s\n", err)

		}
	}
}

// Benchmark Decoding Alaw data
func BenchmarkADecode(b *testing.B) {
	alawData, err := ioutil.ReadFile("testing/speech.alaw")
	if err != nil {
		b.Fatalf("Failed to read test data: %s\n", err)

	}
	b.SetBytes(int64(len(alawData)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		decoder, err := NewAlawDecoder(bytes.NewReader(alawData))
		if err != nil {
			b.Fatalf("Failed to create Reader: %s\n", err)

		}
		_, err = io.Copy(ioutil.Discard, decoder)
		if err != nil {
			b.Fatalf("Decoding failed: %s\n", err)

		}
	}
}
