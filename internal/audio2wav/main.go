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
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bhojpur/speech/go/portaudio"
	"github.com/bhojpur/speech/go/wave"
)

func main() {

	// Check arguments and resolve the .wav file name
	checkArgs()
	audioFileName := audioFileName()

	waveFile, err := os.Create(audioFileName)
	if err != nil {
		panic(err)
	}

	// Processing settings
	inChannels := 1
	outChannels := 0
	sampleRate := 44100
	framesPerBuffer := make([]byte, 64)

	// Init PortAudio
	if err := portaudio.Initialize(); err != nil {
		panic(err)
	}

	stream, err := portaudio.OpenDefaultStream(inChannels, outChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
	if err != nil {
		panic(err)
	}

	// Create the Wave writer
	ww, err := waveWriter(waveFile, inChannels, sampleRate)
	if err != nil {
		panic(err)
	}

	// Start recording
	fmt.Println("Recording (ESC to abort, 's' to save).")
	go startRecording(stream, ww, framesPerBuffer)

	// Wait for a signal to either abort or stop the recording
	awaitStopRecording(stream, ww)

	// Analyze the Wave file
	runPythonSoundAnalyzer("sound_analyzer.py", audioFileName)
}

func checkArgs() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <audiofilename.wav>\n", os.Args[0])
		os.Exit(0)
	}
}

func audioFileName() string {
	fn := os.Args[1]

	if !strings.HasSuffix(fn, ".wav") {
		fn += ".wav"
	}

	return fn
}

// waveWriter setup and return a Wave writer.
func waveWriter(waveFile *os.File, inChannels, sampleRate int) (*wave.Writer, error) {
	param := wave.WriterParam{
		Out:           waveFile,
		Channel:       inChannels,
		SampleRate:    sampleRate,
		BitsPerSample: 8, // if 16, change to WriteSample16()
	}

	return wave.NewWriter(param)
}

func startRecording(stream *portaudio.Stream, ww *wave.Writer, framesPerBuffer []byte) {

	// Ticker to show some progress
	ticker := []string{
		"-",
		"\\",
		"/",
		"|",
	}
	rand.Seed(time.Now().UnixNano())

	// Start reading from microphone
	if err := stream.Start(); err != nil {
		panic(err)
	}

	for {
		fmt.Printf("\rRecording is live now. [%v]", ticker[rand.Intn(len(ticker)-1)])

		err := stream.Read()
		if err != nil {
			panic(err)
		}

		// Write to the Wave file
		_, err = ww.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
		if err != nil {
			panic(err)
		}
	}
}

func awaitStopRecording(stream *portaudio.Stream, ww *wave.Writer) {

	key := "recording"
	var err error
	for string([]byte(key)[0]) != "s" {
		reader := bufio.NewReader(os.Stdin)

		key, err = reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
	}

	err = ww.Close()
	if err != nil {
		panic(err)
	}

	err = stream.Close()
	if err != nil {
		panic(err)
	}

	err = portaudio.Terminate()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nRecording finished.")
}

func runPythonSoundAnalyzer(script, waveFileName string) {
	c := exec.Command("python3", script, waveFileName)

	var out []byte
	var err error
	if out, err = c.Output(); err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(out))
}
