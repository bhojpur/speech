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

// It decodes 8bit G711 PCM data to 16 Bit signed LPCM raw data

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/speech/pkg/wave/g711"
)

func main() {
	log.Println("Bhojpur Speech G.711 Decoder utility")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	if len(os.Args) == 1 || os.Args[1] == "help" || os.Args[1] == "--help" {
		fmt.Printf("%s Decodes 8bit G711 PCM data to raw 16 Bit signed LPCM\n", os.Args[0])
		fmt.Println("The program takes as input a list A-law or u-law encoded files")
		fmt.Println("decodes them to LPCM and saves the files with a \".raw\" extension.")
		fmt.Printf("\nUsage: %s [files]\n", os.Args[0])
		os.Exit(1)
	}
	var exitCode int
	for _, file := range os.Args[1:] {
		err := decodeG711(file)
		if err != nil {
			log.Println(err)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func decodeG711(file string) error {
	input, err := os.Open(file)
	if err != nil {
		return err
	}
	defer input.Close()

	extension := strings.ToLower(filepath.Ext(file))
	decoder := new(g711.Decoder)
	if extension == ".alaw" || extension == ".al" {
		decoder, err = g711.NewAlawDecoder(input)
		if err != nil {
			return err
		}
	} else if extension == ".ulaw" || extension == ".ul" {
		decoder, err = g711.NewUlawDecoder(input)
		if err != nil {
			return err
		}
	} else {
		err = fmt.Errorf("Unrecognised format for file: %s", file)
		return err
	}
	outName := strings.TrimSuffix(file, filepath.Ext(file)) + ".raw"
	outFile, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, decoder)
	return err
}
