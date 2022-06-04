package espeak

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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTextToSpeech(t *testing.T) {
	tmp, err := ioutil.TempDir("", "bhojpur-speech-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	p := NewParameters(WithDir(tmp))
	t.Run("success", func(t *testing.T) {
		samples, err := TextToSpeech("test speech", nil, "test", p)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if samples == 0 {
			t.Errorf("0 samples written")
		}
	})
	t.Run("errors", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			text   string
			params *Parameters
			voice  *Voice
			want   error
		}{
			{"empty text", "", nil, nil, ErrEmptyText},
		} {
			t.Run(tt.name, func(t *testing.T) {
				s, err := TextToSpeech(tt.text, nil, "test", p)
				if s != 0 {
					t.Errorf("expected return samples 0 got %d", s)
				}
				if !errors.Is(err, tt.want) {

				}
				if err == nil {
					t.Error("expected an error but didn't get one")
				}
			})
		}
	})

}

func TestEnsureWavSuffix(t *testing.T) {
	for _, tt := range []struct {
		in, want string
	}{
		{"outfile.wav", "outfile.wav"},
		{"out", "out.wav"},
		{"out.mp4", "out.mp4.wav"},
		{"out.", "out.wav"},
		{"out....................", "out.wav"},
		{"out...____.", "out...____.wav"},
		{"out.", "out.wav"},
	} {
		t.Run(fmt.Sprintf("ensureWavSuffix(%q)==%q", tt.in, tt.want), func(t *testing.T) {
			got := ensureWavSuffix(tt.in)
			if got != tt.want {
				t.Errorf("expected %q got %q", tt.want, got)
			}
		})
	}
}