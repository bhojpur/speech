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

	langdetection "github.com/bhojpur/speech/pkg/service/lang-detection"
	"github.com/bhojpur/speech/pkg/synthesis"
	voices "github.com/bhojpur/speech/pkg/voices"
)

func main() {
	fmt.Println("Bhojpur Speech translate utility")
	fmt.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	fmt.Println("All rights reserved.")

	if len(os.Args) < 2 {
		fmt.Println("\nNo input text provided")
		fmt.Printf("Usage: translate [TEXT]\n")
		os.Exit(1)
	} else {
		speech := synthesis.Speech{
			Folder:   "audios",
			Language: voices.English,
			Volume:   0,
			Speed:    1}

		detectionService := langdetection.NewLingualDetectionService(langdetection.DefaultLanguages)
		langDetected, _ := detectionService.Detect(os.Args[1])
		lang := string(*langDetected)
		fmt.Printf("Detected Language: %s\n", lang)

		speech.Speak(os.Args[1])
	}
}
