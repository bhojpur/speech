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

import "regexp"

const maxNgramLength = 5

var japaneseCharacterSet = regexp.MustCompile("^[\\p{Hiragana}\\p{Katakana}\\p{Han}]+$")
var letter = regexp.MustCompile("^\\p{L}+$")
var multipleWhitespace = regexp.MustCompile("\\s+")
var noLetter = regexp.MustCompile("^[^\\p{L}]+$")
var numbers = regexp.MustCompile("\\p{N}")
var punctuation = regexp.MustCompile("\\p{P}")
var languagesSupportingLogograms = []Language{Chinese, Japanese, Korean}

var charsToLanguagesMapping = map[string][]Language{
	"Ãã":     {Portuguese, Vietnamese},
	"ĄąĘę":   {Lithuanian, Polish},
	"Żż":     {Polish, Romanian},
	"Îî":     {French, Romanian},
	"Ññ":     {Basque, Spanish},
	"ŇňŤť":   {Czech, Slovak},
	"Ăă":     {Romanian, Vietnamese},
	"İıĞğ":   {Azerbaijani, Turkish},
	"ЈјЉљЊњ": {Macedonian, Serbian},
	"ẸẹỌọ":   {Vietnamese, Yoruba},
	"ÐðÞþ":   {Icelandic, Turkish},
	"Ûû":     {French, Hungarian},
	"Ōō":     {Maori, Yoruba},

	"ĀāĒēĪī": {Latvian, Maori, Yoruba},
	"Şş":     {Azerbaijani, Romanian, Turkish},
	"Ďď":     {Czech, Romanian, Slovak},
	"Ćć":     {Bosnian, Croatian, Polish},
	"Đđ":     {Bosnian, Croatian, Vietnamese},
	"Іі":     {Belarusian, Kazakh, Ukrainian},
	"Ìì":     {Italian, Vietnamese, Yoruba},
	"Øø":     {Bokmal, Danish, Nynorsk},

	"Ūū":     {Latvian, Lithuanian, Maori, Yoruba},
	"Ëë":     {Afrikaans, Albanian, Dutch, French},
	"ÈèÙù":   {French, Italian, Vietnamese, Yoruba},
	"Êê":     {Afrikaans, French, Portuguese, Vietnamese},
	"Õõ":     {Estonian, Hungarian, Portuguese, Vietnamese},
	"Ôô":     {French, Portuguese, Slovak, Vietnamese},
	"ЁёЫыЭэ": {Belarusian, Kazakh, Mongolian, Russian},
	"ЩщЪъ":   {Bulgarian, Kazakh, Mongolian, Russian},
	"Òò":     {Catalan, Italian, Vietnamese, Yoruba},
	"Ææ":     {Bokmal, Danish, Icelandic, Nynorsk},
	"Åå":     {Bokmal, Danish, Nynorsk, Swedish},

	"Ýý": {Czech, Icelandic, Slovak, Turkish, Vietnamese},
	"Ää": {Estonian, Finnish, German, Slovak, Swedish},
	"Àà": {Catalan, French, Italian, Portuguese, Vietnamese},
	"Ââ": {French, Portuguese, Romanian, Turkish, Vietnamese},

	"Üü":     {Azerbaijani, Catalan, Estonian, German, Hungarian, Spanish, Turkish},
	"ČčŠšŽž": {Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene},
	"Çç":     {Albanian, Azerbaijani, Basque, Catalan, French, Portuguese, Turkish},

	"Öö": {Azerbaijani, Estonian, Finnish, German, Hungarian, Icelandic, Swedish, Turkish},

	"Óó":     {Catalan, Hungarian, Icelandic, Irish, Polish, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
	"ÁáÍíÚú": {Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},

	"Éé": {Catalan, Czech, French, Hungarian, Icelandic, Irish, Italian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
}
