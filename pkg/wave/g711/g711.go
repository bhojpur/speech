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

// It implements encoding and decoding of G711 PCM sound data.
// G.711 is an ITU-T standard for audio companding.

import (
	"errors"
	"io"
)

const (
	// Input and output formats
	Alaw = iota // Alaw G711 encoded PCM data
	Ulaw        // Ulaw G711  encoded PCM data
	Lpcm        // Lpcm 16bit signed linear data
)

// Decoder reads G711 PCM data and decodes it to 16bit 8000Hz LPCM
type Decoder struct {
	decode func([]byte) []byte // decoding function
	source io.Reader           // source data
}

// Encoder encodes 16bit 8000Hz LPCM data to G711 PCM or
// directly transcodes between A-law and u-law
type Encoder struct {
	input       int                 // input format
	encode      func([]byte) []byte // encoding function
	transcode   func([]byte) []byte // transcoding function
	destination io.Writer           // output data
}

// NewAlawDecoder returns a pointer to a Decoder that implements an io.Reader.
// It takes as input the source data Reader.
func NewAlawDecoder(reader io.Reader) (*Decoder, error) {
	if reader == nil {
		return nil, errors.New("io.Reader is nil")
	}
	r := Decoder{
		decode: DecodeAlaw,
		source: reader,
	}
	return &r, nil
}

// NewUlawDecoder returns a pointer to a Decoder that implements an io.Reader.
// It takes as input the source data Reader.
func NewUlawDecoder(reader io.Reader) (*Decoder, error) {
	if reader == nil {
		return nil, errors.New("io.Reader is nil")
	}
	r := Decoder{
		decode: DecodeUlaw,
		source: reader,
	}
	return &r, nil
}

// NewAlawEncoder returns a pointer to an Encoder that implements an io.Writer.
// It takes as input the destination data Writer and the input encoding format.
func NewAlawEncoder(writer io.Writer, input int) (*Encoder, error) {
	if writer == nil {
		return nil, errors.New("io.Writer is nil")
	}
	if input != Ulaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	w := Encoder{
		input:       input,
		encode:      EncodeAlaw,
		transcode:   Ulaw2Alaw,
		destination: writer,
	}
	return &w, nil
}

// NewUlawEncoder returns a pointer to an Encoder that implements an io.Writer.
// It takes as input the destination data Writer and the input encoding format.
func NewUlawEncoder(writer io.Writer, input int) (*Encoder, error) {
	if writer == nil {
		return nil, errors.New("io.Writer is nil")
	}
	if input != Alaw && input != Lpcm {
		return nil, errors.New("Invalid input format")
	}
	w := Encoder{
		input:       input,
		encode:      EncodeUlaw,
		transcode:   Alaw2Ulaw,
		destination: writer,
	}
	return &w, nil
}

// Reset discards the Decoder state. This permits reusing a Decoder rather than allocating a new one.
func (r *Decoder) Reset(reader io.Reader) error {
	if reader == nil {
		return errors.New("io.Reader is nil")
	}
	r.source = reader
	return nil
}

// Reset discards the Encoder state. This permits reusing an Encoder rather than allocating a new one.
func (w *Encoder) Reset(writer io.Writer) error {
	if writer == nil {
		return errors.New("io.Writer is nil")
	}
	w.destination = writer
	return nil
}

// Read decodes G711 data. Reads up to len(p) bytes into p, returns the number
// of bytes read and any error encountered.
func (r *Decoder) Read(p []byte) (i int, err error) {
	if len(p) == 0 {
		return
	}
	b := make([]byte, len(p)/2)
	i, err = r.source.Read(b)
	copy(p, r.decode(b))
	i *= 2 // Report back the correct number of bytes
	return
}

// Write encodes G711 Data. Writes len(p) bytes from p to the underlying data stream,
// returns the number of bytes written from p (0 <= n <= len(p)) and any error encountered
// that caused the write to stop early.
func (w *Encoder) Write(p []byte) (i int, err error) {
	if len(p) == 0 {
		return
	}
	if w.input == Lpcm { // Encode LPCM data to G711
		i, err = w.destination.Write(w.encode(p))
		if err == nil && len(p)%2 != 0 {
			err = errors.New("Odd number of LPCM bytes, incomplete frame")
		}
		i *= 2 // Report back the correct number of bytes written from p
	} else { // Trans-code
		i, err = w.destination.Write(w.transcode(p))
	}
	return
}
