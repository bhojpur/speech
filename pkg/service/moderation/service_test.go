package moderation

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

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	filterDefault := NewFilterDefault("", "")
	result := filterDefault.Moderate(Message{From: "shit", Text: "shit"})
	assert.Equal(t, " ", result)
}

func TestName2(t *testing.T) {
	filterDefault := NewFilterDefault("", "")
	result := filterDefault.Moderate(Message{From: "wtf", Text: "wtf"})
	assert.Equal(t, " ", result)
}

func TestName3(t *testing.T) {
	filterDefault := NewFilterDefault("", "")
	urlFilter := NewUrlFilterDecorator(filterDefault)
	type testData struct {
		input    Message
		expected string
	}

	tests := []testData{
		{expected: "http:::/not.valid/a//a??a?b=&&c#hi  ", input: Message{From: "1", Text: "http:::/not.valid/a//a??a?b=&&c#hi"}},
		{expected: " ", input: Message{From: "1", Text: "http//google.com"}},
		{expected: " ", input: Message{From: "1", Text: "google.com"}},
		{expected: " hello  ", input: Message{From: "1", Text: "wtf google.com hello"}},
		{expected: "/foo/bar  ", input: Message{From: "1", Text: "/foo/bar"}},
		{expected: "http://  ", input: Message{From: "1", Text: "http://"}},
		{expected: " message send by me  ", input: Message{From: "1", Text: " message send by me"}},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, urlFilter.Moderate(test.input))
	}
}

func TestName4(t *testing.T) {
	filterDefault := NewFilterDefault("", "")
	urlFilter := NewUrlFilterDecorator(filterDefault)
	result := urlFilter.Moderate(Message{From: "", Text: "Adjust position, velocity, accel?"})
	assert.Equal(t, "Adjust position, velocity, accel?  ", result)
	userFilter := NewUserFilterDecorator(urlFilter, []string{"spin"})
	result = userFilter.Moderate(Message{From: "spin", Text: "Adjust position, velocity, accel?"})
	assert.Equal(t, "", result)
}
