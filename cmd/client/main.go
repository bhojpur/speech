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
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	pb "github.com/bhojpur/speech/pkg/api/v1/stream"
	"github.com/bhojpur/speech/pkg/portaudio"
	"github.com/bhojpur/speech/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	log.Println("Bhojpur Speech streaming client (MP3)")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	wd, _ := os.Getwd()
	certFile := filepath.Join(wd, "ssl", "cert.pem")
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("Bhojpur Speech: Error creating credentials: %s\n", err)
	}

	serverAddr := fmt.Sprintf(
		"%s:%s",
		utils.GetenvDefault("ADDR", pb.ADDR),
		utils.GetenvDefault("PORT", strconv.Itoa(pb.PORT)),
	)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatalf("Bhojpur Speech: Fail to dial: %s\n", err)
	}

	defer conn.Close()
	client := pb.NewStreamerClient(conn)

	stream, err := client.Audio(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal("Bhojpur Speech: audio client error: ", err)
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int16, 8192)
	var portAudioStream *portaudio.Stream

	for {
		time.Sleep(50 * time.Millisecond)
		utils.CallClear()
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("Bhojpur Speech: cannot receive response: ", err)
		}
		log.Printf("Bhojpur Speech: client playing: %d - %s", res.GetSequence(), res.GetFilename())

		// log.Printf("audio data: ", res.GetData())

		if portAudioStream == nil {
			portAudioStream, err = portaudio.OpenDefaultStream(0, int(res.GetChannels()), float64(res.GetRate()), len(out), &out)
			utils.Chk(err)
			defer portAudioStream.Close()

			utils.Chk(portAudioStream.Start())
			defer portAudioStream.Stop()
		}

		utils.Chk(binary.Read(bytes.NewBuffer(res.GetData()), binary.LittleEndian, out))
		utils.Chk(portAudioStream.Write())
	}
}
