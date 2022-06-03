package database

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

	"github.com/bhojpur/speech/pkg/service/moderation"
	"github.com/bhojpur/speech/pkg/utils"
)

type SettingDB struct {
	Id                      int
	ReplacementWordPair     string
	IgnoreWords             string
	Language                string
	LanguageDetectorEnabled bool
	UserBanList             string
	ChannelsToListen        string
	Volume                  float64
}

func (s *SettingDB) SetReplacementWordPair(filter moderation.FilterMap) {
	filterMap := filter.Range()
	result := ""
	for key, value := range filterMap {
		result += key + ":" + value + ","
	}
	s.ReplacementWordPair = result[:len(result)-1]
}

func (s *SettingDB) SetIgnoreWords(words []string) {
	var result string
	for _, item := range words {
		result += item + ","
	}
	s.IgnoreWords = result[:len(result)-1]
}

func (s *SettingDB) SetUserBanList(users []string) {
	var result string
	for _, item := range users {
		result += item + ","
	}
	s.UserBanList = result[:len(result)-1]
}

func (s *SettingDB) SetChannelsToListen(list []string) {
	var result string
	for _, item := range list {
		result += item + ","
	}
	s.ChannelsToListen = result[:len(result)-1]
}

type Setting struct {
	Id                      int
	ReplacementWordPair     moderation.FilterMap
	IgnoreWords             []string
	Language                string
	LanguageDetectorEnabled bool
	UserBanList             []string
	ChannelsToListen        []string
	Volume                  int
}

func (s *Setting) SetIgnoreWords(str string) {
	s.IgnoreWords = strings.Split(str, ",")
}

func (s *Setting) StoreIgnoreWord(word string) {
	s.IgnoreWords = utils.ArrayStore(s.IgnoreWords, word)
}

func (s *Setting) DeleteIgnoreWord(word string) {
	s.IgnoreWords = utils.ArrayDelete(s.IgnoreWords, word)
}

func (s *Setting) SetUserBanList(str string) {
	s.UserBanList = strings.Split(str, ",")
}

func (s *Setting) StoreUserBanList(user string) {
	s.UserBanList = utils.ArrayStore(s.UserBanList, user)
}

func (s *Setting) DeleteUserBanList(user string) {
	s.UserBanList = utils.ArrayDelete(s.UserBanList, user)
}

func (s *Setting) SetChannelsToListen(str string) {
	s.ChannelsToListen = strings.Split(str, ",")
}

func (s *Setting) StoreChannelsToListen(user string) {
	s.ChannelsToListen = utils.ArrayStore(s.ChannelsToListen, user)
}

func (s *Setting) DeleteChannelsToListen(user string) {
	s.ChannelsToListen = utils.ArrayDelete(s.ChannelsToListen, user)
}

func (s *Setting) SetReplacementPair(key, value string) {
	s.ReplacementWordPair.Set(key, value)
}

func (s *Setting) RemoveReplacementPair(key string) {
	s.ReplacementWordPair.Remove(key)
}

func (s *Setting) GetReplacementPair(key string) (string, bool) {
	return s.ReplacementWordPair.Get(key)
}

func (s *Setting) SetReplacementWordPair(str string) {
	builder := moderation.FilterMapBuilderImpl{}
	s.ReplacementWordPair = *builder.Build(str, "")
}
