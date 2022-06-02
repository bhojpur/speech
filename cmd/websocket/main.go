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
	"io"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const Host = "localhost"
const Port = "2700"
const buffsize = 8000

type Message struct {
	Result []struct {
		Conf  float64
		End   float64
		Start float64
		Word  string
	}
	Text string
}

var m Message

func main() {

	if len(os.Args) < 2 {
		panic("Please specify second argument")
	}

	u := url.URL{Scheme: "ws", Host: Host + ":" + Port, Path: ""}
	log.Info("connecting to ", u.String())

	// Opening websocket connection
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	check(err)
	defer c.Close()

	f, err := os.Open(os.Args[1])
	check(err)

	for {
		buf := make([]byte, buffsize)
		dat, err := f.Read(buf)

		if dat == 0 && err == io.EOF {
			err = c.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
			check(err)
			break
		}
		check(err)

		err = c.WriteMessage(websocket.BinaryMessage, buf)
		check(err)

		// Read message from server
		_, _, err = c.ReadMessage()
		check(err)
	}

	// Read final message from server
	_, msg, err := c.ReadMessage()
	check(err)

	// Closing websocket connection
	c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// Unmarshalling received message
	err = json.Unmarshal(msg, &m)
	check(err)
	log.Info(m.Text)
}

func check(err error) {

	if err != nil {
		log.Error(err)
	}
}
