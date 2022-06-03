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
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CreateAndWriteLanguageModelFiles creates language model files for
// accuracy report generation and writes them to a directory.
//
// `inputFilePath` is the path to a txt file used for language model creation.
// The assumed enconding of the txt file is UTF-8.
//
// `outputDirectoryPath` is the path to an existing directory where the language
// model files are to be written.
//
// `language` is the language for which to create the language models.
//
// `charClass` is a regex character class such as `\\p{L}` to restrict the
// set of characters that the language models are built from.
//
// An error is returned if:
//
// - the input file path is not absolute or does not point to an existing txt file
//
// - the input file's encoding is not UTF-8
//
// - the output directory path is not absolute or does not point to an existing directory
//
// Panics if the character class cannot be compiled to a valid regular expression.
func CreateAndWriteLanguageModelFiles(
	inputFilePath string,
	outputDirectoryPath string,
	language Language,
	charClass string,
) error {
	err := checkInputFilePath(inputFilePath)
	if err != nil {
		return err
	}
	err = checkOutputDirectoryPath(outputDirectoryPath)
	if err != nil {
		return err
	}

	unigramModel, err := createLanguageModel(
		inputFilePath,
		language,
		1,
		charClass,
		map[ngram]uint32{},
	)
	if err != nil {
		return err
	}

	bigramModel, err := createLanguageModel(
		inputFilePath,
		language,
		2,
		charClass,
		unigramModel.absoluteFrequencies,
	)
	if err != nil {
		return err
	}

	trigramModel, err := createLanguageModel(
		inputFilePath,
		language,
		3,
		charClass,
		bigramModel.absoluteFrequencies,
	)
	if err != nil {
		return err
	}

	quadrigramModel, err := createLanguageModel(
		inputFilePath,
		language,
		4,
		charClass,
		trigramModel.absoluteFrequencies,
	)
	if err != nil {
		return err
	}

	fivegramModel, err := createLanguageModel(
		inputFilePath,
		language,
		5,
		charClass,
		quadrigramModel.absoluteFrequencies,
	)
	if err != nil {
		return err
	}

	err = writeCompressedLanguageModel(
		unigramModel,
		outputDirectoryPath,
		"unigrams.json",
	)
	if err != nil {
		return err
	}

	err = writeCompressedLanguageModel(
		bigramModel,
		outputDirectoryPath,
		"bigrams.json",
	)
	if err != nil {
		return err
	}

	err = writeCompressedLanguageModel(
		trigramModel,
		outputDirectoryPath,
		"trigrams.json",
	)
	if err != nil {
		return err
	}

	err = writeCompressedLanguageModel(
		quadrigramModel,
		outputDirectoryPath,
		"quadrigrams.json",
	)
	if err != nil {
		return err
	}

	err = writeCompressedLanguageModel(
		fivegramModel,
		outputDirectoryPath,
		"fivegrams.json",
	)
	if err != nil {
		return err
	}

	return nil
}

// CreateAndWriteTestDataFiles creates test data files for accuracy report
// generation and writes them to a directory.
//
// `inputFilePath` is the path to a txt file used for test data creation.
// The assumed enconding of the txt file is UTF-8.
//
// `outputDirectoryPath` is the path to an existing directory where the test
// data files are to be written.
//
// `charClass` is a regex character class such as `\\p{L}` to restrict the
// set of characters that the language models are built from.
//
// `maximumLines` is the maximum number of lines each test data file should have.
//
// An error is returned if:
//
// - the input file path is not absolute or does not point to an existing txt file
//
// - the input file's encoding is not UTF-8
//
// - the output directory path is not absolute or does not point to an existing directory
//
// Panics if the character class cannot be compiled to a valid regular expression.
func CreateAndWriteTestDataFiles(
	inputFilePath string,
	outputDirectoryPath string,
	charClass string,
	maximumLines int,
) error {
	err := checkInputFilePath(inputFilePath)
	if err != nil {
		return err
	}
	err = checkOutputDirectoryPath(outputDirectoryPath)
	if err != nil {
		return err
	}

	err = createAndWriteSentencesFile(
		inputFilePath,
		outputDirectoryPath,
		maximumLines,
	)
	if err != nil {
		return err
	}

	singleWords, err := createAndWriteSingleWordsFile(
		inputFilePath,
		outputDirectoryPath,
		charClass,
		maximumLines,
	)
	if err != nil {
		return err
	}

	err = createAndWriteWordPairsFile(
		singleWords,
		outputDirectoryPath,
		maximumLines,
	)
	if err != nil {
		return err
	}

	return nil
}

func createAndWriteSentencesFile(
	inputFilePath string,
	outputDirectoryPath string,
	maximumLines int,
) error {
	sentencesFilePath := filepath.Join(outputDirectoryPath, "sentences.txt")
	if _, err := os.Stat(sentencesFilePath); errors.Is(err, os.ErrExist) {
		err = os.Remove(sentencesFilePath)
		if err != nil {
			return err
		}
	}
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	var inputLines []string
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			inputLines = append(inputLines, line)
		}
	}

	sentencesFile, err := os.Create(sentencesFilePath)
	if err != nil {
		return err
	}
	defer sentencesFile.Close()

	lineCounter := 0

	for _, line := range inputLines {
		normalizedWhitespace := multipleWhitespace.ReplaceAllString(line, " ")
		removedQuotes := strings.ReplaceAll(normalizedWhitespace, "\"", "")

		if lineCounter < maximumLines {
			_, err = sentencesFile.WriteString(fmt.Sprintf("%s\n", removedQuotes))
			if err != nil {
				return err
			}
			lineCounter++
		} else {
			break
		}
	}

	return nil
}

func createAndWriteSingleWordsFile(
	inputFilePath string,
	outputDirectoryPath string,
	charClass string,
	maximumLines int,
) ([]string, error) {
	singleWordsFilePath := filepath.Join(outputDirectoryPath, "single-words.txt")
	wordRegex, err := regexp.Compile(fmt.Sprintf("[%s]{5,}", charClass))
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(singleWordsFilePath); errors.Is(err, os.ErrExist) {
		err = os.Remove(singleWordsFilePath)
		if err != nil {
			return nil, err
		}
	}

	var words []string

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	var inputLines []string
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			inputLines = append(inputLines, line)
		}
	}

	singleWordsFile, err := os.Create(singleWordsFilePath)
	if err != nil {
		return nil, err
	}
	defer singleWordsFile.Close()

	lineCounter := 0

	for _, line := range inputLines {
		removedPunctuation := punctuation.ReplaceAllString(line, "")
		removedNumbers := numbers.ReplaceAllString(removedPunctuation, "")
		normalizedWhitespace := multipleWhitespace.ReplaceAllString(removedNumbers, " ")
		removedQuotes := strings.ReplaceAll(normalizedWhitespace, "\"", "")

		for _, word := range strings.Split(removedQuotes, " ") {
			word = strings.ToLower(strings.TrimSpace(word))
			if wordRegex.MatchString(word) {
				words = append(words, word)
			}
		}
	}

	for _, word := range words {
		if lineCounter < maximumLines {
			_, err = singleWordsFile.WriteString(fmt.Sprintf("%s\n", word))
			if err != nil {
				return nil, err
			}
			lineCounter++
		} else {
			break
		}
	}

	return words, nil
}

func createAndWriteWordPairsFile(
	words []string,
	outputDirectoryPath string,
	maximumLines int,
) error {
	wordPairsFilePath := filepath.Join(outputDirectoryPath, "word-pairs.txt")

	if _, err := os.Stat(wordPairsFilePath); errors.Is(err, os.ErrExist) {
		err = os.Remove(wordPairsFilePath)
		if err != nil {
			return err
		}
	}

	var wordPairs []string

	for i := 0; i <= len(words)-2; i++ {
		if i%2 == 0 {
			wordPairs = append(wordPairs, strings.Join(words[i:i+2], " "))
		}
	}

	wordPairsFile, err := os.Create(wordPairsFilePath)
	if err != nil {
		return err
	}
	defer wordPairsFile.Close()

	lineCounter := 0

	for _, wordPair := range wordPairs {
		if lineCounter < maximumLines {
			_, err = wordPairsFile.WriteString(fmt.Sprintf("%s\n", wordPair))
			if err != nil {
				return err
			}
			lineCounter++
		} else {
			break
		}
	}

	return nil
}

func createLanguageModel(
	inputFilePath string,
	language Language,
	ngramLength int,
	charClass string,
	lowerNgramAbsoluteFrequencies map[ngram]uint32,
) (trainingDataLanguageModel, error) {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return trainingDataLanguageModel{}, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) > 0 {
			lines = append(lines, line)
		}
	}
	return newTrainingDataLanguageModel(
		lines,
		language,
		ngramLength,
		charClass,
		lowerNgramAbsoluteFrequencies,
	), nil
}

func writeCompressedLanguageModel(
	model trainingDataLanguageModel,
	outputDirectoryPath string,
	fileName string,
) error {
	zipFileName := fmt.Sprintf("%s.zip", fileName)
	zipFilePath := filepath.Join(outputDirectoryPath, zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	jsonFileWriter, err := zipWriter.Create(fileName)
	if err != nil {
		return err
	}
	_, err = jsonFileWriter.Write(model.toJson())
	if err != nil {
		return err
	}
	return nil
}

func checkInputFilePath(inputFilePath string) error {
	if !filepath.IsAbs(inputFilePath) {
		return fmt.Errorf("input file path '%s' is not absolute", inputFilePath)
	}

	fileInfo, err := os.Stat(inputFilePath)

	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("input file '%s' does not exist", inputFilePath)
	}
	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("input file path '%s' does not represent a regular file", inputFilePath)
	}

	return nil
}

func checkOutputDirectoryPath(outputDirectoryPath string) error {
	if !filepath.IsAbs(outputDirectoryPath) {
		return fmt.Errorf("output directory path '%s' is not absolute", outputDirectoryPath)
	}

	fileInfo, err := os.Stat(outputDirectoryPath)

	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("output directory '%s' does not exist", outputDirectoryPath)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("output directory path '%s' does not exist", outputDirectoryPath)
	}

	return nil
}