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
	"encoding/binary"
	"io"
)

type wavHeader [44]byte

func newWavHeader() *wavHeader {
	var w = &wavHeader{
		'R', 'I', 'F', 'F', // 0-3: Marks file as riff file.
		0, 0, 0, 0, // 4-7: Size of overall file - 8 bytes as int32.
		'W', 'A', 'V', 'E', // 8-11: File Type Header. Always equals "WAVE"
		'f', 'm', 't', ' ', //  12-15: Format chunk marker.
		0x10, 0, 0, 0, // 16-19: Length of format data
		1, 0, // 20-21: Type of format (1 is PCM) - 2 byte integer
		1, 0, // 22-23: Number of channels - 2 byte integer
		0xa, 0x0c, 0x004, 0x0004, // 24-27: Sample Rate - 32 byte integer. Common values are 44100 (CD), 48000 (DAT). Sample Rate = Number of Samples per second, or Hertz.
		0, 0, 0, 0, // 28-31: (Sample Rate * BitsPerSample * Channels) / 8.
		2, 0, // 32-33:  (BitsPerSample * Channels) / 8.1 - 8 bit mono2 - 8 bit stereo/16 bit mono4 - 16 bit stereo
		0x10, 0, // 34-35: Bits per sample
		'd', 'a', 't', 'a', // 36-39: "data" chunk header. Marks the beginning of the data section.
		0, 0, 0, 0, // 40-43: Size of the data section.
	}
	return w
}

func (w *wavHeader) littleEndianInt32ToBytes(pos int, value int32) {
	(*w)[pos] = byte(value & 0xff)
	(*w)[pos+1] = byte((value >> 8) & 0xff)
	(*w)[pos+2] = byte((value >> 16) & 0xff)
	(*w)[pos+3] = byte((value >> 24) & 0xff)
}
func (w *wavHeader) littleEndianInt32ToBytesBinary(pos int, value int32) {
	binary.LittleEndian.PutUint32(w[pos:(pos+4)], uint32(value))
}

func (w *wavHeader) writeSize(size int32) {
	w.littleEndianInt32ToBytes(4, size)
}

func (w *wavHeader) writeSampleRate(rate int32) {
	w.littleEndianInt32ToBytes(24, rate)
}
func (w *wavHeader) writeByteRate(rate int32) {
	w.littleEndianInt32ToBytes(28, rate)
}

func (w *wavHeader) writeDataBytes(bytes int32) {
	w.littleEndianInt32ToBytes(40, bytes)
}

// Writer has utilities for writing to .wav files.
type Writer struct {
	out          io.Writer
	err          error
	bytesWritten uint64
	sampleRate   int32
}

// NewWriter returns a *Writer.
func NewWriter(w io.Writer, sampleRate int32) *Writer {
	return &Writer{out: w, sampleRate: sampleRate}
}

// Write implements the io.Writer interface.
func (w *Writer) Write(data []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	n, err := w.out.Write(data)
	w.bytesWritten += uint64(n)
	w.err = err
	return n, err
}

// WriteSamples writes the .wav header and an []int16 to the file.
func (w *Writer) WriteSamples(data []int16) (uint64, error) {
	h := newWavHeader()
	h.writeSampleRate(w.sampleRate)
	h.writeByteRate(w.sampleRate * 2)
	h.writeDataBytes(int32(len(data) * 2))
	h.writeSize(int32(len(data)*2+binary.Size(h)) - 8)

	binary.Write(w, binary.LittleEndian, h)
	w.bytesWritten += uint64(binary.Size(h))
	binary.Write(w, binary.LittleEndian, data)

	return w.bytesWritten, w.err
}
