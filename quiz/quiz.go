package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Question struct {
	problem string
	answer  string
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as csv for "+filePath, err)
	}
	return records
}

func readQuiz(records [][]string, shuffle bool) []Question {
	quiz := make([]Question, len(records))
	for idx, record := range records {
		quiz[idx] = Question{problem: strings.TrimSpace(record[0]), answer: strings.TrimSpace(record[1])}
	}

	if shuffle {
		quiz = shuffleQuiz(quiz)
	}
	return quiz
}

func readAnswer() <-chan string {
	ansChan := make(chan string)
	go func() {
		var ans string
		fmt.Scanf("%s\n", &ans)
		ansChan <- ans
	}()
	return ansChan
}

func shuffleQuiz(quiz []Question) []Question {
	shuffled := make([]Question, len(quiz))
	perm := rand.Perm(len(quiz))
	for i, v := range perm {
		shuffled[v] = quiz[i]
	}
	return shuffled
}

func readQuizParameters() (timeLmt int, filepath string, shuffle bool) {
	flag.IntVar(&timeLmt, "limit", 10, "Time limit for the quiz")
	flag.StringVar(&filepath, "file", "data/problems.csv", "File for the quiz")
	flag.BoolVar(&shuffle, "shuffle", true, "shuffle questions?")
	flag.Parse()
	return
}

func getScore(answers []string, quiz []Question) int {
	correct := 0
	for idx := range answers {
		if answers[idx] == quiz[idx].answer {
			correct++
		}
	}
	return correct
}

func playQuiz(quiz []Question, timeLmt int) int {
	fmt.Printf("Press enter to start the quiz")
	fmt.Scanf("%s\n")
	answers := make([]string, len(quiz))
	timeout := time.After(time.Duration(timeLmt) * time.Second)
problemLoop:
	for idx, question := range quiz {
		fmt.Printf("Problem #%v : %s = ", idx+1, question.problem)
		select {
		case <-timeout:
			break problemLoop
		case answers[idx] = <-readAnswer():
		}
	}

	return getScore(answers, quiz)
}

func main() {
	timeLmt, filepath, shuffle := readQuizParameters()
	quiz := readQuiz(readCsvFile(filepath), shuffle)
	score := playQuiz(quiz, timeLmt)
	fmt.Printf("\nYou scored %v out of %v\n", score, len(quiz))
}
