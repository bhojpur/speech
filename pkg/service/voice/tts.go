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

import (
	"strings"

	lang_detection "github.com/bhojpur/speech/pkg/service/lang-detection"
	"github.com/bhojpur/speech/pkg/service/moderation"
	"github.com/bhojpur/speech/pkg/synthesis"
	"github.com/bhojpur/speech/pkg/utils"
	"github.com/bhojpur/speech/pkg/utils/repo"
	"github.com/bhojpur/speech/pkg/voices"
)

type SetterFromService interface {
	SetFrom(from string)
}

type SpeechVoiceService interface {
	SetterFromService
	Speak(text string) error
}

type GoTtsService struct {
	speech              synthesis.Speech
	language            string
	filter              moderation.Filter
	volume              float64
	from                string
	repo                repo.SettingRepo
	langDetector        lang_detection.LanguageDetectionService
	langDetectorEnabled bool
}

func NewGoTtsService(language string, filter moderation.Filter, volume float64, repo repo.SettingRepo, langDetectorEnabled bool, langDetector lang_detection.LanguageDetectionService) *GoTtsService {
	s := new(GoTtsService)
	if len(language) == 0 {
		language = voices.English
	}
	if volume < 0 || volume > 15 {
		return nil
	}

	s.volume = volume
	s.speech = NewSpeech(language, volume)
	s.language = language
	s.filter = filter
	s.repo = repo
	s.langDetector = langDetector
	s.langDetectorEnabled = langDetectorEnabled
	return s
}

func (s *GoTtsService) SetFrom(from string) {
	s.from = from
}

func (s *GoTtsService) setLanguage(language string) {
	s.speech = NewSpeech(language, s.volume)
	s.language = language
}

func (s *GoTtsService) detectLanguage(text string) error {
	langDetected, err := s.langDetector.Detect(text)
	if err != nil {
		return err
	}
	lang := string(*langDetected)
	s.setLanguage(lang)
	if s.repo != nil {
		settingDb, err := s.repo.GetSettings()
		if err != nil {
			return err
		}
		s.langDetectorEnabled = settingDb.LanguageDetectorEnabled
		settingDb.Language = lang
		return s.repo.SaveSettings(settingDb)
	}
	return nil
}

func (s *GoTtsService) Speak(text string) error {
	if s.repo != nil {
		settingDb, err := s.repo.GetSettings()
		if err != nil {
			return err
		}
		s.filter = moderation.NewDefaultFilter(settingDb.ReplacementWordPair, settingDb.IgnoreWords, utils.StrEnumerationToArray(settingDb.UserBanList))
		if s.volume != settingDb.Volume {
			s.speech = NewSpeech(s.language, s.volume)
		}
		s.volume = settingDb.Volume
		s.langDetectorEnabled = settingDb.LanguageDetectorEnabled
		if err = s.repo.SaveSettings(settingDb); err != nil {
			return err
		}
		s.setLanguage(settingDb.Language)
	}

	if s.langDetectorEnabled {
		if err := s.detectLanguage(text); err != nil {
			return err
		}
	}

	result := s.filter.Moderate(moderation.Message{From: s.from, Text: text})
	result = strings.Trim(result, " ")
	fromLen := len(s.from)
	if fromLen > len(result) && len(result) == 0 {
		return nil
	}
	check := result[fromLen:]
	if strings.HasSuffix(check, "say    !") {
		return nil
	}
	return s.speech.Speak(result)
}

func NewSpeech(language string, volume float64) synthesis.Speech {
	return synthesis.Speech{Folder: "audios", Language: language, Volume: volume, Speed: 1}
}
