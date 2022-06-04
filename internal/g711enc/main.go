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

// It encodes 16bit 8kHz LPCM data to 8bit G711 PCM.
// It works with wav or raw files as input.

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/speech/pkg/wave/g711"
)

func main() {
	log.Println("Bhojpur Speech G.711 Encoder utility")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	if len(os.Args) < 3 || os.Args[1] == "help" || os.Args[1] == "--help" || (os.Args[1] != "ulaw" && os.Args[1] != "alaw") {
		fmt.Printf("%s Encodes 16bit 8kHz LPCM data to 8bit G711 PCM\n", os.Args[0])
		fmt.Println("The program takes as input a list or wav or raw files, encodes them")
		fmt.Println("to G711 PCM and saves them with the proper extension.")
		fmt.Printf("\nUsage: %s [encoding format] [files]\n", os.Args[0])
		fmt.Println("encoding format can be either alaw or ulaw")
		os.Exit(1)
	}
	var exitCode int
	format := os.Args[1]
	for _, file := range os.Args[2:] {
		err := encodeG711(file, format)
		if err != nil {
			log.Println(err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func encodeG711(file, format string) error {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	extension := strings.ToLower(filepath.Ext(file))
	if extension != ".wav" && extension != ".raw" && extension != ".sln" {
		err = fmt.Errorf("Unrecognised format for input file: %s", file)
		return err
	}
	outName := strings.TrimSuffix(file, filepath.Ext(file)) + "." + format
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	encoder := new(g711.Encoder)
	if format == "alaw" {
		encoder, err = g711.NewAlawEncoder(outFile, g711.Lpcm)
		if err != nil {
			return err
		}
	} else {
		encoder, err = g711.NewUlawEncoder(outFile, g711.Lpcm)
		if err != nil {
			return err
		}
	}
	if extension == ".wav" {
		_, err = encoder.Write(input[44:]) // Skip WAV header
		return err
	}
	_, err = encoder.Write(input)
	return err
}
