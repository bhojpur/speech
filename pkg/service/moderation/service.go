package moderation

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
	"net/url"
	"strings"
)

type Message struct {
	From string
	Text string
}

type Filter interface {
	Moderate(text Message) string
	SetFilterMap(filterMap FilterMap)
}

type BaseFilterDecorator struct {
	filter Filter
}

func (f *BaseFilterDecorator) SetFilterMap(filterMap FilterMap) {
	f.filter.SetFilterMap(filterMap)
}

func (f *BaseFilterDecorator) Moderate(text Message) string {
	return f.filter.Moderate(text)
}

var mostPopularTLD = []string{
	".com",
	".net",
	".org",
	".de",
	".icu",
	".uk",
	".ru",
	".info",
	".top",
	".xyz",
	".tk",
	".cn",
	".ga",
	".cf",
	".nl",
}

func NewDefaultFilter(moderationPair, ignoreString string, users []string) *UserFilterDecorator {
	filterDefault := NewFilterDefault(moderationPair, ignoreString)
	urlFilter := NewUrlFilterDecorator(filterDefault)
	return NewUserFilterDecorator(urlFilter, users)
}

type UserFilterDecorator struct {
	BaseFilterDecorator
	users map[string]struct{}
}

func NewUserFilterDecorator(filter Filter, users []string) *UserFilterDecorator {
	u := map[string]struct{}{}
	for _, user := range users {
		if len(user) > 3 {
			u[user] = struct{}{}
		}
	}
	decorator := UserFilterDecorator{
		BaseFilterDecorator: BaseFilterDecorator{
			filter: filter,
		},
		users: u,
	}
	return &decorator
}

func (f *UserFilterDecorator) SetFilterMap(filterMap FilterMap) {
	f.filter.SetFilterMap(filterMap)
}

func (f *UserFilterDecorator) Moderate(msg Message) string {
	_, ok := f.users[msg.From]
	if ok {
		return ""
	}
	return f.filter.Moderate(msg)
}

type UrlFilterDecorator struct {
	BaseFilterDecorator
}

func NewUrlFilterDecorator(filter Filter) *UrlFilterDecorator {
	return &UrlFilterDecorator{
		BaseFilterDecorator{
			filter: filter,
		},
	}
}

func (f *UrlFilterDecorator) SetFilterMap(filterMap FilterMap) {
	f.filter.SetFilterMap(filterMap)
}

func (f *UrlFilterDecorator) Moderate(msg Message) string {
	split := strings.Split(msg.Text, " ")
	var result string
	for _, word := range split {
		if !f.isValidUrl(word) && !f.isContainsTopLevelDomain(word) {
			result += word + " "
		}
	}
	msg.Text = result
	return f.filter.Moderate(msg)
}

func (f *UrlFilterDecorator) isValidUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (f *UrlFilterDecorator) isContainsTopLevelDomain(str string) bool {
	for _, item := range mostPopularTLD {
		if strings.Contains(str, item) {
			return true
		}
	}
	return false
}

type FilterDefault struct {
	filterMap FilterMap
}

const defaultIgnore = "shit,slut,spunk,whore,fuck,nigger,sex,pussy,queer,sh1t,wank,wtf,anal,bitch,poop,tosser,vagina,balls,Goddamn,muff,clitoris,knobend,knob end,ballsack,bastard,bum,penis,arse,dick,f u c k,God damn,pube,anus,cunt,fellate,feck,felching,lmao,nigga,omg,bollok,dildo,fag,homo,turd,bugger,buttplug,dyke,bollock,flange,blowjob,boob,crap,labia,scrotum,s hit,smegma,ass,biatch,coon,lmfao,boner,fudge packer,jizz,hell,jerk,piss,tit,twat,bloody,butt,damn,blow job,cock,fellatio,fudgepacker,prick"

func NewFilterDefault(moderationPair, ignoreString string) *FilterDefault {
	f := new(FilterDefault)
	if len(moderationPair) > 0 || len(ignoreString) > 0 {
		builder := FilterMapBuilderImpl{}
		f.filterMap = *builder.Build(moderationPair, ignoreString+defaultIgnore)
	} else {
		f.filterMap = DefaultFilterMap
	}
	return f
}

func (f *FilterDefault) SetFilterMap(filterMap FilterMap) {
	f.filterMap = filterMap
}

func (f *FilterDefault) Moderate(msg Message) string {
	words := strings.Split(msg.Text, " ")
	var result string
	var val string
	var ok bool

	for _, word := range words {
		val, ok = f.filterMap.Get(strings.ToLower(word))
		if ok {
			result += val
		} else {
			result += word
		}
		result += " "
	}

	return result
}
