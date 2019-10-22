package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Question represents a simple problem,answer relation
type Question struct {
	problem, answer string
}

//var done chan int = make(chan int)

// ParseFile parses a csv file
func ParseFile(filename string) (questions []Question) {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file %s", filename))
	}
	defer file.Close()

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	questions = make([]Question, len(lines))

	for i, q := range lines {
		questions[i] = Question{
			q[0],
			strings.TrimSpace(q[1]),
		}
	}

	return
}

// Quizz displays the quizz to the user and computes the score
func Quizz(questions []Question, timeLimit int) (score int) {
	if len(questions) == 0 {
		exit("Empty quizz file")
	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	for i, q := range questions {
		fmt.Printf("Problem #%v: %v = ", i+1, q.problem)
		answerCh := make(chan string)
		go func() {
			var userAnswer string
			fmt.Scanf("%s\n", &userAnswer)
			answerCh <- userAnswer
		}()
		select {
		case <-timer.C:
			return
		case userAnswer := <-answerCh:
			if q.answer == strings.ToLower(userAnswer) {
				score++
			}
		}
	}

	return
}

func exit(s string) {
	fmt.Println(s)
	os.Exit(1)
}

func main() {
	var filename = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	var timeLimit = flag.Int("limit", 15, "the time limit for the quizz in seconds")
	var shuffle = flag.Bool("shuffle", false, "if the questions order has to be shuffled")
	flag.Parse()

	questions := ParseFile(*filename)
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}
	score := Quizz(questions, *timeLimit)

	fmt.Printf("\nYou scored %v out of %v.\n", score, len(questions))
}
