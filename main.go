package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	ENTROPY_THRESHOLD = 4.0
)

type NGramMatrix map[string]map[rune]float64

func shannonEntropy(input string) float64 {
	counts := make(map[rune]int)
	for _, c := range input {
		counts[c]++
	}

	result := 0.0
	listLen := float64(len(counts))
	for _, count := range counts {
		freq := float64(count) / listLen
		result -= freq * math.Log(freq)
	}
	return result
}

func ngram(input string, ngramSize int) []string {
	var result []string
	input = strings.ToLower(input)
	for i := 0; i < len(input)-ngramSize+1; i++ {
		for j := 0; j < ngramSize; j++ {
			if !validChar(rune(input[i+j])) {
				continue
			}
		}
		result = append(result, input[i:i+ngramSize])
	}
	return result
}

func validChar(r rune) bool {
	return r >= 32 && r <= 126
}

func markovTrain(path string, ngramSize int) (float64, NGramMatrix) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	counts := make(NGramMatrix)

	scanner := bufio.NewScanner(file)
	padding := strings.Repeat(" ", 2)
	for scanner.Scan() {
		// Padding the string with whitespace helps generalize the probabilities for markov chains
		line := padding + strings.TrimSpace(scanner.Text()) + padding
		ngrams := ngram(line, ngramSize)
		for _, ngram := range ngrams {
			bigram := ngram[:ngramSize-1]
			nextChar := rune(ngram[ngramSize-1])
			if _, ok := counts[bigram]; !ok {
				counts[bigram] = make(map[rune]float64)
			}
			counts[bigram][nextChar]++
		}
	}

	// normalize
	probs := make(NGramMatrix)
	for bigram, nextCharCounts := range counts {
		total := 0.0
		for _, count := range nextCharCounts {
			total += count
		}
		probs[bigram] = make(map[rune]float64)
		for char, count := range nextCharCounts {
			probs[bigram][char] = math.Log(count / total)
		}
	}

	goodProbs := calculateProbabilities("good.txt", probs)
	badProbs := calculateProbabilities("bad.txt", probs)

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
	threshold := (minGood + maxBad) / 2.0

	return threshold, counts

}

func calculateProbabilities(path string, transitionProbs NGramMatrix) []float64 {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	var outputProbs []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		prob, err := averageTransitionProbability(line, transitionProbs)
		if err != nil {
			continue
		}

		outputProbs = append(outputProbs, prob)
	}
	return outputProbs
}

func averageTransitionProbability(input string, probs NGramMatrix) (float64, error) {
	logProb := 0.0

	ngramSize := 0
	for k, _ := range probs {
		ngramSize = len(k) + 1
		break
	}

	ngrams := ngram(input, ngramSize)
	for _, ngram := range ngrams {
		bigram := ngram[:ngramSize-1]
		nextChar := rune(ngram[ngramSize-1])
		if _, ok := probs[bigram]; ok {
			// The probabilities are stored as log probabilities in the model to avoid underflow
			if prob, ok := probs[bigram][nextChar]; ok {
				logProb += prob
			} else {
				return 0, fmt.Errorf("No probability for %v", ngram)
			}
		} else {
			return 0, fmt.Errorf("No probability for %v", ngram)
		}
	}

	return logProb, nil
}

// Would this fail for unicode?
func markovCheck(input string, threshold float64, probs NGramMatrix) bool {
	prob, err := averageTransitionProbability(input, probs)
	if err != nil {
		return false
	}
	return prob > threshold
}

func main() {
	ngramSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	threshold, probs := markovTrain(os.Args[2], ngramSize)
	readFile, err := os.Open(os.Args[3])

	defer readFile.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fileScanner := bufio.NewScanner(readFile)

	good, total := 0, 0

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if strings.ToLower(line) == "null);" {
			continue
		}

		valid := markovCheck(line, threshold, probs)
		if valid {
			good++
			fmt.Println(line)
		} else {
			//fmt.Println(line)
		}
		total++
		// ent := shannonEntropy(line)
		// fmt.Printf("entropy(%v:%v) markov(%v) -> %+v\n", ent, ent > ENTROPY_THRESHOLD, valid, line)
	}
	fmt.Printf("Good: %d, Total: %d, Accuracy: %f", good, total, float64(good)/float64(total))

}
