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
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bhojpur/speech/pkg/oscilloscope/connector"
	myrender "github.com/bhojpur/speech/pkg/oscilloscope/render"
	cli "github.com/jawher/mow.cli"
)

func washData(input []byte) []float32 {

	tmp := strings.Split(string(input), ":")
	res := make([]float32, 0)
	for i := 0; i < len(tmp); i++ {
		tmp64, _ := strconv.ParseFloat(string(tmp[i]), 32)
		res = append(res, float32(tmp64))
	}
	return res
}

func cmdEmulator(cmd *cli.Cmd) {
	cmd.Action = func() {
		ctx := context.Background()
		conn := connector.NewEmulator(ctx)
		rd := myrender.NewRender(ctx, 1280, 640, conn)
		rd.Start()
	}
}

func cmdPortAudio(cmd *cli.Cmd) {
	cmd.Action = func() {
		ctx := context.Background()
		conn := connector.NewPortAudio(ctx)
		rd := myrender.NewRender(ctx, 1280, 640, conn)
		rd.Start()
	}
}

func cmdSerial(cmd *cli.Cmd) {
	cmd.Spec = "PORT_NAME... BAUD_RATE"
	pn := cmd.StringArg("PORT_NAME", "", "The serial port name")
	br := cmd.StringArg("BAUD_RATE", "", "The serial baud rate")
	if pn == nil || br == nil {
		fmt.Println("PORT_NAME and BAUD_RATE not set, using default value")
		//"/dev/cu.usbserial-12BP0136"
		return
	}
	brn, err := strconv.ParseInt(*br, 10, 32)
	if err != nil {
		return
	}

	cmd.Action = func() {
		ctx := context.Background()
		conn := connector.NewSerial(ctx, *pn, int(brn))
		rd := myrender.NewRender(ctx, 1280, 640, conn)
		rd.Start()
	}
}

func main() {
	log.Println("Bhojpur Speech oscilloscope utility")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	app := cli.App("speechview", "An audio signal analysis ocilloscope using standard I/O ports")

	app.Command("source", "need specific a signal source", func(cmd *cli.Cmd) {
		cmd.Command("serial", "use a serial input", cmdSerial)
		cmd.Command("portaudio", "use a audio(microphone) input", cmdPortAudio)
		cmd.Command("emulator", "A sine wave emulator", cmdEmulator)
	})

	app.Run(os.Args)
}
