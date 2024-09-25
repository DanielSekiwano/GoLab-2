package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// question struct stores a single question and its corresponding answer.
type question struct {
	q, a string
}

type score int

// check handles a potential error.
// It stops execution of the program ("panics") if an error has happened.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// questions reads in questions and corresponding answers from a CSV file into a slice of question structs.
func questions() []question {
	f, err := os.Open("quiz-questions.csv")
	check(err)
	reader := csv.NewReader(f)
	table, err := reader.ReadAll()
	check(err)
	var questions []question
	for _, row := range table {
		questions = append(questions, question{q: row[0], a: row[1]})
	}
	return questions
}

// ask asks a question and returns an updated score depending on the answer.
func ask(sC chan score, question question, curScore score) {

	fmt.Println(question.q)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter answer: ")
	scanner.Scan()
	text := scanner.Text()

	if strings.Compare(text, question.a) == 0 {
		fmt.Println("Correct!")
		curScore++
	} else {
		fmt.Println("Incorrect :-(")
	}

	sC <- curScore
}

func main() {
	sC := make(chan score, 1)
	s := score(0)
	timeout := time.After(time.Second * 5)

	for _, q := range questions() {
		go ask(sC, q, s)
		select {
		case score := <-sC:
			s = score
		case <-timeout:
			fmt.Println("Final score", s)
			return
		}
	}
}
