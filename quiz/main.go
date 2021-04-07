package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var csvFilePath = flag.String("csv", "", "path to csv file in format 'question, answer'")
var timeout = flag.Duration("timeout", time.Second*15, "quiz answer timeout")
var shuffle = flag.Bool("shuffle", false, "shuffle questions")

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	problemProvider := &CsvProblemProvider{CsvFilePath: *csvFilePath}
	answerProvider := &StdinAnswerProvider{}

	quiz := &Quiz{
		ProblemProvider: problemProvider,
		AnswerProvider:  answerProvider,
		Timeout:         *timeout,
		Shuffle:         *shuffle,
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
	Shuffle bool
}

type QuizResult struct {
	Solved int
	Total  int
}

func (qr QuizResult) String() string {
	return fmt.Sprintf("solved %d of %d", qr.Solved, qr.Total)
}

func swapProblemArr(arr []Problem) func(i, j int) {
	return func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func (q *Quiz) Run() QuizResult {
	problems := q.ProblemProvider.Problems()
	answers := q.AnswerProvider.Answer()

	problemsArr := []Problem{}
	for problem := range problems {
		problemsArr = append(problemsArr, problem)
	}
	if q.Shuffle {
		rand.Shuffle(len(problemsArr), swapProblemArr(problemsArr))
	}

	quizResult := QuizResult{}

	fmt.Printf("Starting quiz with timeout=%s\n", q.Timeout)
	fmt.Printf("Please, press enter when you'll be ready.\n")
	<-answers

	quizFailed := false

	for _, problem := range problemsArr {
		quizResult.Total++

		if quizFailed {
			// to draw all problems and set total to valid value
			continue
		}

		fmt.Println(problem.Question)

		select {
		case <-time.After(q.Timeout):
			fmt.Printf("Time out.\n")
			quizFailed = true
		case answer := <-answers:
			if problem.Answer != answer {
				fmt.Printf("Incorrect.\n")
				quizFailed = true
			} else {
				fmt.Printf("Correct.\n")
				quizResult.Solved++
			}
		}
	}
	return quizResult
}
