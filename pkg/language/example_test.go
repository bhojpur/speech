package language_test

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

	spoken "github.com/bhojpur/speech/pkg/language"
)

func Example_basic() {
	languages := []spoken.Language{
		spoken.English,
		spoken.French,
		spoken.German,
		spoken.Spanish,
	}

	detector := spoken.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	if language, exists := detector.DetectLanguageOf("languages are awesome"); exists {
		fmt.Println(language)
	}

	// Output: English
}

// By default, Detector returns the most likely language for a given input text.
// However, there are certain words that are spelled the same in more than one
// language. The word `prologue`, for instance, is both a valid English and
// French word. Detector would output either English or French which might be
// wrong in the given context. For cases like that, it is possible to specify a
// minimum relative distance that the logarithmized and summed up probabilities
// for each possible language have to satisfy. It can be stated as seen below.
//
// Be aware that the distance between the language probabilities is dependent on
// the length of the input text. The longer the input text, the larger the
// distance between the languages. So if you want to classify very short text
// phrases, do not set the minimum relative distance too high. Otherwise Unknown
// will be returned most of the time as in the example below. This is the return
// value for cases where language detection is not reliably possible.
func Example_minimumRelativeDistance() {
	languages := []spoken.Language{
		spoken.English,
		spoken.French,
		spoken.German,
		spoken.Spanish,
	}

	detector := spoken.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithMinimumRelativeDistance(0.25).
		Build()

	language, exists := detector.DetectLanguageOf("languages are awesome")

	fmt.Println(language)
	fmt.Println(exists)

	// Output:
	// Unknown
	// false
}

// Knowing about the most likely language is nice but how reliable is the
// computed likelihood? And how less likely are the other examined languages in
// comparison to the most likely one? In the example below, a slice of
// ConfidenceValue is returned containing all possible languages sorted by their
// confidence value in descending order. The values that this method computes are
// part of a relative confidence metric, not of an absolute one. Each value is a
// number between 0.0 and 1.0. The most likely language is always returned with
// value 1.0. All other languages get values assigned which are lower than 1.0,
// denoting how less likely those languages are in comparison to the most likely
// language.
//
// The slice returned by this method does not necessarily contain all
// languages which the calling instance of LanguageDetector was built from.
// If the rule-based engine decides that a specific language is truly
// impossible, then it will not be part of the returned slice. Likewise,
// if no ngram probabilities can be found within the detector's languages
// for the given input text, the returned slice will be empty.
// The confidence value for each language not being part of the returned
// slice is assumed to be 0.0.
func Example_confidenceValues() {
	languages := []spoken.Language{
		spoken.English,
		spoken.French,
		spoken.German,
		spoken.Spanish,
	}

	detector := spoken.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	confidenceValues := detector.ComputeLanguageConfidenceValues("languages are awesome")

	for _, elem := range confidenceValues {
		fmt.Printf("%s: %.2f\n", elem.Language(), elem.Value())
	}

	// Output:
	// English: 1.00
	// French: 0.79
	// German: 0.75
	// Spanish: 0.72
}

// By default, Detector uses lazy-loading to load only those language models on
// demand which are considered relevant by the rule-based filter engine. For web
// services, for instance, it is rather beneficial to preload all language models
// into memory to avoid unexpected latency while waiting for the service response.
// If you want to enable the eager-loading mode, you can do it as seen below.
// Multiple instances of LanguageDetector share the same language models in
// memory which are accessed asynchronously by the instances.
func Example_eagerLoading() {
	spoken.NewLanguageDetectorBuilder().
		FromAllLanguages().
		WithPreloadedLanguageModels().
		Build()
}

// There might be classification tasks where you know beforehand that your language
// data is definitely not written in Latin, for instance. The detection accuracy
// can become better in such cases if you exclude certain languages from the
// decision process or just explicitly include relevant languages.
func Example_builderApi() {
	// Including all languages available in the library
	// consumes at least 2GB of memory and might
	// lead to slow runtime performance.
	spoken.NewLanguageDetectorBuilder().FromAllLanguages()

	// Include only languages that are not yet extinct
	// (= currently excludes Latin).
	spoken.NewLanguageDetectorBuilder().FromAllSpokenLanguages()

	// Include only languages written with Cyrillic script.
	spoken.NewLanguageDetectorBuilder().FromAllLanguagesWithCyrillicScript()

	// Exclude only the Spanish language from the decision algorithm.
	spoken.NewLanguageDetectorBuilder().FromAllLanguagesWithout(spoken.Spanish)

	// Only decide between English and German.
	spoken.NewLanguageDetectorBuilder().FromLanguages(spoken.English, spoken.German)

	// Select languages by ISO 639-1 code.
	spoken.NewLanguageDetectorBuilder().FromIsoCodes639_1(spoken.EN, spoken.DE)

	// Select languages by ISO 639-3 code.
	spoken.NewLanguageDetectorBuilder().FromIsoCodes639_3(spoken.ENG, spoken.DEU)
}
