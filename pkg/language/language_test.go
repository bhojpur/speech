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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllLanguages(t *testing.T) {
	assert.Equal(
		t,
		[]Language{
			Afrikaans,
			Albanian,
			Arabic,
			Armenian,
			Azerbaijani,
			Basque,
			Belarusian,
			Bengali,
			Bokmal,
			Bosnian,
			Bulgarian,
			Catalan,
			Chinese,
			Croatian,
			Czech,
			Danish,
			Dutch,
			English,
			Esperanto,
			Estonian,
			Finnish,
			French,
			Ganda,
			Georgian,
			German,
			Greek,
			Gujarati,
			Hebrew,
			Hindi,
			Hungarian,
			Icelandic,
			Indonesian,
			Irish,
			Italian,
			Japanese,
			Kazakh,
			Korean,
			Latin,
			Latvian,
			Lithuanian,
			Macedonian,
			Malay,
			Maori,
			Marathi,
			Mongolian,
			Nynorsk,
			Persian,
			Polish,
			Portuguese,
			Punjabi,
			Romanian,
			Russian,
			Serbian,
			Shona,
			Slovak,
			Slovene,
			Somali,
			Sotho,
			Spanish,
			Swahili,
			Swedish,
			Tagalog,
			Tamil,
			Telugu,
			Thai,
			Tsonga,
			Tswana,
			Turkish,
			Ukrainian,
			Urdu,
			Vietnamese,
			Welsh,
			Xhosa,
			Yoruba,
			Zulu,
		},
		AllLanguages())
}

func TestAllSpokenLanguages(t *testing.T) {
	assert.Equal(
		t,
		[]Language{
			Afrikaans,
			Albanian,
			Arabic,
			Armenian,
			Azerbaijani,
			Basque,
			Belarusian,
			Bengali,
			Bokmal,
			Bosnian,
			Bulgarian,
			Catalan,
			Chinese,
			Croatian,
			Czech,
			Danish,
			Dutch,
			English,
			Esperanto,
			Estonian,
			Finnish,
			French,
			Ganda,
			Georgian,
			German,
			Greek,
			Gujarati,
			Hebrew,
			Hindi,
			Hungarian,
			Icelandic,
			Indonesian,
			Irish,
			Italian,
			Japanese,
			Kazakh,
			Korean,
			Latvian,
			Lithuanian,
			Macedonian,
			Malay,
			Maori,
			Marathi,
			Mongolian,
			Nynorsk,
			Persian,
			Polish,
			Portuguese,
			Punjabi,
			Romanian,
			Russian,
			Serbian,
			Shona,
			Slovak,
			Slovene,
			Somali,
			Sotho,
			Spanish,
			Swahili,
			Swedish,
			Tagalog,
			Tamil,
			Telugu,
			Thai,
			Tsonga,
			Tswana,
			Turkish,
			Ukrainian,
			Urdu,
			Vietnamese,
			Welsh,
			Xhosa,
			Yoruba,
			Zulu,
		},
		AllSpokenLanguages())
}

func TestAllLanguagesWithArabicScript(t *testing.T) {
	assert.Equal(t, []Language{Arabic, Persian, Urdu}, AllLanguagesWithArabicScript())
}

func TestAllLanguagesWithCyrillicScript(t *testing.T) {
	assert.Equal(
		t,
		[]Language{
			Belarusian, Bulgarian, Kazakh, Macedonian, Mongolian, Russian, Serbian, Ukrainian,
		},
		AllLanguagesWithCyrillicScript())
}

func TestAllLanguagesWithDevanagariScript(t *testing.T) {
	assert.Equal(t, []Language{Hindi, Marathi}, AllLanguagesWithDevanagariScript())
}

func TestAllLanguagesWithLatinScript(t *testing.T) {
	assert.Equal(
		t,
		[]Language{
			Afrikaans,
			Albanian,
			Azerbaijani,
			Basque,
			Bokmal,
			Bosnian,
			Catalan,
			Croatian,
			Czech,
			Danish,
			Dutch,
			English,
			Esperanto,
			Estonian,
			Finnish,
			French,
			Ganda,
			German,
			Hungarian,
			Icelandic,
			Indonesian,
			Irish,
			Italian,
			Latin,
			Latvian,
			Lithuanian,
			Malay,
			Maori,
			Nynorsk,
			Polish,
			Portuguese,
			Romanian,
			Shona,
			Slovak,
			Slovene,
			Somali,
			Sotho,
			Spanish,
			Swahili,
			Swedish,
			Tagalog,
			Tsonga,
			Tswana,
			Turkish,
			Vietnamese,
			Welsh,
			Xhosa,
			Yoruba,
			Zulu,
		},
		AllLanguagesWithLatinScript())
}

func TestLanguage_MarshalJSON(t *testing.T) {
	language, err := json.Marshal(German)
	assert.Equal(t, "\"GERMAN\"", string(language))
	assert.Equal(t, nil, err)
}

func TestLanguage_UnmarshalJSON(t *testing.T) {
	var language Language
	err := json.Unmarshal([]byte("\"GERMAN\""), &language)
	assert.Equal(t, German, language)
	assert.Equal(t, nil, err)

	err = json.Unmarshal([]byte("\"GERM\""), &language)
	assert.Equal(t, fmt.Errorf("string \"GERM\" cannot be unmarshalled to an instance of type Language"), err)
}
