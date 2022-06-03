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
	"fmt"
	"regexp"
)

type alphabet int

const (
	arabic alphabet = iota
	armenian
	bengali
	cyrillic
	devanagari
	georgian
	greek
	gujarati
	gurmukhi
	han
	hangul
	hebrew
	hiragana
	katakana
	latin
	tamil
	telugu
	thai
)

func (alphabet alphabet) matches(text string) bool {
	switch alphabet {
	case arabic:
		return arabicChars.MatchString(text)
	case armenian:
		return armenianChars.MatchString(text)
	case bengali:
		return bengaliChars.MatchString(text)
	case cyrillic:
		return cyrillicChars.MatchString(text)
	case devanagari:
		return devanagariChars.MatchString(text)
	case georgian:
		return georgianChars.MatchString(text)
	case greek:
		return greekChars.MatchString(text)
	case gujarati:
		return gujaratiChars.MatchString(text)
	case gurmukhi:
		return gurmukhiChars.MatchString(text)
	case han:
		return hanChars.MatchString(text)
	case hangul:
		return hangulChars.MatchString(text)
	case hebrew:
		return hebrewChars.MatchString(text)
	case hiragana:
		return hiraganaChars.MatchString(text)
	case katakana:
		return katakanaChars.MatchString(text)
	case latin:
		return latinChars.MatchString(text)
	case tamil:
		return tamilChars.MatchString(text)
	case telugu:
		return teluguChars.MatchString(text)
	case thai:
		return thaiChars.MatchString(text)
	default:
		return false
	}
}

func (alphabet alphabet) supportedLanguages() (languages []Language) {
	for _, language := range AllLanguages() {
		for _, script := range language.alphabets() {
			if script == alphabet {
				languages = append(languages, language)
			}
		}
	}
	return
}

func allAlphabetsSupportingSingleLanguage() map[alphabet]Language {
	alphabets := make(map[alphabet]Language)
	for _, alphabet := range allAlphabets() {
		supportedLanguages := alphabet.supportedLanguages()
		if len(supportedLanguages) == 1 {
			alphabets[alphabet] = supportedLanguages[0]
		}
	}
	return alphabets
}

func allAlphabets() []alphabet {
	alphabets := make([]alphabet, thai+1)
	for i := 0; i <= int(thai); i++ {
		alphabets[i] = alphabet(i)
	}
	return alphabets
}

func containsAlphabet(alphabets []alphabet, alphabet alphabet) bool {
	for _, a := range alphabets {
		if a == alphabet {
			return true
		}
	}
	return false
}

var (
	arabicChars     = createRegexp("Arabic")
	armenianChars   = createRegexp("Armenian")
	bengaliChars    = createRegexp("Bengali")
	cyrillicChars   = createRegexp("Cyrillic")
	devanagariChars = createRegexp("Devanagari")
	georgianChars   = createRegexp("Georgian")
	greekChars      = createRegexp("Greek")
	gujaratiChars   = createRegexp("Gujarati")
	gurmukhiChars   = createRegexp("Gurmukhi")
	hanChars        = createRegexp("Han")
	hangulChars     = createRegexp("Hangul")
	hebrewChars     = createRegexp("Hebrew")
	hiraganaChars   = createRegexp("Hiragana")
	katakanaChars   = createRegexp("Katakana")
	latinChars      = createRegexp("Latin")
	tamilChars      = createRegexp("Tamil")
	teluguChars     = createRegexp("Telugu")
	thaiChars       = createRegexp("Thai")
)

func createRegexp(charClass string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("^\\p{%v}+$", charClass))
}
