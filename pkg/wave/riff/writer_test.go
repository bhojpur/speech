package riff

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
	"testing"
)

func TestWriteRIFF(t *testing.T) {
	testFiles := []testFile{
		testFile{
			"a.wav",
			3,
			243800,
			"WAVE"},
		testFile{
			"1_webp_a.webp",
			3,
			23396,
			"WEBP"}}

	for _, testFile := range testFiles {
		file, err := fixtureFile(testFile.Name)

		if err != nil {
			t.Fatalf("Failed to open fixture file")
		}

		reader := NewReader(file)
		riff, err := reader.Read()

		if err != nil {
			t.Fatal(err)
		}

		outfile, err := ioutil.TempFile("/tmp", "outfile")

		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			outfile.Close()
			os.Remove(outfile.Name())
		}()

		writer := NewWriter(outfile, riff.FileType, riff.FileSize)

		for _, chunk := range riff.Chunks {
			writer.WriteChunk(chunk.ChunkID, chunk.ChunkSize, func(w io.Writer) {
				buf, err := ioutil.ReadAll(chunk)

				if err != nil {
					t.Fatal(err)
				}

				w.Write(buf)
			})
		}

		outfile.Close()

		// reopen to check file content
		file, err = os.Open(outfile.Name())

		if err != nil {
			t.Fatal(err)
		}

		reader = NewReader(file)
		riff, err = reader.Read()

		if err != nil {
			t.Fatal(err)
		}

		for _, chunk := range riff.Chunks {
			t.Logf("Chunk ID: %s", string(chunk.ChunkID[:]))
		}

		if len(riff.Chunks) != testFile.ChunkSize {
			t.Fatalf("Invalid length of chunks")
		}

		if riff.FileSize != testFile.FileSize {
			t.Fatalf("File size is invalid: %d", riff.FileSize)
		}

		if string(riff.FileType[:]) != testFile.FileType {
			t.Fatalf("File type is invalid: %s", riff.FileType)
		}
	}
}
