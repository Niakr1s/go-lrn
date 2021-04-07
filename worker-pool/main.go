package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result struct {
	Url  string
	Body string
}

func main() {
	urlCh := make(chan string)
	resCh := make(chan Result)
	endCh := make(chan bool)

	go scan(urlCh, endCh)
	go workerPool(urlCh, resCh, 2)
	go resultPrinter(resCh)

	<-endCh
}

func scan(urlCh chan<- string, endCh chan<- bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		if url == "exit" {
			break
		}
		url = strings.TrimSpace(url)
		if !strings.HasSuffix(url, "http") {
			url = "http://" + url
		}
		urlCh <- url
	}
	endCh <- true
}

func workerPool(urlCh <-chan string, resCh chan<- Result, max int) {
	semCh := make(chan bool, max)

	for url := range urlCh {
		url := url
		go func() {
			semCh <- true
			fmt.Printf("worker %s start\n", url)
			res, err := worker(url)
			if err != nil {
				fmt.Printf("worker produced an error")
			}
			resCh <- res
			fmt.Printf("worker %s end\n", url)
			<-semCh
		}()
	}
}

func worker(url string) (Result, error) {
	<-time.After(time.Second * 5)
	res, err := http.DefaultClient.Get(url)
	if err != nil {
		return Result{}, fmt.Errorf("couldn't load url %s, error: %v", url, err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Result{}, fmt.Errorf("couldn't read body, url: %s, error: %v", url, err)
	}
	result := Result{Url: url, Body: string(body)}
	return result, nil
}

func resultPrinter(resultCh <-chan Result) {
	for result := range resultCh {
		fmt.Printf("%s: %d\n", result.Url, len(result.Body))
	}
}
