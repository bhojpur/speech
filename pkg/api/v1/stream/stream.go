package stream

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
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/bhojpur/speech/go/mpg123"
	"github.com/bhojpur/speech/go/portaudio"
	"github.com/bhojpur/speech/pkg/utils"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	ADDR string = "localhost"
	PORT int    = 4000
)

type StreamServer struct{}

func NewServer() *StreamServer {
	server := &StreamServer{}
	return server
}

func (s *StreamServer) Audio(empty *emptypb.Empty, a Streamer_AudioServer) error {
	files, err := ioutil.ReadDir("./audios")
	utils.Chk(err)

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(files))
	fmt.Println("Sequence: ", randomIndex+1)
	file := files[randomIndex]

	fmt.Println("Playing: ", file.Name())
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	utils.Chk(err)

	utils.Chk(decoder.Open("./audios/" + file.Name()))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	utils.Chk(err)
	defer stream.Close()

	utils.Chk(stream.Start())
	defer stream.Stop()
	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		utils.Chk(err)

		a.Send(&Data{
			Sequence: int32(randomIndex + 1),
			Filename: file.Name(),
			Rate:     rate,
			Channels: int64(channels),
			Data:     audio,
		})
	}
	return nil
}

func (s *StreamServer) mustEmbedUnimplementedStreamerServer() {}
