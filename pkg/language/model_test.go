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
	"bufio"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const text = `These sentences are intended for testing purposes.
 ⚠ Do not use them in production
 By the way, they consist of 23 words in total.`

func linesOfText(text string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		lines = append(lines, strings.ToLower(scanner.Text()))
	}
	return lines
}

var expectedUnigrams = mapStringsToNgrams(
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "l",
	"m", "n", "o", "p", "r", "s", "t", "u", "w", "y")

var expectedUnigramAbsoluteFrequencies = mapKeysToNgrams(map[string]uint32{
	"a": 3, "b": 1, "c": 3, "d": 5, "e": 14, "f": 2, "g": 1, "h": 4, "i": 6,
	"l": 1, "m": 1, "n": 10, "o": 10, "p": 3, "r": 5, "s": 10, "t": 13,
	"u": 3, "w": 2, "y": 3,
})

var expectedUnigramRelativeFrequencies = mapKeysToNgramsAndValuesToRats(map[string]string{
	"a": "3/100", "b": "1/100", "c": "3/100", "d": "1/20", "e": "7/50",
	"f": "1/50", "g": "1/100", "h": "1/25", "i": "3/50", "l": "1/100",
	"m": "1/100", "n": "1/10", "o": "1/10", "p": "3/100", "r": "1/20",
	"s": "1/10", "t": "13/100", "u": "3/100", "w": "1/50", "y": "3/100",
})

func expectedUnigramJsonRelativeFrequencies() map[ngram]float64 {
	frequencies := make(map[ngram]float64)
	for ngram, rat := range expectedUnigramRelativeFrequencies {
		frequencies[ngram], _ = rat.Float64()
	}
	return frequencies
}

var expectedBigrams = mapStringsToNgrams(
	"de", "pr", "pu", "do", "uc", "ds", "du", "ur", "us", "ed", "in",
	"io", "em", "en", "is", "al", "es", "ar", "rd", "re", "ey", "nc", "nd",
	"ay", "ng", "ro", "rp", "no", "ns", "nt", "fo", "wa", "se", "od", "si",
	"by", "of", "wo", "on", "st", "ce", "or", "os", "ot", "co", "ta", "te",
	"ct", "th", "ti", "to", "he", "po")

var expectedBigramAbsoluteFrequencies = mapKeysToNgrams(map[string]uint32{
	"de": 1, "pr": 1, "pu": 1, "do": 1, "uc": 1, "ds": 1, "du": 1,
	"ur": 1, "us": 1, "ed": 1, "in": 4, "io": 1, "em": 1, "en": 3,
	"is": 1, "al": 1, "es": 4, "ar": 1, "rd": 1, "re": 1, "ey": 1,
	"nc": 1, "nd": 1, "ay": 1, "ng": 1, "ro": 1, "rp": 1, "no": 1,
	"ns": 1, "nt": 2, "fo": 1, "wa": 1, "se": 4, "od": 1, "si": 1,
	"of": 1, "by": 1, "wo": 1, "on": 2, "st": 2, "ce": 1, "or": 2,
	"os": 1, "ot": 2, "co": 1, "ta": 1, "ct": 1, "te": 3, "th": 4,
	"ti": 2, "to": 1, "he": 4, "po": 1,
})

var expectedBigramRelativeFrequencies = mapKeysToNgramsAndValuesToRats(map[string]string{
	"de": "1/5", "pr": "1/3", "pu": "1/3", "do": "1/5", "uc": "1/3",
	"ds": "1/5", "du": "1/5", "ur": "1/3", "us": "1/3", "ed": "1/14",
	"in": "2/3", "io": "1/6", "em": "1/14", "en": "3/14", "is": "1/6",
	"al": "1/3", "es": "2/7", "ar": "1/3", "rd": "1/5", "re": "1/5",
	"ey": "1/14", "nc": "1/10", "nd": "1/10", "ay": "1/3", "ng": "1/10",
	"ro": "1/5", "rp": "1/5", "no": "1/10", "ns": "1/10", "nt": "1/5",
	"fo": "1/2", "wa": "1/2", "se": "2/5", "od": "1/10", "si": "1/10",
	"of": "1/10", "by": "1/1", "wo": "1/2", "on": "1/5", "st": "1/5",
	"ce": "1/3", "or": "1/5", "os": "1/10", "ot": "1/5", "co": "1/3",
	"ta": "1/13", "ct": "1/3", "te": "3/13", "th": "4/13", "ti": "2/13",
	"to": "1/13", "he": "1/1", "po": "1/3",
})

var expectedTrigrams = mapStringsToNgrams(
	"rds", "ose", "ded", "con", "use", "est", "ion", "ist", "pur",
	"hem", "hes", "tin", "cti", "tio", "wor", "ten", "hey", "ota", "tal",
	"tes", "uct", "sti", "pro", "odu", "nsi", "rod", "for", "ces", "nce",
	"not", "are", "pos", "tot", "end", "enc", "sis", "sen", "nte", "ses",
	"ord", "ing", "ent", "int", "nde", "way", "the", "rpo", "urp", "duc",
	"ons", "ese")

var expectedTrigramAbsoluteFrequencies = mapKeysToNgrams(map[string]uint32{
	"rds": 1, "ose": 1, "ded": 1, "con": 1, "use": 1, "est": 1, "ion": 1,
	"ist": 1, "pur": 1, "hem": 1, "hes": 1, "tin": 1, "cti": 1, "wor": 1,
	"tio": 1, "ten": 2, "ota": 1, "hey": 1, "tal": 1, "tes": 1, "uct": 1,
	"sti": 1, "pro": 1, "odu": 1, "nsi": 1, "rod": 1, "for": 1, "ces": 1,
	"nce": 1, "not": 1, "pos": 1, "are": 1, "tot": 1, "end": 1, "enc": 1,
	"sis": 1, "sen": 1, "nte": 2, "ord": 1, "ses": 1, "ing": 1, "ent": 1,
	"way": 1, "nde": 1, "int": 1, "rpo": 1, "the": 4, "urp": 1, "duc": 1,
	"ons": 1, "ese": 1,
})

var expectedTrigramRelativeFrequencies = mapKeysToNgramsAndValuesToRats(map[string]string{
	"rds": "1/1", "ose": "1/1", "ded": "1/1", "con": "1/1", "use": "1/1",
	"est": "1/4", "ion": "1/1", "ist": "1/1", "pur": "1/1", "hem": "1/4",
	"hes": "1/4", "tin": "1/2", "cti": "1/1", "wor": "1/1", "tio": "1/2",
	"ten": "2/3", "ota": "1/2", "hey": "1/4", "tal": "1/1", "tes": "1/3",
	"uct": "1/1", "sti": "1/2", "pro": "1/1", "odu": "1/1", "nsi": "1/1",
	"rod": "1/1", "for": "1/1", "ces": "1/1", "nce": "1/1", "not": "1/1",
	"pos": "1/1", "are": "1/1", "tot": "1/1", "end": "1/3", "enc": "1/3",
	"sis": "1/1", "sen": "1/4", "nte": "1/1", "ord": "1/2", "ses": "1/4",
	"ing": "1/4", "ent": "1/3", "way": "1/1", "nde": "1/1", "int": "1/4",
	"rpo": "1/1", "the": "1/1", "urp": "1/1", "duc": "1/1", "ons": "1/2",
	"ese": "1/4",
})

var expectedQuadrigrams = mapStringsToNgrams(
	"onsi", "sist", "ende", "ords", "esti", "tenc", "nces", "oduc",
	"tend", "thes", "rpos", "ting", "nten", "nsis", "they", "tota", "cons",
	"tion", "prod", "ence", "test", "otal", "pose", "nded", "oses", "inte",
	"urpo", "them", "sent", "duct", "stin", "ente", "ucti", "purp", "ctio",
	"rodu", "word", "hese")

var expectedQuadrigramAbsoluteFrequencies = mapKeysToNgrams(map[string]uint32{
	"onsi": 1, "sist": 1, "ende": 1, "ords": 1, "esti": 1, "oduc": 1,
	"nces": 1, "tenc": 1, "tend": 1, "thes": 1, "rpos": 1, "ting": 1,
	"nsis": 1, "nten": 2, "tota": 1, "they": 1, "cons": 1, "tion": 1,
	"prod": 1, "otal": 1, "test": 1, "ence": 1, "pose": 1, "oses": 1,
	"nded": 1, "inte": 1, "them": 1, "urpo": 1, "duct": 1, "sent": 1,
	"stin": 1, "ucti": 1, "ente": 1, "purp": 1, "ctio": 1, "rodu": 1,
	"word": 1, "hese": 1,
})

var expectedQuadrigramRelativeFrequencies = mapKeysToNgramsAndValuesToRats(map[string]string{
	"onsi": "1/1", "sist": "1/1", "ende": "1/1", "ords": "1/1",
	"esti": "1/1", "oduc": "1/1", "nces": "1/1", "tenc": "1/2",
	"tend": "1/2", "thes": "1/4", "rpos": "1/1", "ting": "1/1",
	"nsis": "1/1", "nten": "1/1", "tota": "1/1", "they": "1/4",
	"cons": "1/1", "tion": "1/1", "prod": "1/1", "otal": "1/1",
	"test": "1/1", "ence": "1/1", "pose": "1/1", "oses": "1/1",
	"nded": "1/1", "inte": "1/1", "them": "1/4", "urpo": "1/1",
	"duct": "1/1", "sent": "1/1", "stin": "1/1", "ucti": "1/1",
	"ente": "1/1", "purp": "1/1", "ctio": "1/1", "rodu": "1/1",
	"word": "1/1", "hese": "1/1",
})

var expectedFivegrams = mapStringsToNgrams(
	"testi", "sente", "ences", "tende", "these", "ntenc", "ducti",
	"ntend", "onsis", "total", "uctio", "enten", "poses", "ction", "produ",
	"inten", "nsist", "words", "sting", "tence", "purpo", "estin", "roduc",
	"urpos", "ended", "rpose", "oduct", "consi")

var expectedFivegramAbsoluteFrequencies = mapKeysToNgrams(map[string]uint32{
	"testi": 1, "sente": 1, "ences": 1, "tende": 1, "ducti": 1,
	"ntenc": 1, "these": 1, "onsis": 1, "ntend": 1, "total": 1,
	"uctio": 1, "enten": 1, "poses": 1, "ction": 1, "produ": 1,
	"inten": 1, "nsist": 1, "words": 1, "sting": 1, "purpo": 1,
	"tence": 1, "estin": 1, "roduc": 1, "urpos": 1, "rpose": 1,
	"ended": 1, "oduct": 1, "consi": 1,
})

var expectedFivegramRelativeFrequencies = mapKeysToNgramsAndValuesToRats(map[string]string{
	"testi": "1/1", "sente": "1/1", "ences": "1/1", "tende": "1/1",
	"ducti": "1/1", "ntenc": "1/2", "these": "1/1", "onsis": "1/1",
	"ntend": "1/2", "total": "1/1", "uctio": "1/1", "enten": "1/1",
	"poses": "1/1", "ction": "1/1", "produ": "1/1", "inten": "1/1",
	"nsist": "1/1", "words": "1/1", "sting": "1/1", "purpo": "1/1",
	"tence": "1/1", "estin": "1/1", "roduc": "1/1", "urpos": "1/1",
	"rpose": "1/1", "ended": "1/1", "oduct": "1/1", "consi": "1/1",
})

func TestNewTrainingDataLanguageModel(t *testing.T) {
	params := []struct {
		ngramLength                   int
		expectedAbsoluteFrequencies   map[ngram]uint32
		expectedRelativeFrequencies   map[ngram]*big.Rat
		lowerNgramAbsoluteFrequencies map[ngram]uint32
	}{
		{1, expectedUnigramAbsoluteFrequencies, expectedUnigramRelativeFrequencies, map[ngram]uint32{}},
		{2, expectedBigramAbsoluteFrequencies, expectedBigramRelativeFrequencies, expectedUnigramAbsoluteFrequencies},
		{3, expectedTrigramAbsoluteFrequencies, expectedTrigramRelativeFrequencies, expectedBigramAbsoluteFrequencies},
		{4, expectedQuadrigramAbsoluteFrequencies, expectedQuadrigramRelativeFrequencies, expectedTrigramAbsoluteFrequencies},
		{5, expectedFivegramAbsoluteFrequencies, expectedFivegramRelativeFrequencies, expectedQuadrigramAbsoluteFrequencies},
	}
	for i := range params {
		model := newTrainingDataLanguageModel(
			linesOfText(text),
			English,
			params[i].ngramLength,
			"\\p{L}&&\\p{Latin}",
			params[i].lowerNgramAbsoluteFrequencies,
		)
		assert.Equal(t, English, model.language)
		assert.Equal(t, params[i].expectedAbsoluteFrequencies, model.absoluteFrequencies)
		assert.Equal(t, params[i].expectedRelativeFrequencies, model.relativeFrequencies)

	}
}

func TestNewTrainingDataLanguageModelFromJson(t *testing.T) {
	model := trainingDataLanguageModel{
		language:                English,
		absoluteFrequencies:     nil,
		relativeFrequencies:     expectedUnigramRelativeFrequencies,
		jsonRelativeFrequencies: nil,
	}
	deserialized := newTrainingDataLanguageModelFromJson(model.toJson())
	assert.Equal(t, English, deserialized.language)
	assert.Equal(t, map[ngram]uint32(nil), deserialized.absoluteFrequencies)
	assert.Equal(t, map[ngram]*big.Rat(nil), deserialized.relativeFrequencies)
	assert.Equal(t, expectedUnigramJsonRelativeFrequencies(), deserialized.jsonRelativeFrequencies)
}

func TestNewTestDataLanguageModel(t *testing.T) {
	params := []struct {
		ngramLength    int
		expectedNgrams map[ngram]bool
	}{
		{1, expectedUnigrams},
		{2, expectedBigrams},
		{3, expectedTrigrams},
		{4, expectedQuadrigrams},
		{5, expectedFivegrams},
	}
	for i := range params {
		model := newTestDataLanguageModel(strings.ToLower(text), params[i].ngramLength)
		assert.Equal(t, params[i].expectedNgrams, model.ngrams)
	}
}

func mapStringsToNgrams(strings ...string) map[ngram]bool {
	ngrams := make(map[ngram]bool)
	for _, s := range strings {
		ngrams[newNgram(s)] = true
	}
	return ngrams
}

func mapKeysToNgrams(m map[string]uint32) map[ngram]uint32 {
	ngrams := make(map[ngram]uint32)
	for key, value := range m {
		ngrams[newNgram(key)] = value
	}
	return ngrams
}

func mapKeysToNgramsAndValuesToRats(m map[string]string) map[ngram]*big.Rat {
	ngrams := make(map[ngram]*big.Rat)
	for key, value := range m {
		rat := new(big.Rat)
		rat.SetString(value)
		ngrams[newNgram(key)] = rat
	}
	return ngrams
}
