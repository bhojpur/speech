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
	"fmt"
	"os"

	"github.com/bhojpur/speech/pkg/espeak"
)

func main() {
	fmt.Println("Bhojpur Speech text-to-speech utility")
	fmt.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	fmt.Printf("All rights reserved.\n")

	if len(os.Args) < 2 {
		fmt.Println("No input audio text provided.")
		os.Exit(1)
	}

	// need to call terminate so eSpeak can clean itself out
	defer espeak.Terminate()
	params := espeak.NewParameters().WithDir(".")
	var written uint64
	written, _ = espeak.TextToSpeech(
		os.Args[1],        // Text to speak in English language
		nil,               // voice to use, nil == DefaultVoice (en-us male)
		"audios/test.wav", // if "" or "play", it plays to default audio out
		params,            // Parameters for voice modulation, nil == DefaultParameters
	)
	fmt.Printf("Bhojpur Speech bytes written to audios/test.wav:\t%d\n", written)

	// get a random Hindi voice
	v, _ := espeak.VoiceFromSpec(&espeak.Voice{Languages: "hi"})
	written, _ = espeak.TextToSpeech(os.Args[1], v, "audios/test_hi.wav", params)
	fmt.Printf("Bhojpur Speech bytes written to audios/test_hi.wav:\t%d\n", written)
}
