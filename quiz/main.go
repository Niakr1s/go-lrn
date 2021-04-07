package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

var csvFilePath = flag.String("csv", "", "path to csv file in format 'question, answer'")

func init() {
	flag.Parse()
}

func main() {
	problemProvider := &CsvProblemProvider{CsvFilePath: *csvFilePath}
	answerProvider := &StdinAnswerProvider{}

	RunQuiz(problemProvider, answerProvider)
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

func RunQuiz(problemProvider ProblemProvider, answerProvider AnswerProvider) {
	problems := problemProvider.Problems()
	answers := answerProvider.Answer()

	solved := 0
	for problem := range problems {
		fmt.Println(problem.Question)
		answer := <-answers
		if problem.Answer != answer {
			fmt.Printf("Incorrect.\n")
			break
		}
		fmt.Printf("Correct.\n")
		solved++
	}
	fmt.Printf("You have solved %d problems.\n", solved)
}
