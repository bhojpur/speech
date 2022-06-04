package espeak_test

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

import "github.com/bhojpur/speech/pkg/espeak"

func ExampleTextToSpeech() {
	espeak.TextToSpeech("Hello world!", espeak.DefaultVoice, "play", nil)
	// or set an outfile name to save it
	// TextToSpeech("Hello world!", ENUSFemale, "hello-world.wav", nil)
}

// ExampleTextToSpeech_second show usage with a non-default voice.
func ExampleTextToSpeech_customVoice() {
	// output of
	//     ~$ espeak --voices=el
	//     Pty Language Age/Gender VoiceName          File          Other Languages
	//     5  el             M  greek                europe/el
	//     7  el             M  greek-mbrola-1       mb/mb-gr2
	greek := espeak.Voice{
		Languages:  "el",
		Gender:     espeak.Male,
		Name:       "greek",
		Identifier: "europe/el",
	}
	espeak.TextToSpeech("Γειά σου Κόσμε!", &greek, "play", nil)
}