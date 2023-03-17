package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

const NGRAM_SIZE = 2

type NGramMatrix map[rune]map[rune]float64

func characterEntropy(input string) int {
	counts := make(map[rune]int)
	for _, c := range input {
		counts[c]++
	}
	return len(counts)
}

func shannonEntropy(input string) float64 {
	counts := make(map[rune]int)
	for _, c := range input {
		counts[c]++
	}

	result := 0.0
	log2 := math.Log(2)
	listLen := float64(len(counts))
	for _, count := range counts {
		freq := float64(count) / listLen
		result -= freq * (math.Log(freq) / log2)
	}
	return result
}

func ngram(input string, n int) []string {
	var result []string
	input = strings.ToLower(input)
	for i := 0; i < len(input)-n+1; i++ {
		result = append(result, input[i:i+n])
	}
	return result
}

func validChar(r rune) bool {
	return r >= 32 && r <= 122
}

func markovTrain(path string) (float64, NGramMatrix) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	counts := make(NGramMatrix)

	for i := 32; i <= 122; i++ {
		counts[rune(i)] = make(map[rune]float64)
		for j := 32; j <= 122; j++ {
			counts[rune(i)][rune(j)] = 10
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		ngrams := ngram(line, NGRAM_SIZE)
		for i := 0; i < len(ngrams)-1; i++ {
			if !validChar(rune(ngrams[i][0])) || !validChar(rune(ngrams[i][1])) {
				continue
			}

			counts[rune(ngrams[i][0])][rune(ngrams[i][1])]++
		}
	}

	// normalize
	for row := range counts {
		sum := 0.0
		for col := range counts[row] {
			sum += float64(col)
		}
		for col := range counts[row] {
			counts[row][col] = math.Log(counts[row][col] / sum)
		}
	}

	goodProbs := calculateProbabilities("good.txt", counts)
	badProbs := calculateProbabilities("bad.txt", counts)

	minGood := goodProbs[0]
	for _, prob := range goodProbs {
		if prob < minGood {
			minGood = prob
		}
	}

	maxBad := badProbs[0]
	for _, prob := range badProbs {
		if prob > maxBad {
			maxBad = prob
		}
	}

	// TODO save the model to disk
	threshold := (minGood + maxBad) / 2

	fmt.Println(threshold)
	return threshold, counts

}

func calculateProbabilities(path string, counts NGramMatrix) []float64 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	var probs []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		probs = append(probs, averageTransitionProbability(line, counts))
	}
	return probs
}

func averageTransitionProbability(input string, counts NGramMatrix) float64 {
	logProb := 0.0
	transitionCt := 0

	for i := 0; i < len(input)-1; i++ {
		if !validChar(rune(input[i])) || !validChar(rune(input[i+1])) {
			log.Fatalf("Invalid character in input string %c %c", input[i], input[i+1])
		}
	}

	ngrams := ngram(input, 2)
	for i := 0; i < len(ngrams)-1; i++ {
		transitionCt++
		logProb += counts[rune(ngrams[i][0])][rune(ngrams[i+1][1])]
	}

	if transitionCt == 0 {
		transitionCt = 1
	}
	return math.Exp(logProb / float64(transitionCt))
}

// Would this fail for unicode?
func markovCheck(input string, threshold float64, counts NGramMatrix) bool {
	return averageTransitionProbability(input, counts) > threshold
}

func main() {
	threshold, counts := markovTrain("big.txt")
	var testStrings = []string{
		"Hello, World!",
		"Hello, World! Hello, World!",
		"Hello, World! Hello, World! Hello, World!",
		"abcabc",
		"abcabcabc",
		"alksjdlkajsdlkajsd",
		"crapdad",
		"crappad",
		"crappadcrappad",
		"foobar",
		"fooooooooobar",
		"Art Of Craft",
		"gbxhbncjn",
		"vcwvbxbxb",
		"cxtyjei153",
		"herpdederp.store",
		"Fadil",
		"tryu",
		"sdfsdf",
		"aaaaaaaaaaabbbbbbbbbbbbbbbbbbaaaaaaaaaaaaaaaaaaaaaaa",
		"Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch",
	}

	for i := 0; i < len(testStrings); i++ {
		//fmt.Printf("%s -> linearEntropy score %+v\n", testStrings[i], characterEntropy(testStrings[i]))
		// fmt.Printf("%s -> shannonEntropy score %+v\n", testStrings[i], shannonEntropy(testStrings[i]))
		fmt.Printf("%s -> valid? %+v\n", testStrings[i], markovCheck(testStrings[i], threshold, counts))
	}
}
