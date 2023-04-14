package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/mattdietz/gibberish/markov"
)

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

func main() {
	app := &cli.App{
		Name: "gibberish",
		Commands: []*cli.Command{
			{
				Name:   "train",
				Action: trainModel,
			},
			{
				Name:   "trainWords",
				Action: trainModelWords,
			},
			{
				Name:   "test",
				Action: testModel,
			},
			{
				Name:   "testBoth",
				Action: testBoth,
			},
			{
				Name:   "generate",
				Action: generateText,
			},
		},
	}

	app.Run(os.Args)
}

func trainModel(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) != 5 {
		fmt.Println("Usage: gibberish train <ngram size> <training file> <positive test file> <negative test file> <output file>")
		return nil
	}

	ngramArg := ctx.Args().Get(0)
	trainingPath := ctx.Args().Get(1)
	positiveTestPath := ctx.Args().Get(2)
	negativeTestPath := ctx.Args().Get(3)
	outputPath := ctx.Args().Get(4)

	ngramSize, err := strconv.Atoi(ngramArg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	model := markov.Train(trainingPath, positiveTestPath, negativeTestPath, ngramSize)
	model.Save(outputPath)

	return nil
}

func trainModelWords(ctx *cli.Context) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}

func testModel(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) != 2 {
		fmt.Println("Usage: gibberish test <model file> <test file>")
		return nil
	}

	modelPath := ctx.Args().Get(0)
	testPath := ctx.Args().Get(1)
	readFile, err := os.Open(testPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer readFile.Close()

	model := markov.LoadModel(modelPath)

	fileScanner := bufio.NewScanner(readFile)

	good, total := 0, 0

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if len(line) < 3 {
			continue
		}

		valid, _ := model.Test(line)
		if valid {
			good++
			fmt.Println(line)
		} else {
			fmt.Println(line)
		}
		total++
		// ent := shannonEntropy(line)
		// fmt.Printf("entropy(%v:%v) markov(%v) -> %+v\n", ent, ent > ENTROPY_THRESHOLD, valid, line)
	}
	fmt.Printf("Good: %d, Total: %d, Accuracy: %f", good, total, float64(good)/float64(total))
	return nil
}

func testBoth(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) != 3 {
		fmt.Println("Usage: gibberish test <good model> <bad model> <test file>")
		return nil
	}

	goodModelPath := ctx.Args().Get(0)
	badModelPath := ctx.Args().Get(1)
	testPath := ctx.Args().Get(2)
	readFile, err := os.Open(testPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer readFile.Close()

	goodModel := markov.LoadModel(goodModelPath)
	badModel := markov.LoadModel(badModelPath)

	fileScanner := bufio.NewScanner(readFile)

	good, bad := 0, 0

	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if len(line) < 3 {
			continue
		}

		goodProb, _ := goodModel.Probability(line)
		badProb, _ := badModel.Probability(line)
		if goodProb > badProb {
			good++
		} else if badProb > goodProb {
			bad++
		} else {
			fmt.Println("TIE", line)
		}
	}
	total := good + bad
	goodPct := float64(good) / float64(total)
	badPct := float64(bad) / float64(total)
	fmt.Printf("Rated Good: (%d/%v), Rated Bad: (%d/%v), Total: %d\n", good, goodPct, bad, badPct, total)
	return nil
}

func generateText(ctx *cli.Context) error {
	if len(ctx.Args().Slice()) != 2 {
		fmt.Println("Usage: gibberish generate <model file> <length>")
		return nil
	}
	modelPath := ctx.Args().Get(0)
	strLen, err := strconv.Atoi(ctx.Args().Get(1))
	if err != nil {
		fmt.Println(err)
		return err
	}
	model := markov.LoadModel(modelPath)
	fmt.Println(model.Generate(strLen, "  "))
	return nil
}
