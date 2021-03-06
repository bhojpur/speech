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
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ##############################
// MOCKS
// ##############################

type mockedTrainingDataLanguageModel struct {
	mock.Mock
}

func (m *mockedTrainingDataLanguageModel) getRelativeFrequency(ngram ngram) float64 {
	return m.Called(ngram).Get(0).(float64)
}

func createTrainingModelMock(data map[string]float64) *mockedTrainingDataLanguageModel {
	model := new(mockedTrainingDataLanguageModel)
	for ngram, probability := range data {
		model.On("getRelativeFrequency", newNgram(ngram)).Return(probability)
	}
	return model
}

// ##############################
// LANGUAGE MODELS FOR ENGLISH
// ##############################

func unigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"a": 0.01,
		"l": 0.02,
		"t": 0.03,
		"e": 0.04,
		"r": 0.05,
		// unknown unigrams
		"w": 0.0,
	})

}

func bigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"al": 0.11,
		"lt": 0.12,
		"te": 0.13,
		"er": 0.14,
		// unknown bigrams
		"aq": 0.0,
		"wx": 0.0,
	})
}

func trigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alt": 0.19,
		"lte": 0.2,
		"ter": 0.21,
		// unknown trigrams
		"aqu": 0.0,
		"tez": 0.0,
		"wxy": 0.0,
	})
}

func quadrigramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alte": 0.25,
		"lter": 0.26,
		// unknown quadrigrams
		"aqua": 0.0,
		"wxyz": 0.0,
	})
}

func fivegramModelForEnglish() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alter": 0.29,
		// unknown fivegrams
		"aquas": 0.0,
	})
}

// ##############################
// LANGUAGE MODELS FOR GERMAN
// ##############################

func unigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"a": 0.06,
		"l": 0.07,
		"t": 0.08,
		"e": 0.09,
		"r": 0.1,
		// unknown unigrams
		"w": 0.0,
	})
}

func bigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"al": 0.15,
		"lt": 0.16,
		"te": 0.17,
		"er": 0.18,
		// unknown bigrams
		"wx": 0.0,
	})
}

func trigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alt": 0.22,
		"lte": 0.23,
		"ter": 0.24,
		// unknown trigrams
		"wxy": 0.0,
	})
}

func quadrigramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alte": 0.27,
		"lter": 0.28,
		// unknown quadrigrams
		"wxyz": 0.0,
	})
}

func fivegramModelForGerman() languageModel {
	return createTrainingModelMock(map[string]float64{
		"alter": 0.3,
	})
}

// ##############################
// TEST DATA MODELS
// ##############################

func testDataModel(strs []string) testDataLanguageModel {
	ngrams := make(map[ngram]bool)
	for _, s := range strs {
		ngrams[newNgram(s)] = true
	}
	return testDataLanguageModel{ngrams}
}

// ##############################
// DETECTORS
// ##############################

func newDetectorForEnglishAndGerman() languageDetector {
	var unigramLanguageModels sync.Map
	unigramLanguageModels.Store(English, unigramModelForEnglish())
	unigramLanguageModels.Store(German, unigramModelForGerman())

	var bigramLanguageModels sync.Map
	bigramLanguageModels.Store(English, bigramModelForEnglish())
	bigramLanguageModels.Store(German, bigramModelForGerman())

	var trigramLanguageModels sync.Map
	trigramLanguageModels.Store(English, trigramModelForEnglish())
	trigramLanguageModels.Store(German, trigramModelForGerman())

	var quadrigramLanguageModels sync.Map
	quadrigramLanguageModels.Store(English, quadrigramModelForEnglish())
	quadrigramLanguageModels.Store(German, quadrigramModelForGerman())

	var fivegramLanguageModels sync.Map
	fivegramLanguageModels.Store(English, fivegramModelForEnglish())
	fivegramLanguageModels.Store(German, fivegramModelForGerman())

	return languageDetector{
		languages:                     []Language{English, German},
		minimumRelativeDistance:       0.0,
		languagesWithUniqueCharacters: []Language{},
		oneLanguageAlphabets:          map[alphabet]Language{},
		unigramLanguageModels:         &unigramLanguageModels,
		bigramLanguageModels:          &bigramLanguageModels,
		trigramLanguageModels:         &trigramLanguageModels,
		quadrigramLanguageModels:      &quadrigramLanguageModels,
		fivegramLanguageModels:        &fivegramLanguageModels,
	}
}

func newDetectorForAllLanguages() languageDetector {
	languages := AllLanguages()
	var emptyLanguageModels sync.Map
	return languageDetector{
		languages:                     languages,
		minimumRelativeDistance:       0.0,
		languagesWithUniqueCharacters: collectLanguagesWithUniqueCharacters(languages),
		oneLanguageAlphabets:          collectOneLanguageAlphabets(languages),
		unigramLanguageModels:         &emptyLanguageModels,
		bigramLanguageModels:          &emptyLanguageModels,
		trigramLanguageModels:         &emptyLanguageModels,
		quadrigramLanguageModels:      &emptyLanguageModels,
		fivegramLanguageModels:        &emptyLanguageModels,
	}
}

var detectorForEnglishAndGerman = newDetectorForEnglishAndGerman()
var detectorForAllLanguages = newDetectorForAllLanguages()

// ##############################
// TESTS
// ##############################

var delta = 0.00000000000001

func TestCleanUpInputText(t *testing.T) {
	text := `Weltweit    gibt es ungef??hr 6.000 Sprachen,
	 wobei laut Sch??tzungen zufolge ungef??hr 90  Prozent davon
	 am Ende dieses Jahrhunderts verdr??ngt sein werden.`

	expectedCleanedText := "weltweit gibt es ungef??hr sprachen wobei laut sch??tzungen " +
		"zufolge ungef??hr prozent davon am ende dieses jahrhunderts verdr??ngt sein werden"

	assert.Equal(t, expectedCleanedText, detectorForAllLanguages.cleanUpInputText(text))
}

func TestSplitTextIntoWords(t *testing.T) {
	testCases := []struct {
		text          string
		expectedWords []string
	}{
		{
			"this is a sentence",
			[]string{"this", "is", "a", "sentence"},
		},
		{
			"sentence",
			[]string{"sentence"},
		},
		{
			"?????????????????????????????? this is a sentence",
			[]string{"???", "???", "???", "???", "???", "???", "???", "???", "???", "???", "this", "is", "a", "sentence"},
		},
	}
	for _, testCase := range testCases {
		assert.Equal(
			t,
			testCase.expectedWords,
			detectorForAllLanguages.splitTextIntoWords(testCase.text),
			fmt.Sprintf("unexpected tokenization for text '%s'", testCase.text),
		)
	}
}

func TestLookUpNgramProbability(t *testing.T) {
	testCases := []struct {
		language            Language
		ngram               string
		expectedProbability float64
	}{
		{English, "a", 0.01},
		{English, "lt", 0.12},
		{English, "ter", 0.21},
		{English, "alte", 0.25},
		{English, "alter", 0.29},
		{German, "t", 0.08},
		{German, "er", 0.18},
		{German, "alt", 0.22},
		{German, "lter", 0.28},
		{German, "alter", 0.3},
	}
	for _, testCase := range testCases {
		probability := detectorForEnglishAndGerman.lookUpNgramProbability(testCase.language, newNgram(testCase.ngram))
		message := fmt.Sprintf(
			"expected probability %v for language %v and ngram '%s', got %v",
			testCase.expectedProbability,
			testCase.language,
			testCase.ngram,
			probability,
		)
		assert.Equal(t, testCase.expectedProbability, probability, message)
	}

	assert.Panicsf(t, func() {
		detectorForEnglishAndGerman.lookUpNgramProbability(English, newNgram(""))
	}, "zerogram detected")
}

func TestComputeSumOfNgramProbabilities(t *testing.T) {
	testCases := []struct {
		ngrams                     []string
		expectedSumOfProbabilities float64
	}{
		{
			[]string{"a", "l", "t", "e", "r"},
			math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05),
		},
		{
			// back off unknown Trigram("tez") to known Bigram("te")
			[]string{"alt", "lte", "tez"},
			math.Log(0.19) + math.Log(0.2) + math.Log(0.13),
		},
		{
			// back off unknown Fivegram("aquas") to known Unigram("a")
			[]string{"aquas"},
			math.Log(0.01),
		},
	}
	for _, testCase := range testCases {
		mappedNgrams := make(map[ngram]bool)
		for _, ngram := range testCase.ngrams {
			mappedNgrams[newNgram(ngram)] = true
		}
		sumOfProbabilities := detectorForEnglishAndGerman.computeSumOfNgramProbabilities(English, mappedNgrams)
		message := fmt.Sprintf(
			"expected sum %v for language %v and ngrams %v, got %v",
			testCase.expectedSumOfProbabilities,
			English,
			testCase.ngrams,
			sumOfProbabilities,
		)
		assert.InDelta(t, testCase.expectedSumOfProbabilities, sumOfProbabilities, delta, message)
	}
}

func TestComputeLanguageProbabilities(t *testing.T) {
	testCases := []struct {
		testDataModel         testDataLanguageModel
		expectedProbabilities map[Language]float64
	}{
		{
			testDataModel([]string{"a", "l", "t", "e", "r"}),
			map[Language]float64{
				English: math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05),
				German:  math.Log(0.06) + math.Log(0.07) + math.Log(0.08) + math.Log(0.09) + math.Log(0.1),
			},
		},
		{
			testDataModel([]string{"alt", "lte", "ter", "wxy"}),
			map[Language]float64{
				English: math.Log(0.19) + math.Log(0.2) + math.Log(0.21),
				German:  math.Log(0.22) + math.Log(0.23) + math.Log(0.24),
			},
		},
		{
			testDataModel([]string{"alte", "lter", "wxyz"}),
			map[Language]float64{
				English: math.Log(0.25) + math.Log(0.26),
				German:  math.Log(0.27) + math.Log(0.28),
			},
		},
	}
	languages := []Language{English, German}
	for _, testCase := range testCases {
		probabilities := detectorForEnglishAndGerman.computeLanguageProbabilities(testCase.testDataModel, languages)

		for language, probability := range probabilities {
			expectedProbability := testCase.expectedProbabilities[language]
			message := fmt.Sprintf(
				"expected probability %v for language %v, got %v",
				expectedProbability,
				language,
				probability,
			)
			assert.InDelta(t, expectedProbability, probability, delta, message)
		}
	}
}

func TestComputeLanguageConfidenceValues(t *testing.T) {
	unigramCountForBothLanguages := 5.0
	totalProbabilityForGerman := (
	// unigrams
	math.Log(0.06) + math.Log(0.07) + math.Log(0.08) + math.Log(0.09) + math.Log(0.1) +
		// bigrams
		math.Log(0.15) + math.Log(0.16) + math.Log(0.17) + math.Log(0.18) +
		// trigrams
		math.Log(0.22) + math.Log(0.23) + math.Log(0.24) +
		// quadrigrams
		math.Log(0.27) + math.Log(0.28) +
		// fivegrams
		math.Log(0.3)) / unigramCountForBothLanguages

	totalProbabilityForEnglish := (
	// unigrams
	math.Log(0.01) + math.Log(0.02) + math.Log(0.03) + math.Log(0.04) + math.Log(0.05) +
		// bigrams
		math.Log(0.11) + math.Log(0.12) + math.Log(0.13) + math.Log(0.14) +
		// trigrams
		math.Log(0.19) + math.Log(0.2) + math.Log(0.21) +
		// quadrigrams
		math.Log(0.25) + math.Log(0.26) +
		// fivegrams
		math.Log(0.29)) / unigramCountForBothLanguages

	expectedConfidenceForGerman := 1.0
	expectedConfidenceForEnglish := totalProbabilityForGerman / totalProbabilityForEnglish

	confidenceValues := detectorForEnglishAndGerman.ComputeLanguageConfidenceValues("Alter")

	assert.Equal(
		t,
		2,
		len(confidenceValues),
		fmt.Sprintf("expected 2 confidence values, got %v", len(confidenceValues)),
	)

	first, second := confidenceValues[0], confidenceValues[1]

	assert.Equal(t, German, first.Language())
	assert.Equal(t, expectedConfidenceForGerman, first.Value())

	assert.Equal(t, English, second.Language())
	assert.InDelta(t, expectedConfidenceForEnglish, second.Value(), delta)
}

func TestDetectLanguage(t *testing.T) {
	language, exists := detectorForEnglishAndGerman.DetectLanguageOf("Alter")
	assert.Equal(t, German, language)
	assert.True(t, exists)

	language, exists = detectorForEnglishAndGerman.DetectLanguageOf("??????????????????")
	assert.Equal(t, Unknown, language)
	assert.False(t, exists)
}

func TestDetectLanguageWithRules(t *testing.T) {
	testCases := []struct {
		word             string
		expectedLanguage Language
	}{
		// words with unique characters
		{"m??h??rr??m", Azerbaijani},
		{"substitu??ts", Catalan},
		{"rozd??lit", Czech},
		{"tvo??en", Czech},
		{"subjekt??", Czech},
		{"nesufi??econ", Esperanto},
		{"intermiksi??is", Esperanto},
		{"mona??inoj", Esperanto},
		{"kreita??oj", Esperanto},
		{"??pinante", Esperanto},
		{"apena??", Esperanto},
		{"gro??", German},
		{"????????????", Greek},
		{"fekv??", Hungarian},
		{"meggy??r??zni", Hungarian},
		{"????????????????????????", Japanese},
		{"????????", Kazakh},
		{"??????????????????????", Kazakh},
		{"????????", Kazakh},
		{"????????", Kazakh},
		{"??????????????", Kazakh},
		{"teolo??iska", Latvian},
		{"bla??ene", Latvian},
		{"ce??ojumiem", Latvian},
		{"numuri??u", Latvian},
		{"mergel??s", Lithuanian},
		{"??rengus", Lithuanian},
		{"slegiam??", Lithuanian},
		{"??????????????", Macedonian},
		{"????????????", Macedonian},
		{"??????????", Macedonian},
		{"??????????????", Macedonian},
		{"???????????????", Marathi},
		{"????????????", Mongolian},
		{"??????????", Mongolian},
		{"zmieni??y", Polish},
		{"pa??stwowych", Polish},
		{"mniejszo??ci", Polish},
		{"gro??ne", Polish},
		{"ialomi??a", Romanian},
		{"??????????????????????", Serbian},
		{"????????????????????????????", Serbian},
		{"pod??a", Slovak},
		{"poh??ade", Slovak},
		{"m??tvych", Slovak},
		{"????????????????????", Ukrainian},
		{"????????????????", Ukrainian},
		{"????????????????", Ukrainian},
		{"c???m", Vietnamese},
		{"th???n", Vietnamese},
		{"ch???ng", Vietnamese},
		{"qu???y", Vietnamese},
		{"s???n", Vietnamese},
		{"nh???n", Vietnamese},
		{"d???t", Vietnamese},
		{"ch???t", Vietnamese},
		{"?????p", Vietnamese},
		{"m???n", Vietnamese},
		{"h???u", Vietnamese},
		{"hi???n", Vietnamese},
		{"l???n", Vietnamese},
		{"bi???u", Vietnamese},
		{"k???m", Vietnamese},
		{"di???m", Vietnamese},
		{"ph???", Vietnamese},
		{"vi???c", Vietnamese},
		{"ch???nh", Vietnamese},
		{"tr??", Vietnamese},
		{"rav???", Vietnamese},
		{"th??", Vietnamese},
		{"ngu???n", Vietnamese},
		{"th???", Vietnamese},
		{"s???i", Vietnamese},
		{"t???ng", Vietnamese},
		{"nh???", Vietnamese},
		{"m???i", Vietnamese},
		{"b???i", Vietnamese},
		{"t???t", Vietnamese},
		{"gi???i", Vietnamese},
		{"m???t", Vietnamese},
		{"h???p", Vietnamese},
		{"h??ng", Vietnamese},
		{"t???ng", Vietnamese},
		{"c???a", Vietnamese},
		{"s???", Vietnamese},
		{"c??ng", Vietnamese},
		{"nh???ng", Vietnamese},
		{"ch???c", Vietnamese},
		{"d???ng", Vietnamese},
		{"th???c", Vietnamese},
		{"k???", Vietnamese},
		{"k???", Vietnamese},
		{"m???", Vietnamese},
		{"m???", Vietnamese},
		{"a???iw??r??", Yoruba},
		{"???aaju", Yoruba},
		{"????????????????", Unknown},
		{"??????????????????????????", Unknown},
		{"house", Unknown},

		// words with unique alphabet
		{"????????????", Armenian},
		{"??????????????????", Bengali},
		{"????????????????????????", Georgian},
		{"??????????????????", Greek},
		{"????????????????????????", Gujarati},
		{"????????????????", Hebrew},
		{"??????", Japanese},
		{"???????????????", Korean},
		{"?????????????????????????????????", Punjabi},
		{"??????????????????????????????", Tamil},
		{"???????????????????????????????????????", Telugu},
		{"????????????????????????????????????????????????", Thai},
	}
	for _, testCase := range testCases {
		detectedLanguage := detectorForAllLanguages.detectLanguageWithRules([]string{testCase.word})
		message := fmt.Sprintf(
			"expected %v for word '%s', got %v",
			testCase.expectedLanguage,
			testCase.word,
			detectedLanguage,
		)
		assert.Equal(t, testCase.expectedLanguage, detectedLanguage, message)
	}
}

func TestFilterLanguagesByRules(t *testing.T) {
	testCases := []struct {
		word              string
		expectedLanguages []Language
	}{
		{"????????????????", []Language{Arabic, Persian, Urdu}},
		{"??????????????????????????", []Language{
			Belarusian, Bulgarian, Kazakh, Macedonian, Mongolian, Russian, Serbian, Ukrainian},
		},
		{"??????????????????", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"????????", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"??????????", []Language{Belarusian, Kazakh, Mongolian, Russian}},
		{"??????????????", []Language{Bulgarian, Kazakh, Mongolian, Russian}},
		{"????????????????", []Language{Bulgarian, Kazakh, Mongolian, Russian}},
		{"??????????????", []Language{Belarusian, Kazakh, Ukrainian}},
		{"??????????????????????", []Language{Macedonian, Serbian}},
		{"??????????????????????????", []Language{Macedonian, Serbian}},
		{"????????????????????", []Language{Macedonian, Serbian}},
		{"aizkl??t??", []Language{Latvian, Maori, Yoruba}},
		{"sist??mas", []Language{Latvian, Maori, Yoruba}},
		{"pal??dzi", []Language{Latvian, Maori, Yoruba}},
		{"nh???n", []Language{Vietnamese, Yoruba}},
		{"ch???n", []Language{Vietnamese, Yoruba}},
		{"prihva??anju", []Language{Bosnian, Croatian, Polish}},
		{"na??ete", []Language{Bosnian, Croatian, Vietnamese}},
		{"vis??o", []Language{Portuguese, Vietnamese}},
		{"wyst??pi??", []Language{Lithuanian, Polish}},
		{"budow??", []Language{Lithuanian, Polish}},
		{"neb??sime", []Language{Latvian, Lithuanian, Maori, Yoruba}},
		{"afi??ate", []Language{Azerbaijani, Romanian, Turkish}},
		{"kradzie??ami", []Language{Polish, Romanian}},
		{"??nviat", []Language{French, Romanian}},
		{"venerd??", []Language{Italian, Vietnamese, Yoruba}},
		{"a??os", []Language{Basque, Spanish}},
		{"rozoh??uje", []Language{Czech, Slovak}},
		{"rtu??", []Language{Czech, Slovak}},
		{"preg??tire", []Language{Romanian, Vietnamese}},
		{"je??te", []Language{Czech, Romanian, Slovak}},
		{"minjaver??ir", []Language{Icelandic, Turkish}},
		{"??agnarskyldu", []Language{Icelandic, Turkish}},
		{"neb??tu", []Language{French, Hungarian}},
		{"hashemid??ve", []Language{Afrikaans, Albanian, Dutch, French}},
		{"for??t", []Language{Afrikaans, French, Portuguese, Vietnamese}},
		{"succ??dent", []Language{French, Italian, Vietnamese, Yoruba}},
		{"o??", []Language{French, Italian, Vietnamese, Yoruba}},
		{"t??eliseks", []Language{Estonian, Hungarian, Portuguese, Vietnamese}},
		{"vi??iem", []Language{Catalan, Italian, Vietnamese, Yoruba}},
		{"contr??le", []Language{French, Portuguese, Slovak, Vietnamese}},
		{"direkt??r", []Language{Bokmal, Danish, Nynorsk}},
		{"v??voj", []Language{Czech, Icelandic, Slovak, Turkish, Vietnamese}},
		{"p??ralt", []Language{Estonian, Finnish, German, Slovak, Swedish}},
		{"lab??k", []Language{French, Portuguese, Romanian, Turkish, Vietnamese}},
		{"pr??ctiques", []Language{Catalan, French, Italian, Portuguese, Vietnamese}},
		{"??berrascht", []Language{
			Azerbaijani, Catalan, Estonian, German, Hungarian, Spanish, Turkish},
		},
		{"indeb??rer", []Language{Bokmal, Danish, Icelandic, Nynorsk}},
		{"m??ned", []Language{Bokmal, Danish, Nynorsk, Swedish}},
		{"zaru??en", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"zkou??kou", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"navr??en", []Language{Bosnian, Czech, Croatian, Latvian, Lithuanian, Slovak, Slovene}},
		{"fa??onnage", []Language{
			Albanian, Azerbaijani, Basque, Catalan, French, Portuguese, Turkish},
		},
		{"h??her", []Language{
			Azerbaijani, Estonian, Finnish, German, Hungarian, Icelandic, Swedish, Turkish},
		},
		{"catedr??ticos", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"pol??tica", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"m??sica", []Language{
			Catalan, Czech, Icelandic, Irish, Hungarian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"contradicci??", []Language{
			Catalan, Hungarian, Icelandic, Irish, Polish, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"nom??s", []Language{
			Catalan, Czech, French, Hungarian, Icelandic, Irish,
			Italian, Portuguese, Slovak, Spanish, Vietnamese, Yoruba},
		},
		{"house", []Language{
			Afrikaans, Albanian, Azerbaijani, Basque, Bokmal, Bosnian, Catalan, Croatian, Czech, Danish, Dutch, English,
			Esperanto, Estonian, Finnish, French, Ganda, German, Hungarian, Icelandic, Indonesian, Irish, Italian,
			Latin, Latvian, Lithuanian, Malay, Maori, Nynorsk, Polish, Portuguese, Romanian, Shona, Slovak, Slovene,
			Somali, Sotho, Spanish, Swahili, Swedish, Tagalog, Tsonga, Tswana, Turkish, Vietnamese, Welsh, Xhosa,
			Yoruba, Zulu},
		},
	}
	for _, testCase := range testCases {
		filteredLanguages := detectorForAllLanguages.filterLanguagesByRules([]string{testCase.word})
		message := fmt.Sprintf("expected %v for word '%s', got %v", testCase.expectedLanguages, testCase.word, filteredLanguages)
		assert.ElementsMatch(t, testCase.expectedLanguages, filteredLanguages, message)
	}
}

func TestDetectionOfInvalidStrings(t *testing.T) {
	testCases := []string{"", " \n  \t;", "3<856%)??"}
	for _, testCase := range testCases {
		language, exists := detectorForAllLanguages.DetectLanguageOf(testCase)
		assert.Equal(t, Unknown, language)
		assert.False(t, exists)
	}
}

func TestLanguageDetectionIsDeterministic(t *testing.T) {
	testCases := []struct {
		text      string
		languages []Language
	}{
		{
			"???? ???? ???? ???????? ?????????? ???????? ???????????????????? i vote for bts ( _ ) as the _ via ( _ )",
			[]Language{English, Urdu},
		},
		{
			"Az elm??lt h??tv??g??n 12-re emelkedett az elhunyt koronav??rus-fert??z??ttek sz??ma Szlov??ki??ban. Mindegyik szoci??lis otthon dolgoz??j??t letesztelik, Matovi?? szerint az ing??z??knak m??g v??rniuk kellene a tesztel??ssel",
			[]Language{Hungarian, Slovak},
		},
	}
	for _, testCase := range testCases {
		detector := NewLanguageDetectorBuilder().
			FromLanguages(testCase.languages...).
			WithPreloadedLanguageModels().
			Build()
		detectedLanguages := make(map[Language]bool)
		for i := 0; i < 100; i++ {
			language, _ := detector.DetectLanguageOf(testCase.text)
			detectedLanguages[language] = true
		}
		assert.Len(
			t,
			detectedLanguages,
			1,
			fmt.Sprintf("language detector is non-deterministic for languages %v", testCase.languages),
		)
	}
}
