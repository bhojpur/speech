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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const content = `These sentences are intended for testing purposes.
 Do not use them in production!
 By the way, they consist of 23 words in total.`

const corpus = `There are many attributes associated with good software.
 Some of these can be mutually contradictory, and different customers and participants may have different priorities.
 Weinberg provides an example of how different goals can have a dramatic effect on both effort required and efficiency.
 Furthermore, he notes that programmers will generally aim to achieve any explicit goals which may be set, probably at the expense of any other quality attributes.
 Sommerville has identified four generalised attributes which are not concerned with what a program does, but how well the program does it:
 Maintainability, Dependability, Efficiency, Usability`

const expectedSentencesFileContent = `There are many attributes associated with good software.
 Some of these can be mutually contradictory, and different customers and participants may have different priorities.
 Weinberg provides an example of how different goals can have a dramatic effect on both effort required and efficiency.
 Furthermore, he notes that programmers will generally aim to achieve any explicit goals which may be set, probably at the expense of any other quality attributes.
 `

const expectedSingleWordsFileContent = `there
 attributes
 associated
 software
 `

const expectedWordPairsFileContent = `there attributes
 associated software
 these mutually
 contradictory different
 `

const expectedUnigramJsonModel = `
 {
	 "language":"ENGLISH",
	 "ngrams":{
		 "1/10":"n o s",
		 "1/100":"b g l m",
		 "1/20":"d r",
		 "1/25":"h",
		 "1/50":"f w",
		 "13/100":"t",
		 "3/100":"a c p u y",
		 "3/50":"i",
		 "7/50":"e"
	 }
 }`

const expectedBigramJsonModel = `
 {
	 "language":"ENGLISH",
	 "ngrams":{
		 "1/1":"by he",
		 "1/10":"nc nd ng no ns od of os si",
		 "1/13":"ta to",
		 "1/14":"ed em ey",
		 "1/2":"fo wa wo",
		 "1/3":"al ar ay ce co ct po pr pu uc ur us",
		 "1/5":"de do ds du nt on or ot rd re ro rp st",
		 "1/6":"io is",
		 "2/13":"ti",
		 "2/3":"in",
		 "2/5":"se",
		 "2/7":"es",
		 "3/13":"te",
		 "3/14":"en",
		 "4/13":"th"
	 }
 }`

const expectedTrigramJsonModel = `
 {
	 "language":"ENGLISH",
	 "ngrams":{
		 "1/1":"are ces con cti ded duc for ion ist nce nde not nsi nte odu ose pos pro pur rds rod rpo sis tal the tot uct urp use way wor",
		 "1/2":"ons ord ota sti tin tio",
		 "1/3":"enc end ent tes",
		 "1/4":"ese est hem hes hey ing int sen ses",
		 "2/3":"ten"
	 }
 }`

const expectedQuadrigramJsonModel = `
 {
	 "language":"ENGLISH",
	 "ngrams":{
		 "1/1":"cons ctio duct ence ende ente esti hese inte nces nded nsis nten oduc onsi ords oses otal pose prod purp rodu rpos sent sist stin test ting tion tota ucti urpo word",
		 "1/2":"tenc tend",
		 "1/4":"them thes they"
	 }
 }`

const expectedFivegramJsonModel = `
 {
	 "language":"ENGLISH",
	 "ngrams":{
		 "1/1":"consi ction ducti ences ended enten estin inten nsist oduct onsis poses produ purpo roduc rpose sente sting tence tende testi these total uctio urpos words",
		 "1/2":"ntenc ntend"
	 }
 }`

func TestCreateAndWriteLanguageModelFiles(t *testing.T) {
	inputFilePath := createTempInputFile(content)
	outputDirectoryPath, _ := ioutil.TempDir("", "lingualOutputDirectory")

	err := CreateAndWriteLanguageModelFiles(inputFilePath, outputDirectoryPath, English, "\\p{L}")

	assert.Nil(t, err, "language model files could not be created")

	files, _ := os.ReadDir(outputDirectoryPath)

	assert.Equal(t, 5, len(files), "number of language model files is not correct")

	unigramsFile := files[4]
	bigramsFile := files[0]
	trigramsFile := files[3]
	quadrigramsFile := files[2]
	fivegramsFile := files[1]

	assert.Equal(t, "unigrams.json.zip", unigramsFile.Name())
	assert.Equal(t, "bigrams.json.zip", bigramsFile.Name())
	assert.Equal(t, "trigrams.json.zip", trigramsFile.Name())
	assert.Equal(t, "quadrigrams.json.zip", quadrigramsFile.Name())
	assert.Equal(t, "fivegrams.json.zip", fivegramsFile.Name())

	assertLanguageModelFileContent(t, outputDirectoryPath, unigramsFile.Name(), "unigrams.json", expectedUnigramJsonModel)
	assertLanguageModelFileContent(t, outputDirectoryPath, bigramsFile.Name(), "bigrams.json", expectedBigramJsonModel)
	assertLanguageModelFileContent(t, outputDirectoryPath, trigramsFile.Name(), "trigrams.json", expectedTrigramJsonModel)
	assertLanguageModelFileContent(t, outputDirectoryPath, quadrigramsFile.Name(), "quadrigrams.json", expectedQuadrigramJsonModel)
	assertLanguageModelFileContent(t, outputDirectoryPath, fivegramsFile.Name(), "fivegrams.json", expectedFivegramJsonModel)

	cleanUp(inputFilePath, outputDirectoryPath)
}

func TestCreateAndWriteTestDataFiles(t *testing.T) {
	inputFilePath := createTempInputFile(corpus)
	outputDirectoryPath, _ := ioutil.TempDir("", "lingualOutputDirectory")

	err := CreateAndWriteTestDataFiles(inputFilePath, outputDirectoryPath, "\\p{L}", 4)

	assert.Nil(t, err, "test data files could not be created")

	files, _ := os.ReadDir(outputDirectoryPath)

	assert.Equal(t, 3, len(files), "number of test data files is not correct")

	assertTestDataFileContent(
		t,
		files[0],
		outputDirectoryPath,
		"sentences.txt",
		expectedSentencesFileContent,
	)

	assertTestDataFileContent(
		t,
		files[1],
		outputDirectoryPath,
		"single-words.txt",
		expectedSingleWordsFileContent,
	)

	assertTestDataFileContent(
		t,
		files[2],
		outputDirectoryPath,
		"word-pairs.txt",
		expectedWordPairsFileContent,
	)

	cleanUp(inputFilePath, outputDirectoryPath)
}

func createTempInputFile(content string) string {
	inputFile, _ := ioutil.TempFile("", "text.*.txt")
	defer inputFile.Close()
	_, _ = inputFile.WriteString(content)
	return inputFile.Name()
}

func assertLanguageModelFileContent(
	t *testing.T,
	outputDirectoryPath,
	zipFileName,
	expectedFileName,
	expectedFileContent string,
) {
	zipFilePath := filepath.Join(outputDirectoryPath, zipFileName)
	zipFile, _ := zip.OpenReader(zipFilePath)
	defer zipFile.Close()

	assert.Equal(t, 1, len(zipFile.File), "zip archive does not contain exactly one file")

	jsonFile := zipFile.File[0]

	assert.Equal(t, expectedFileName, jsonFile.Name)

	jsonFileReader, _ := jsonFile.Open()
	defer jsonFileReader.Close()
	jsonFileContent, _ := io.ReadAll(jsonFileReader)

	assert.Equal(t, minify(expectedFileContent), string(jsonFileContent))
}

func assertTestDataFileContent(
	t *testing.T,
	file os.DirEntry,
	outputDirectoryPath,
	expectedFileName,
	expectedFileContent string,
) {
	assert.False(t, file.IsDir())
	assert.Equal(t, expectedFileName, file.Name())

	filePath := filepath.Join(outputDirectoryPath, file.Name())
	testDataFileContent, _ := os.ReadFile(filePath)

	assert.Equal(t, expectedFileContent, string(testDataFileContent))
}

func cleanUp(inputFilePath, outputDirectoryPath string) {
	err := os.Remove(inputFilePath)
	if err != nil {
		panic(fmt.Errorf("temporary file could not be removed: %v", err))
	}
	err = os.RemoveAll(outputDirectoryPath)
	if err != nil {
		panic(fmt.Errorf("temporary directory could not be removed: %v", err))
	}
}
