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

const missingLanguageMessage = "LanguageDetector needs at least 2 languages to choose from"

// UnconfiguredLanguageDetectorBuilder is the interface describing the methods
// for specifying which languages will be used to build an instance of
// LanguageDetector. All methods return an implementation of the interface
// LanguageDetectorBuilder.
type UnconfiguredLanguageDetectorBuilder interface {
	// FromAllLanguages configures the LanguageDetectorBuilder
	// to use all built-in languages.
	FromAllLanguages() LanguageDetectorBuilder

	// FromAllSpokenLanguages configures the LanguageDetectorBuilder
	// to use all built-in spoken languages.
	FromAllSpokenLanguages() LanguageDetectorBuilder

	// FromAllLanguagesWithArabicScript configures the LanguageDetectorBuilder
	// to use all built-in languages supporting the Arabic script.
	FromAllLanguagesWithArabicScript() LanguageDetectorBuilder

	// FromAllLanguagesWithCyrillicScript configures the LanguageDetectorBuilder
	// to use all built-in languages supporting the Cyrillic script.
	FromAllLanguagesWithCyrillicScript() LanguageDetectorBuilder

	// FromAllLanguagesWithDevanagariScript configures the LanguageDetectorBuilder
	// to use all built-in languages supporting the Devanagari script.
	FromAllLanguagesWithDevanagariScript() LanguageDetectorBuilder

	// FromAllLanguagesWithLatinScript configures the LanguageDetectorBuilder
	// to use all built-in languages supporting the Latin script.
	FromAllLanguagesWithLatinScript() LanguageDetectorBuilder

	// FromAllLanguagesWithout configures the LanguageDetectorBuilder
	// to use all built-in languages except for those specified as arguments
	// passed to this method. Panics if less than two languages are used to
	// build the LanguageDetector.
	FromAllLanguagesWithout(languages ...Language) LanguageDetectorBuilder

	// FromLanguages configures the LanguageDetectorBuilder to use the languages
	// specified as arguments passed to this method. Panics if less than two
	// languages are specified.
	FromLanguages(languages ...Language) LanguageDetectorBuilder

	// FromIsoCodes639_1 configures the LanguageDetectorBuilder to use those
	// languages whose ISO 639-1 codes are specified as arguments passed to
	// this method. Panics if less than two iso codes are specified.
	FromIsoCodes639_1(isoCodes ...IsoCode639_1) LanguageDetectorBuilder

	// FromIsoCodes639_3 configures the LanguageDetectorBuilder to use those
	// languages whose ISO 639-3 codes are specified as arguments passed to
	// this method. Panics if less than two iso codes are specified.
	FromIsoCodes639_3(isoCodes ...IsoCode639_3) LanguageDetectorBuilder
}

// LanguageDetectorBuilder is the interface that defines any other settings
// to use for building an instance of LanguageDetector, except for the languages
// to use.
type LanguageDetectorBuilder interface {
	// WithMinimumRelativeDistance sets the desired value for the minimum
	// relative distance measure.
	//
	// By default, Detector returns the most likely language for a given
	// input text. However, there are certain words that are spelled the
	// same in more than one language. The word `prologue`, for instance,
	// is both a valid English and French word. Detector would output either
	// English or French which might be wrong in the given context.
	// For cases like that, it is possible to specify a minimum relative
	// distance that the logarithmized and summed up probabilities for
	// each possible language have to satisfy.
	//
	// Be aware that the distance between the language probabilities is
	// dependent on the length of the input text. The longer the input
	// text, the larger the distance between the languages. So if you
	// want to classify very short text phrases, do not set the minimum
	// relative distance too high. Otherwise you will get most results
	// returned as Unknown which is the return value for cases
	// where language detection is not reliably possible.
	//
	// Panics if distance is smaller than 0.0 or greater than 0.99.
	WithMinimumRelativeDistance(distance float64) LanguageDetectorBuilder

	// WithPreloadedLanguageModels configures LanguageDetectorBuilder to
	// preload all language models when creating the instance of LanguageDetector.
	//
	// By default, Detector uses lazy-loading to load only those language
	// models on demand which are considered relevant by the rule-based
	// filter engine. For web services, for instance, it is rather
	// beneficial to preload all language models into memory to avoid
	// unexpected latency while waiting for the service response. This
	// method allows to switch between these two loading modes.
	WithPreloadedLanguageModels() LanguageDetectorBuilder

	// Build creates and returns the configured instance of LanguageDetector.
	Build() LanguageDetector
	getLanguages() []Language
	getMinimumRelativeDistance() float64
}

type languageDetectorBuilder struct {
	languages                     []Language
	minimumRelativeDistance       float64
	isEveryLanguageModelPreloaded bool
}

// NewLanguageDetectorBuilder returns a new instance that implements the
// interface UnconfiguredLanguageDetectorBuilder.
func NewLanguageDetectorBuilder() UnconfiguredLanguageDetectorBuilder {
	return &languageDetectorBuilder{}
}

func (builder *languageDetectorBuilder) FromAllLanguages() LanguageDetectorBuilder {
	return builder.from(AllLanguages())
}

func (builder *languageDetectorBuilder) FromAllSpokenLanguages() LanguageDetectorBuilder {
	return builder.from(AllSpokenLanguages())
}

func (builder *languageDetectorBuilder) FromAllLanguagesWithArabicScript() LanguageDetectorBuilder {
	return builder.from(AllLanguagesWithArabicScript())
}

func (builder *languageDetectorBuilder) FromAllLanguagesWithCyrillicScript() LanguageDetectorBuilder {
	return builder.from(AllLanguagesWithCyrillicScript())
}

func (builder *languageDetectorBuilder) FromAllLanguagesWithDevanagariScript() LanguageDetectorBuilder {
	return builder.from(AllLanguagesWithDevanagariScript())
}

func (builder *languageDetectorBuilder) FromAllLanguagesWithLatinScript() LanguageDetectorBuilder {
	return builder.from(AllLanguagesWithLatinScript())
}

func (builder *languageDetectorBuilder) FromAllLanguagesWithout(languages ...Language) LanguageDetectorBuilder {
	languagesToLoad := AllLanguages()
	for _, languageToRemove := range languages {
		for i, currentLanguage := range languagesToLoad {
			if currentLanguage == languageToRemove {
				languagesToLoad = append(languagesToLoad[:i], languagesToLoad[i+1:]...)
				break
			}
		}
	}
	if len(languagesToLoad) < 2 {
		panic(missingLanguageMessage)
	}
	return builder.from(languagesToLoad)
}

func (builder *languageDetectorBuilder) FromLanguages(languages ...Language) LanguageDetectorBuilder {
	for i, language := range languages {
		if language == Unknown {
			languages = append(languages[:i], languages[i+1:]...)
			break
		}
	}
	if len(languages) < 2 {
		panic(missingLanguageMessage)
	}
	return builder.from(languages)
}

func (builder *languageDetectorBuilder) FromIsoCodes639_1(isoCodes ...IsoCode639_1) LanguageDetectorBuilder {
	for i, isoCode := range isoCodes {
		if isoCode == UnknownIsoCode639_1 {
			isoCodes = append(isoCodes[:i], isoCodes[i+1:]...)
			break
		}
	}
	if len(isoCodes) < 2 {
		panic(missingLanguageMessage)
	}
	var languages []Language
	for _, isoCode := range isoCodes {
		languages = append(languages, GetLanguageFromIsoCode639_1(isoCode))
	}
	return builder.from(languages)
}

func (builder *languageDetectorBuilder) FromIsoCodes639_3(isoCodes ...IsoCode639_3) LanguageDetectorBuilder {
	for i, isoCode := range isoCodes {
		if isoCode == UnknownIsoCode639_3 {
			isoCodes = append(isoCodes[:i], isoCodes[i+1:]...)
			break
		}
	}
	if len(isoCodes) < 2 {
		panic(missingLanguageMessage)
	}
	var languages []Language
	for _, isoCode := range isoCodes {
		languages = append(languages, GetLanguageFromIsoCode639_3(isoCode))
	}
	return builder.from(languages)
}

func (builder *languageDetectorBuilder) WithMinimumRelativeDistance(distance float64) LanguageDetectorBuilder {
	if distance < 0.0 || distance > 0.99 {
		panic("Minimum relative distance must lie in between 0.0 and 0.99")
	}
	builder.minimumRelativeDistance = distance
	return builder
}

func (builder *languageDetectorBuilder) WithPreloadedLanguageModels() LanguageDetectorBuilder {
	builder.isEveryLanguageModelPreloaded = true
	return builder
}

func (builder *languageDetectorBuilder) Build() LanguageDetector {
	return newLanguageDetector(
		builder.languages,
		builder.minimumRelativeDistance,
		builder.isEveryLanguageModelPreloaded,
	)
}

func (builder *languageDetectorBuilder) getLanguages() []Language {
	return builder.languages
}

func (builder *languageDetectorBuilder) getMinimumRelativeDistance() float64 {
	return builder.minimumRelativeDistance
}

func (builder *languageDetectorBuilder) from(languages []Language) LanguageDetectorBuilder {
	builder.languages = removeDuplicateLanguages(languages)
	builder.minimumRelativeDistance = 0.0
	builder.isEveryLanguageModelPreloaded = false
	return builder
}

func removeDuplicateLanguages(languages []Language) []Language {
	languageSet := make(map[Language]struct{})
	for _, language := range languages {
		_, exists := languageSet[language]
		if !exists {
			languageSet[language] = struct{}{}
		}
	}
	languageKeys := make([]Language, len(languageSet))
	i := 0
	for k := range languageSet {
		languageKeys[i] = k
		i++
	}
	return languageKeys
}
