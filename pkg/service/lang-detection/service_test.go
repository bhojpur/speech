package lang_detection

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
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input    string
	expected Language
}

func TestDetermineLanguage(t *testing.T) {
	tests := []testCase{
		{" Each language is assigned a two-letter ", "en"},
		{"你好，世界", "zh"},
		{"Hallo Welt", "de"},
		{"A cada idioma se le asigna una letra de dos letras", "es"},
		{"Chaque langue se voit attribuer une lettre à deux lettres", "fr"},
		{"Каждому языку присваивается двухбуквенный", "ru"},
	}

	service := NewLingualDetectionService(DefaultLanguages)
	for _, test := range tests {
		lang, err := service.Detect(test.input)
		fmt.Println(lang, err, test.expected)
		assert.Equal(t, test.expected, *lang)
	}
}
