package language

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

// ConfidenceValue is the interface describing a language's confidence value
// that is computed by LanguageDetector.ComputeLanguageConfidenceValues.
type ConfidenceValue interface {
	// Language returns the language being part of this ConfidenceValue.
	Language() Language

	// Value returns a language's confidence value which lies between 0.0 and 1.0.
	Value() float64
}

type confidenceValue struct {
	language Language
	value    float64
}

func newConfidenceValue(language Language, value float64) ConfidenceValue {
	return confidenceValue{language, value}
}

func (c confidenceValue) Language() Language {
	return c.language
}

func (c confidenceValue) Value() float64 {
	return c.value
}
