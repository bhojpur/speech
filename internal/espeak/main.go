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

/*
#cgo CFLAGS: -I/usr/include/espeak
#cgo LDFLAGS: -lportaudio -lespeak
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <speak_lib.h>

static inline void *userData(espeak_EVENT *event)  {
	if (event != NULL)
		if (event->user_data != NULL)
			return event->user_data;

	return NULL;
}

extern int mySynthCallback(short *wav, int numsamples , espeak_EVENT *events);
*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/bhojpur/speech/pkg/espeak"
)

func main() {
	log.Println("Bhojpur Speech custom text-to-speech")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	MyCustomTTS("Hello world!", "alice")
	MyCustomTTS("Hi there!", "bob")
	log.Printf("Done! Alice has %d samples, Bob has %d samples\n", len(alicesData), len(bobsData))
}

func MyCustomTTS(text, user string) {
	// this id is the address of the data as its acted on by the SynthCallback
	// function, its passed to the callback events.
	espeak.Init(espeak.Synchronous, 1024, nil, espeak.PhonemeEvents)
	espeak.SetVoiceByProps(espeak.DefaultVoice)
	espeak.NewParameters().SetVoiceParams()
	espeak.SetSynthCallback(C.mySynthCallback)
	// espeak internally feeds the id returned by Init to the user_data,
	// but you can also pass arbitrary objects
	espeak.Synth(text, espeak.CharsAuto, 0, 0, espeak.Character, nil, unsafe.Pointer(&user))
	espeak.Synchronize()

	// at this point, the data is populated, write it to a file, distort it or whatever
}

var (
	alicesData = make([]int16, 0)
	bobsData   = make([]int16, 0)
)

//export mySynthCallback
func mySynthCallback(wav *C.short, numsamples C.int, events *C.espeak_EVENT) C.int {
	if wav == nil {
		return 1
	}
	// we passed a *string in Synth, we have to unsafely cast it into a *string
	// to extract it. C.userData is defined in the header (safely dereferences
	// the *C.espeak_EVENT object).
	user := (*string)(unsafe.Pointer(C.userData(events)))
	length := int(numsamples)
	if *user == "alice" {
		alicesData = append(
			alicesData,
			(*[1 << 28]int16)(unsafe.Pointer(wav))[:length:length]...,
		)
		log.Printf("%s, you have %d samples so far\n", *user, len(alicesData))
	} else {
		bobsData = append(
			bobsData,
			(*[1 << 28]int16)(unsafe.Pointer(wav))[:length:length]...,
		)
		log.Printf("%s, you have %d samples so far\n", *user, len(bobsData))
	}
	return 0
}