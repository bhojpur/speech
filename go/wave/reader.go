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
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type WaveReader interface {
	io.Reader
	io.Seeker
	io.ReaderAt
}

type Reader struct {
	input WaveReader

	size int64

	RiffChunk *RiffChunk
	FmtChunk  *FmtChunk
	DataChunk *DataReaderChunk

	originOfAudioData int64
	NumSamples        uint32
	ReadSampleNum     uint32
	SampleTime        int

	// LIST chunk
	extChunkSize int64
}

func NewReader(fileName string) (*Reader, error) {
	// check file size
	fi, err := os.Stat(fileName)
	if err != nil {
		return &Reader{}, err
	}
	if fi.Size() > maxFileSize {
		return &Reader{}, fmt.Errorf("file is too large: %d bytes", fi.Size())
	}

	// open file
	f, err := os.Open(fileName)
	if err != nil {
		return &Reader{}, err
	}
	defer f.Close()

	waveData, err := ioutil.ReadAll(f)
	if err != nil {
		return &Reader{}, err
	}

	reader := new(Reader)
	reader.size = fi.Size()
	reader.input = bytes.NewReader(waveData)

	if err := reader.parseRiffChunk(); err != nil {
		panic(err)
	}
	if err := reader.parseFmtChunk(); err != nil {
		panic(err)
	}
	if err := reader.parseListChunk(); err != nil {
		panic(err)
	}
	if err := reader.parseDataChunk(); err != nil {
		panic(err)
	}

	reader.NumSamples = reader.DataChunk.Size / uint32(reader.FmtChunk.Data.BlockSize)
	reader.SampleTime = int(reader.NumSamples / reader.FmtChunk.Data.SamplesPerSec)

	return reader, nil
}

type csize struct {
	ChunkSize uint32
}

func (rd *Reader) parseRiffChunk() error {
	// RIFF
	chunkId := make([]byte, 4)
	if err := binary.Read(rd.input, binary.BigEndian, chunkId); err != nil {
		return err
	}
	if string(chunkId[:]) != riffChunkToken {
		return fmt.Errorf("file is not RIFF: %s", rd.RiffChunk.ID)
	}

	// RIFF
	chunkSize := &csize{}
	if err := binary.Read(rd.input, binary.LittleEndian, chunkSize); err != nil {
		return err
	}
	if chunkSize.ChunkSize+8 != uint32(rd.size) {
		//		fmt.Println("======================")
		//		fmt.Println("riff chunk size ", rd.riffChunk.ChunkSize)
		//		fmt.Println("file size ", rd.size)
		//		fmt.Println("======================")
		return fmt.Errorf("riff_chunk_size must be whole file size - 8bytes, expected(%d), actual(%d)", chunkSize.ChunkSize+8, rd.size)
	}

	// RIFF 'WAVE'
	format := make([]byte, 4)
	if err := binary.Read(rd.input, binary.BigEndian, format); err != nil {
		return err
	}
	if string(format[:]) != waveFormatType {
		return fmt.Errorf("file is not WAVE: %s", rd.RiffChunk.FormatType)
	}

	riffChunk := RiffChunk{
		ID:         chunkId,
		Size:       chunkSize.ChunkSize,
		FormatType: format,
	}

	rd.RiffChunk = &riffChunk

	return nil
}

func (rd *Reader) parseFmtChunk() error {
	rd.input.Seek(riffChunkSize, os.SEEK_SET)

	// 'fmt '
	chunkId := make([]byte, 4)
	err := binary.Read(rd.input, binary.BigEndian, chunkId)
	if err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	}
	if string(chunkId[:]) != fmtChunkToken {
		return fmt.Errorf("fmt chunk id must be \"%s\" but value is %s", fmtChunkToken, chunkId)
	}

	// fmt_chunk_size 16bit
	chunkSize := &csize{}
	err = binary.Read(rd.input, binary.LittleEndian, chunkSize)
	if err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	}
	if chunkSize.ChunkSize != fmtChunkSize {
		return fmt.Errorf("fmt chunk size must be %d but value is %d", fmtChunkSize, chunkSize.ChunkSize)
	}

	// fmt_chunk_data
	var fmtChunkData WavFmtChunkData
	if err = binary.Read(rd.input, binary.LittleEndian, &fmtChunkData); err != nil {
		return err
	}

	fmtChunk := FmtChunk{
		ID:   chunkId,
		Size: chunkSize.ChunkSize,
		Data: &fmtChunkData,
	}

	rd.FmtChunk = &fmtChunk

	return nil
}

func (rd *Reader) parseListChunk() error {
	rd.input.Seek(listChunkOffset, os.SEEK_SET)

	// 'LIST'
	chunkID := make([]byte, 4)
	if err := binary.Read(rd.input, binary.BigEndian, chunkID); err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	} else if string(chunkID[:]) != listChunkToken {
		// LIST
		return nil
	}

	// 'LIST' 1byte
	chunkSize := make([]byte, 1)
	if err := binary.Read(rd.input, binary.LittleEndian, chunkSize); err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	}

	// header
	// rd.extChunkSize += int64(chunkSize[0]) + 4 + 4
	rd.extChunkSize = int64(chunkSize[0]) + 4 + 4

	return nil
}

// header riffChunkSizeOffset
func (rd *Reader) getRiffChunkSizeOffset() int64 {
	return riffChunkSizeBaseOffset + rd.extChunkSize
}

func (rd *Reader) parseDataChunk() error {
	originOfDataChunk, _ := rd.input.Seek(rd.getRiffChunkSizeOffset(), os.SEEK_SET)

	// 'data'
	chunkId := make([]byte, 4)
	err := binary.Read(rd.input, binary.BigEndian, chunkId)
	if err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	}
	if string(chunkId[:]) != dataChunkToken {
		return fmt.Errorf("data chunk id must be \"%s\" but value is %s", dataChunkToken, chunkId)
	}

	// data_chunk_size
	chunkSize := &csize{}
	err = binary.Read(rd.input, binary.LittleEndian, chunkSize)
	if err == io.EOF {
		return fmt.Errorf("unexpected file end")
	} else if err != nil {
		return err
	}

	// dataChunk ID(4byte) + chunkSize(4byte)
	rd.originOfAudioData = originOfDataChunk + 8
	audioData := io.NewSectionReader(rd.input, rd.originOfAudioData, int64(chunkSize.ChunkSize))

	dataChunk := DataReaderChunk{
		ID:   chunkId,
		Size: chunkSize.ChunkSize,
		Data: audioData,
	}

	rd.DataChunk = &dataChunk

	return nil
}

//
func (rd *Reader) Read(p []byte) (int, error) {
	n, err := rd.DataChunk.Data.Read(p)
	return n, err
}

func (rd *Reader) ReadRawSample() ([]byte, error) {
	size := rd.FmtChunk.Data.BlockSize
	sample := make([]byte, size)
	_, err := rd.Read(sample)
	if err == nil {
		rd.ReadSampleNum += 1
	}
	return sample, err
}

func (rd *Reader) ReadSample() ([]float64, error) {
	raw, err := rd.ReadRawSample()
	channel := int(rd.FmtChunk.Data.Channel)
	ret := make([]float64, channel)
	length := len(raw) / channel // 1 byte

	if err != nil {
		return ret, err
	}

	for i := 0; i < channel; i++ {
		tmp := bytesToInt(raw[length*i : length*(i+1)])
		switch rd.FmtChunk.Data.BitsPerSamples {
		case 8:
			ret[i] = float64(tmp-128) / 128.0
		case 16:
			ret[i] = float64(tmp) / 32768.0
		}
		if err != nil && err != io.EOF {
			return ret, err
		}
	}
	return ret, nil
}

func (rd *Reader) ReadSampleInt() ([]int, error) {
	raw, err := rd.ReadRawSample()
	channels := int(rd.FmtChunk.Data.Channel)
	ret := make([]int, channels)
	length := len(raw) / channels // 1 byte

	if err != nil {
		return ret, err
	}

	for i := 0; i < channels; i++ {
		ret[i] = bytesToInt(raw[length*i : length*(i+1)])
		if err != nil && err != io.EOF {
			return ret, err
		}
	}
	return ret, nil
}

func bytesToInt(b []byte) int {
	var ret int
	switch len(b) {
	case 1:
		// 0 ~ 128 ~ 255
		ret = int(b[0])
	case 2:
		// -32768 ~ 0 ~ 32767
		ret = int(b[0]) + int(b[1])<<8
	//	fmt.Printf("%08b %08b ", b[1], b[0])
	//	fmt.Printf("%016b => %d\n", ret, ret)
	case 3:
		// HiReso / DVDAudio
		ret = int(b[0]) + int(b[1])<<8 + int(b[2])<<16
	default:
		ret = 0
	}
	return ret
}