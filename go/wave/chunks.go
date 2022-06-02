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
	"bytes"
	"io"
)

const (
	maxFileSize             = 2 << 31
	riffChunkSize           = 12
	listChunkOffset         = 36
	riffChunkSizeBaseOffset = 36 // RIFFChunk(12byte) + fmtChunk(24byte) = 36byte
	fmtChunkSize            = 16
)

var (
	riffChunkToken = "RIFF"
	waveFormatType = "WAVE"
	fmtChunkToken  = "fmt "
	listChunkToken = "LIST"
	dataChunkToken = "data"
)

// 12byte
type RiffChunk struct {
	ID         []byte // 'RIFF'
	Size       uint32 // 36bytes + data_chunk_size or whole_file_size - 'RIFF'+ChunkSize (8byte)
	FormatType []byte // 'WAVE'
}

// 8 + 16 = 24byte
type FmtChunk struct {
	ID   []byte // 'fmt '
	Size uint32 // 16
	Data *WavFmtChunkData
}

// 16byte
type WavFmtChunkData struct {
	WaveFormatType uint16 // PCM 1
	Channel        uint16 // monoral or streo
	SamplesPerSec  uint32 // 44100
	BytesPerSec    uint32 // byte
	BlockSize      uint16 // *
	BitsPerSamples uint16 //
}

// data
type DataReader interface {
	io.Reader
	io.ReaderAt
}

type DataReaderChunk struct {
	ID   []byte     // 'data'
	Size uint32     // * channel
	Data DataReader //
}

type DataWriterChunk struct {
	ID   []byte
	Size uint32
	Data *bytes.Buffer
}