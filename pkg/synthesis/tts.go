package synthesis

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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/bhojpur/speech/pkg/utils/handlers"
)

/**
 *
 * Use:
 *
 * speech := tts.Speech{Folder: "audios", Language: "en", Volume: 0, Speed: 1}
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Handler  handlers.PlayerInterface
	Volume   float64
	Speed    float64
}

// Creates a speech file with a given name
func (speech *Speech) CreateSpeechFile(text string, fileName string) (string, error) {
	var err error

	f := speech.Folder + "/" + fileName + ".mp3"
	if err = speech.createFolderIfNotExists(speech.Folder); err != nil {
		return "", err
	}

	if err = speech.downloadIfNotExists(f, text); err != nil {
		return "", err
	}

	return f, nil
}

// Plays an existent .mp3 file
func (speech *Speech) PlaySpeechFile(fileName string) error {
	if speech.Handler == nil {
		speech.Handler = &handlers.BeepPlayer{Volume: speech.Volume, Speed: 1}
	}
	err := speech.Handler.Play(fileName)
	if err != nil {
		return err
	}
	speech.deleteFile(fileName)
	return nil
}

func (speech *Speech) deleteFile(fileName string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(path + "\\" + fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func (speech *Speech) Speak(text string) error {
	var err error
	generatedHashName := speech.generateHashName(text)

	fileName, err := speech.CreateSpeechFile(text, generatedHashName)
	if err != nil {
		return err
	}

	return speech.PlaySpeechFile(fileName)
}

func (speech *Speech) createFolderIfNotExists(folder string) error {
	dir, err := os.Open(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0700)
	}

	dir.Close()
	return nil
}

/**
 * Download the voice file if does not exists.
 */
func (speech *Speech) downloadIfNotExists(fileName string, text string) error {
	f, err := os.Open(fileName)
	if err != nil {
		urlString := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
		fmt.Println(urlString)
		response, err := http.Get(urlString)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		output, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(output, response.Body)
		return err
	}

	defer f.Close()
	return nil
}

func (speech *Speech) generateHashName(name string) string {
	hash := md5.Sum([]byte(name))
	return fmt.Sprintf("%s_%s", speech.Language, hex.EncodeToString(hash[:]))
}
