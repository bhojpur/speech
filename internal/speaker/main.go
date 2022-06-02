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

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bhojpur/speech/go/coqui"
	"github.com/cryptix/wav"
)

var model = flag.String("model", "", "Path to the model (protocol buffer binary file)")
var audio = flag.String("audio", "", "Path to the audio file to run (WAV format)")
var scorer = flag.String("scorer", "", "Path to the external scorer")
var version = flag.Bool("version", false, "Print version and exit")
var extended = flag.Bool("extended", false, "Use extended metadata")
var maxResults = flag.Uint("max-results", 5, "Maximum number of results when -extended is true")
var printSampleRate = flag.Bool("sample-rate", false, "Print model sample rate and exit")

func metadataToStrings(m *coqui.Metadata) []string {
	results := make([]string, 0, m.NumTranscripts())
	for _, tr := range m.Transcripts() {
		var res string
		for _, tok := range tr.Tokens() {
			res += tok.Text()
		}
		res += fmt.Sprintf(" [%0.3f]", tr.Confidence())
		results = append(results, res)
	}
	return results
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *version {
		fmt.Println(coqui.Version())
		return
	}

	if *model == "" || *audio == "" {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	// Initialize Coqui
	m, err := coqui.New(*model)
	if err != nil {
		log.Fatal("Failed initializing model: ", err)
	}
	defer m.Close()

	if *printSampleRate {
		fmt.Println(m.SampleRate())
		return
	}

	if *scorer != "" {
		if err := m.EnableExternalScorer(*scorer); err != nil {
			log.Fatal("Failed enabling external scorer: ", err)
		}
	}

	// Stat audio
	i, err := os.Stat(*audio)
	if err != nil {
		log.Fatal(fmt.Errorf("stating %s failed: %w", *audio, err))
	}

	// Open audio
	f, err := os.Open(*audio)
	if err != nil {
		log.Fatal(fmt.Errorf("opening %s failed: %w", *audio, err))
	}

	// Create reader
	r, err := wav.NewReader(f, i.Size())
	if err != nil {
		log.Fatal(fmt.Errorf("creating new reader failed: %w", err))
	}

	// Read
	var d []int16
	for {
		// Read sample
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(fmt.Errorf("reading sample failed: %w", err))
		}

		// Append
		d = append(d, int16(s))
	}

	// Speech to text
	var results []string
	if *extended {
		metadata, err := m.SpeechToTextWithMetadata(d, *maxResults)
		if err != nil {
			log.Fatal("Failed converting speech to text: ", err)
		}
		defer metadata.Close()
		results = metadataToStrings(metadata)
	} else {
		res, err := m.SpeechToText(d)
		if err != nil {
			log.Fatal("Failed converting speech to text: ", err)
		}
		results = []string{res}
	}
	for _, res := range results {
		fmt.Println("Text:", res)
	}
}