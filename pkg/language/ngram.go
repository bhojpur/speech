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

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

type ngram struct {
	value string
}

type ngramSlice []ngram

func newNgram(value string) ngram {
	charCount := utf8.RuneCountInString(value)
	if charCount > maxNgramLength {
		panic(fmt.Sprintf("length %v of ngram '%v' is greater than %v", charCount, value, maxNgramLength))
	}
	return ngram{value: value}
}

func getNgramNameByLength(ngramLength int) string {
	switch ngramLength {
	case 1:
		return "unigram"
	case 2:
		return "bigram"
	case 3:
		return "trigram"
	case 4:
		return "quadrigram"
	case 5:
		return "fivegram"
	default:
		panic(fmt.Sprintf("ngram length %v is not in range 1..5", ngramLength))
	}
}

func (n ngram) rangeOfLowerOrderNgrams() []ngram {
	var ngrams []ngram
	chars := []rune(n.value)

	for i := len(chars); i > 0; i-- {
		ngrams = append(ngrams, newNgram(string(chars[:i])))
	}

	return ngrams
}

func (n ngram) String() string {
	return n.value
}

func (n ngram) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n *ngram) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	*n = newNgram(s)
	return nil
}

func (ngrams ngramSlice) Len() int {
	return len(ngrams)
}

func (ngrams ngramSlice) Less(i, j int) bool {
	return ngrams[i].value < ngrams[j].value
}

func (ngrams ngramSlice) Swap(i, j int) {
	ngrams[i], ngrams[j] = ngrams[j], ngrams[i]
}
