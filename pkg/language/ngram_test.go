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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeOfLowerOrderNgrams(t *testing.T) {
	n := newNgram("äbcde")
	assert.Equal(
		t,
		[]ngram{
			newNgram("äbcde"),
			newNgram("äbcd"),
			newNgram("äbc"),
			newNgram("äb"),
			newNgram("ä"),
		},
		n.rangeOfLowerOrderNgrams())
}

func TestNgram_MarshalJSON(t *testing.T) {
	serialized, err := json.Marshal(newNgram("äbcde"))
	assert.Equal(t, "\"äbcde\"", string(serialized))
	assert.Equal(t, nil, err)
}

func TestNgram_UnmarshalJSON(t *testing.T) {
	var ngram ngram
	err := json.Unmarshal([]byte("\"äbcde\""), &ngram)
	assert.Equal(t, newNgram("äbcde"), ngram)
	assert.Equal(t, nil, err)
}
