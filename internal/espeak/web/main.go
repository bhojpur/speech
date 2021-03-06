package main

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
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bhojpur/speech/pkg/espeak"
)

var (
	port         string
	audioDir     string
	downloadsDir string
)

func main() {
	log.Println("Bhojpur Speech text-to-speech online")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	flag.Parse()

	os.MkdirAll(audioDir, 0777)
	os.MkdirAll(downloadsDir, 0777)
	http.HandleFunc("/", serve)
	http.HandleFunc("/download", download)
	http.Handle("/downloads/", http.StripPrefix("/downloads/", http.FileServer(http.Dir(downloadsDir))))
	http.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(audioDir))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	go garbageCollect()

	fmt.Println("server engine listening on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	flag.StringVar(&port, "port", "8080", "port to listen at")
	flag.StringVar(&port, "p", "8080", "port to listen at")
	flag.StringVar(&audioDir, "audio-dir", "static/audio", "dir to save the audio files at. will be created if needed")
	flag.StringVar(&audioDir, "a", "static/audio", "dir to save the audio files at. will be created if needed")
	flag.StringVar(&downloadsDir, "downloads-dir", "static/downloads", "dir to save the downloadable audio files at. will be created if needed")
	flag.StringVar(&downloadsDir, "d", "static/downloads", "dir to save the downloadable audio files at. will be created if needed")
}

// garbageCollect cleans audio files from audioDir and downloadsDir
// periodically
func garbageCollect() {
	var err error
	tick := time.Tick(500 * time.Millisecond)
	downloadsTick := time.Tick(5 * time.Second)
	for {
		select {
		case <-tick:
			halfASecondAgo := time.Now().Add(-500 * time.Millisecond)
			err = filepath.Walk(audioDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Println(err)
				}
				if info.Name() == filepath.Base(audioDir) && info.IsDir() {
					return nil
				}
				if info.ModTime().Before(halfASecondAgo) {
					p := filepath.Join(audioDir, info.Name())
					if err := os.Remove(p); err != nil {
						log.Println(err)
					}
					log.Printf("removing %s", p)
				}
				return nil
			})
			if err != nil {
				log.Println(err)
			}
		case <-downloadsTick:
			fifteenSecondsAgo := time.Now().Add(-15 * time.Second)
			err = filepath.Walk(downloadsDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Println(err)
				}
				if info.Name() == filepath.Base(downloadsDir) && info.IsDir() {
					return nil
				}
				if info.ModTime().Before(fifteenSecondsAgo) {
					p := filepath.Join(downloadsDir, info.Name())
					if err := os.Remove(p); err != nil {
						log.Println(err)
					}
					log.Printf("removing %s", p)
				}
				return nil
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

type data struct {
	Error               string
	VoiceName           string `json:"voice"`
	Say                 string `json:"say"`
	Rate                int    `json:"rate"`
	Volume              int    `json:"volume"`
	Pitch               int    `json:"pitch"`
	Range               int    `json:"range"`
	AnnouncePunctuation string `json:"punctuation"`
	AnnounceCapitals    string `json:"capitals"`
	WordGap             int    `json:"word-gap"`
	FileSource          string
	PunctList           string
}

func serve(w http.ResponseWriter, r *http.Request) {

	var params *espeak.Parameters
	var voice *espeak.Voice
	var name string
	var err error
	params, err = getParams(r, audioDir)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	voice, name = getVoice(r)

	d := getData(params, name, r.FormValue("say"))
	if d.Say != "" {
		src := randString(64) + ".wav"
		d.FileSource = "/audio/" + src
		_, err = espeak.TextToSpeech(d.Say, voice, src, params)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := tpl.Execute(w, d); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func download(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	var params *espeak.Parameters
	var voice *espeak.Voice
	var name string
	var err error
	params, err = getParams(r, downloadsDir)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	voice, name = getVoice(r)

	d := getData(params, name, r.FormValue("say"))
	if d.Say != "" {
		src := randString(64) + ".wav"
		d.FileSource = "/downloads/" + src
		_, err = espeak.TextToSpeech(d.Say, voice, src, params)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(d.FileSource)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write(b)
	}
}

func getData(params *espeak.Parameters, vname, say string) *data {
	var d = &data{
		VoiceName: vname,
		Rate:      params.Rate,
		Volume:    params.Volume,
		Pitch:     params.Pitch,
		Range:     params.Range,
		WordGap:   params.WordGap,
		PunctList: params.PunctuationList(),
	}
	switch params.AnnounceCapitals {
	case espeak.CapitalPitchRaise:
		d.AnnounceCapitals = "pitch-raise"
	case espeak.CapitalSoundIcon:
		d.AnnounceCapitals = "sound-icon"
	case espeak.CapitalSpelling:
		d.AnnounceCapitals = "spelling"
	case espeak.CapitalNone:
		d.AnnounceCapitals = "none"
	}
	switch params.AnnouncePunctuation {
	case espeak.PunctNone:
		d.AnnouncePunctuation = "none"
	case espeak.PunctSome:
		d.AnnouncePunctuation = "some"
	case espeak.PunctAll:
		d.AnnouncePunctuation = "all"
	}

	if say != "" {
		d.Say = say
	}
	return d
}

func getVoice(r *http.Request) (*espeak.Voice, string) {
	r.ParseForm()
	switch r.PostFormValue("voice") {
	case "en-us-male":
		return espeak.ENUSMale, "en-us-male"
	case "es-es-male":
		return espeak.ESSpainMale, "es-es-male"
	case "es-lat-male":
		return espeak.ESLatinMale, "es-lat-male"
	case "fr-fr-male":
		return espeak.FRFranceMale, "fr-fr-male"
	default:
		return espeak.DefaultVoice, "en-us-male"
	}
}

func getParams(r *http.Request, dir string) (*espeak.Parameters, error) {

	var n int
	var err error
	var params = espeak.NewParameters(espeak.WithDir(dir))
	r.ParseForm()
	if rate := r.PostFormValue("rate"); rate != "" {
		n, err = strconv.Atoi(rate)
		if err != nil {
			return nil, err
		}
		params.WithRate(n)
	} else {
		params.WithRate(espeak.DefaultParameters.Rate)
	}
	if vol := r.PostFormValue("volume"); vol != "" {
		n, err = strconv.Atoi(vol)
		if err != nil {
			return nil, err
		}
		params.WithVolume(n)
	} else {
		params.WithVolume(espeak.DefaultParameters.Volume)
	}

	if pitch := r.PostFormValue("pitch"); pitch != "" {
		n, err = strconv.Atoi(pitch)
		if err != nil {
			return nil, err
		}
		params.WithPitch(n)
	} else {
		params.WithPitch(espeak.DefaultParameters.Pitch)
	}

	if rng := r.PostFormValue("range"); rng != "" {
		n, err = strconv.Atoi(rng)
		if err != nil {
			return nil, err
		}
		params.WithRange(n)
	} else {
		params.WithRange(espeak.DefaultParameters.Range)
	}

	if wordGap := r.PostFormValue("word-gap"); wordGap != "" {
		n, err = strconv.Atoi(wordGap)
		if err != nil {
			return nil, err
		}
		params.WithWordGap(n)
	} else {
		params.WithWordGap(espeak.DefaultParameters.WordGap)
	}

	params.SetPunctuationList(r.PostFormValue("punctuation-list"))

	switch r.PostFormValue("punctuation") {
	case "all":
		params.WithAnnouncePunctuation(espeak.PunctAll)
	case "some":
		params.WithAnnouncePunctuation(espeak.PunctSome)
	default:
		params.WithAnnouncePunctuation(espeak.PunctNone)
	}

	switch r.PostFormValue("capitals") {
	case "sound-icon":
		params.WithAnnounceCapitals(espeak.CapitalSoundIcon)
	case "spelling":
		params.WithAnnounceCapitals(espeak.CapitalSpelling)
	case "pitch-raise":
		params.WithAnnounceCapitals(espeak.CapitalPitchRaise)
	default:
		params.WithAnnounceCapitals(espeak.CapitalNone)
	}

	return params, nil
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) (out string) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		out += string(chars[rand.Intn(len(chars))])
	}
	return
}
