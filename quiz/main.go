package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var csvFilePath = flag.String("csv", "", "path to csv file in format 'question, answer'")

func init() {
	flag.Parse()
}

func main() {
	problemProvider := &CsvProblemProvider{CsvFilePath: *csvFilePath}
	answerProvider := &StdinAnswerProvider{}

	quiz := &Quiz{
		ProblemProvider: problemProvider,
		AnswerProvider:  answerProvider,
	}
	quizResult := quiz.Run()
	fmt.Printf("Quiz ended: %v\n", quizResult)
}

type AnswerProvider interface {
	Answer() <-chan string
}

type StdinAnswerProvider struct{}

func (u *StdinAnswerProvider) Answer() <-chan string {
	ch := make(chan string)

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			str := scanner.Text()
			ch <- str
		}
	}()
	return ch
}

type Problem struct {
	Question string
	Answer   string
}

type ProblemProvider interface {
	Problems() <-chan Problem
}

type CsvProblemProvider struct {
	CsvFilePath string
}

func (p *CsvProblemProvider) Problems() <-chan Problem {
	file, err := os.Open(p.CsvFilePath)
	if err != nil {
		panic(fmt.Errorf("file open error: %v", err))
	}

	ch := make(chan Problem)
	go func() {
		defer file.Close()

		r := csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(fmt.Errorf("error while reading csv file: %v", err))
			}
			problem := Problem{Question: record[0], Answer: record[1]}
			ch <- problem
		}
		close(ch)
	}()
	return ch
}

type Quiz struct {
	ProblemProvider ProblemProvider
	AnswerProvider  AnswerProvider

	Timeout time.Duration
}

type QuizResult struct {
	Solved int
	Total  int
}

func (qr QuizResult) String() string {
	return fmt.Sprintf("solved %d of %d", qr.Solved, qr.Total)
}

func (q *Quiz) Run() QuizResult {
	problems := q.ProblemProvider.Problems()
	answers := q.AnswerProvider.Answer()

	quizResult := QuizResult{}

	quizFailed := false
	for problem := range problems {
		quizResult.Total++

		if quizFailed {
			// to draw all problems and set total to valid value
			continue
		}

		fmt.Println(problem.Question)
		answer := <-answers
		if problem.Answer != answer {
			fmt.Printf("Incorrect.\n")
			quizFailed = true
			continue
		}
		fmt.Printf("Correct.\n")
		quizResult.Solved++
	}
	return quizResult
}
