package language

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
	"archive/zip"
	"bytes"
	"embed"
	"fmt"
	"io"
	"strings"
)

//go:embed models
var languageModels embed.FS

func loadJson(language Language, ngramLength int) []byte {
	ngramName := getNgramNameByLength(ngramLength)
	isoCode := strings.ToLower(language.IsoCode639_1().String())
	zipFilePath := fmt.Sprintf("models/%s/%ss.json.zip", isoCode, ngramName)
	zipFileBytes, _ := languageModels.ReadFile(zipFilePath)
	zipFile, _ := zip.NewReader(bytes.NewReader(zipFileBytes), int64(len(zipFileBytes)))
	jsonFile, _ := zipFile.File[0].Open()
	defer jsonFile.Close()
	jsonFileContent, _ := io.ReadAll(jsonFile)
	return jsonFileContent
}
