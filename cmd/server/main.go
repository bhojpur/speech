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
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	pb "github.com/bhojpur/speech/pkg/api/v1/stream"
	"github.com/bhojpur/speech/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.Println("Bhojpur Speech streaming server (MP3)")
	log.Println("Copyright (c) 2018 by Bhojpur Consulting Private Limited, India.")
	log.Printf("All rights reserved.\n")

	wd, _ := os.Getwd()
	certFile := filepath.Join(wd, "ssl", "cert.pem")
	keyFile := filepath.Join(wd, "ssl", "private.key")
	creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)

	serverAddr := fmt.Sprintf(
		"%s:%s",
		utils.GetenvDefault("HOST", "localhost"),
		utils.GetenvDefault("PORT", strconv.Itoa(pb.PORT)),
	)
	listen, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("server engine failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterStreamerServer(grpcServer, pb.NewServer())

	log.Printf("server engine listening on gRPC %s\n", serverAddr)
	grpcServer.Serve(listen)
}
