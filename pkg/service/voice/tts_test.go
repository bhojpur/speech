package voice

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
	"testing"

	langdetection "github.com/bhojpur/speech/pkg/service/lang-detection"
	"github.com/bhojpur/speech/pkg/service/moderation"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	detectionService := langdetection.NewLingualDetectionService(langdetection.DefaultLanguages)

	service := NewGoTtsService("", moderation.NewFilterDefault("", ""), 1, nil, true, detectionService)
	tests := []struct {
		input    string
		language string
	}{
		{"Each language is assigned a two-letter!", "english"},
		{"每种语言分配一个两个字母!", "chinese"},
		{"Jeder Sprache wird ein Zweibuchstabe zugewiesen!", "german"},
		{"A cada idioma se le asigna una letra de dos letras!", "spanish"},
		{"Chaque langue se voit attribuer une lettre à deux lettres!", "french"},
		{"Каждому языку присваивается двухбуквенный символ!", "русский"},
	}
	for _, test := range tests {
		assert.NoError(t, service.Speak(test.input))
	}
}
