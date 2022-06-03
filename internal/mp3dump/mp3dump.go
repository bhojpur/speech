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

// It decodes an MP3 file and writes the raw PCM data to a file

import (
	"fmt"
	"github.com/bhojpur/speech/pkg/mpg123"
	"os"
)

func main() {
	// check command-line arguments
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: mp3dump <infile.mp3> <outfile.raw>")
		return
	}

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	if err != nil {
		panic("could not initialize mpg123")
	}

	// open a file with decoder
	err = decoder.Open(os.Args[1])
	if err != nil {
		panic("error opening mp3 file")
	}
	defer decoder.Close()

	// get audio format information
	rate, chans, _ := decoder.GetFormat()
	fmt.Fprintln(os.Stderr, "Encoding: Signed 16bit")
	fmt.Fprintln(os.Stderr, "Sample Rate:", rate)
	fmt.Fprintln(os.Stderr, "Channels:", chans)

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, chans, mpg123.ENC_SIGNED_16)

	// open output file
	o, err := os.Create(os.Args[2])
	if err != nil {
		panic("error opening output file")
	}
	defer o.Close()

	// decode mp3 file and dump output
	buf := make([]byte, 2048*16)
	for {
		len, err := decoder.Read(buf)
		o.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	o.Close()
	decoder.Delete()
}