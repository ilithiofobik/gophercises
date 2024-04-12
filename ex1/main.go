package main

import (
	"fmt"
	"flag"
	"encoding/csv"
	"os"
	"log"
	"time"
	"strings"
	"math/rand"
)

func normalizeString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func timer(seconds int, done chan bool) {
	time.Sleep(time.Duration(seconds) * time.Second)
	done <- true
}

func quizzer(data [][]string, done chan bool, points* int) {
	for i, line := range data {
		question, answer := line[0], line[1]
		
		fmt.Printf("Question %d: %s\n", i, question)

		var response string
		fmt.Scanf("%s\n", &response)

		if normalizeString(response) == normalizeString(answer) {
			*points++
		}
	}

	done <- true
}

func main() {
	filename := flag.String("filename", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("timelimit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions")

	flag.Parse()

	fmt.Printf("Filename: %s\n", *filename)
	fmt.Printf("Time Limit: %d\n", *timeLimit)

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

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
	}

	done := make(chan bool)
	points := 0

	go timer(*timeLimit, done)
	go quizzer(data, done, &points)

	<- done

	fmt.Printf("You scored %d out of %d\n", points, len(data))
}