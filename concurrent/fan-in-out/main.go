package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	name, result int
}

func (r Result) Print() {
	fmt.Printf("worker-%d, result: %d\n", r.name, r.result)
}

func main() {
	const (
		countMessages = 10
		countSplit    = 3
	)

	ch := make(chan int)
	go func() {
		for i := 0; i < countMessages; i++ {
			ch <- i
		}

		close(ch)
	}()

	inChs := split[int](ch, countSplit)
	outChs := make([]chan Result, countSplit)

	for i, sch := range inChs {
		outChs[i] = Do(i, sch)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for v := range merge[Result](outChs...) {
			v.Print()
		}
	}()

	wg.Wait()
}

func split[T any](ch <-chan T, n int) []chan T {
	result := make([]chan T, n)
	for i := range result {
		result[i] = make(chan T, n)
	}

	go func() {
		defer func() {
			for i := range result {
				close(result[i])
			}
		}()
		var i int

		for v := range ch {
			result[i%n] <- v
			i++
		}

	}()

	return result
}

func merge[T any](in ...chan T) chan T {
	out := make(chan T)

	wg := sync.WaitGroup{}
	wg.Add(len(in))

	for _, ch := range in {
		go func(ch chan T) {
			defer wg.Done()

			for v := range ch {
				out <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func Do(name int, in <-chan int) chan Result {
	out := make(chan Result)

	fmt.Printf("worker-%d stated!\n", name)

	go func() {
		defer func() {
			close(out)
		}()

		for v := range in {
			time.Sleep(time.Second / 2)
			out <- Result{name, v}
		}
	}()

	return out
}
