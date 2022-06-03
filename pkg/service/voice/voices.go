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

var supportedLanguage = []string{"en",
	"en-UK",
	"en-AU",
	"ja",
	"de",
	"es",
	"ru",
	"ar",
	"bn",
	"cs",
	"da",
	"nl",
	"fi",
	"el",
	"hi",
	"hu",
	"id",
	"km",
	"la",
	"it",
	"no",
	"pl",
	"sk",
	"sv",
	"th",
	"tr",
	"uk",
	"vi",
	"af",
	"bg",
	"ca",
	"cy",
	"et",
	"fr",
	"gu",
	"is",
	"jv",
	"kn",
	"ko",
	"lv",
	"ml",
	"mr",
	"ms",
	"ne",
	"pt",
	"ro",
	"si",
	"sr",
	"su",
	"ta",
	"te",
	"tl",
	"ur",
	"zh",
	"sw",
	"sq",
	"my",
	"mk",
	"hy",
	"hr",
	"eo",
	"bs"}

func IsSupported(language string) bool {
	for _, item := range supportedLanguage {
		if item == language {
			return true
		}
	}
	return false
}
