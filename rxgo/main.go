package main

import (
	"fmt"
	"sync"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	ch := make(chan rxgo.Item)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- rxgo.Of(i)
		}
		close(ch)
	}()
	obs := rxgo.FromChannel(ch)

	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range obs.Observe() {
				fmt.Println(item.V)
			}
		}()
	}
	wg.Wait()
}
