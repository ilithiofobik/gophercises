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

func normalizeString(s *string) string {
	return strings.TrimSpace(strings.ToLower(*s))
}

func sameString(s1, s2 *string) bool {
	return strings.Compare(normalizeString(s1), normalizeString(s2)) == 0
}

func quizzer(problems []Problem, done chan bool, points *int) {
	var response string

	for i, p := range problems {
		fmt.Printf("Question %d: %s\n", i+1, p.question)
		fmt.Scanf("%s\n", &response)

		if sameString(&response, &p.answer) {
			*points++
		}
	}

	done <- true
}

type Problem struct {
	question string
	answer   string
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))

	for i, line := range lines {
		ret[i] = Problem{
			question: line[0],
			answer:   line[1],
		}
	}

	return ret
}

func readData(filename *string) [][]string {
	file, err := os.Open(*filename)

	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	file.Close()

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func main() {
	filename := flag.String("filename", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("timelimit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions")
	flag.Parse()

	fmt.Printf("Filename: %s\n", *filename)
	fmt.Printf("Time Limit: %d\n", *timeLimit)

	data := readData(filename)
	problems := parseLines(data)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	done := make(chan bool)
	points := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	go quizzer(problems, done, &points)

	select {
	case <-done:
		break
	case <-timer.C:
		break
	}

	fmt.Printf("You scored %d out of %d\n", points, len(problems))
}
