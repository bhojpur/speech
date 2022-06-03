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
	"errors"
	"strings"

	spoken "github.com/bhojpur/speech/pkg/language"
)

var DefaultLanguages = []spoken.Language{
	spoken.English,
	spoken.Hindi,
	spoken.Spanish,
	spoken.French,
	spoken.German,
	spoken.Russian,
}

type Language string

var ErrUnsupportedLanguage = errors.New("unsupported language")

type LanguageDetectionService interface {
	Detect(text string) (*Language, error)
}

type LingualDetectionService struct {
	detector spoken.LanguageDetector
}

func NewLingualDetectionService(languages []spoken.Language) *LingualDetectionService {
	return &BhojpurDetectionService{
		detector: spoken.NewLanguageDetectorBuilder().
			FromLanguages(languages...).
			Build(),
	}
}

func (s *LingualDetectionService) Detect(text string) (*Language, error) {
	lang, ok := s.detector.DetectLanguageOf(text)
	if !ok {
		return nil, ErrUnsupportedLanguage
	}

	language := Language(strings.ToLower(lang.IsoCode639_1().String()))
	return &language, nil
}
