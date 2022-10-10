package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	intStream := repeatFn(ctx.Done(), func() int {
		return rand.Intn(10000) + 1
	})
	// primeStream := primeFinder(ctx.Done(), intStream)

	numFinders := runtime.NumCPU()
	finders := make([]<-chan int, numFinders)
	for i := range finders {
		finders[i] = primeFinder(ctx.Done(), intStream)
	}

	defer func(start time.Time) {
		fmt.Println(time.Since(start))
	}(time.Now())

	for v := range take(ctx.Done(), fanIn(ctx.Done(), finders...), 10) {
		fmt.Println(v)
	}

}

func fanIn[T any](done <-chan struct{}, streams ...<-chan T) <-chan T {
	mergeStream := make(chan T)
	wg := sync.WaitGroup{}

	multiplex := func(inputValues <-chan T) {
		defer wg.Done()

		for v := range inputValues {
			select {
			case <-done:
				return
			case mergeStream <- v:
			}
		}
	}

	wg.Add(len(streams))
	for _, c := range streams {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(mergeStream)
	}()

	return mergeStream
}

func primeFinder(done <-chan struct{}, intStream <-chan int) <-chan int {
	isPrime := func(num int) bool {
		for i := 2; i < num/i+1; i++ {
			if num%i == 0 {
				return false
			}
		}

		return true
	}

	primeStream := make(chan int)

	go func() {
		defer close(primeStream)
		for {
			select {
			case <-done:
				return
			case v := <-intStream:
				if isPrime(v) {
					primeStream <- v
				}
			}
		}
	}()

	return primeStream
}

// helpers
func repeatFn[T any](done <-chan struct{}, fn func() T) <-chan T {
	repeatStream := make(chan T)

	go func() {
		defer close(repeatStream)
		for {
			select {
			case <-done:
				return
			case repeatStream <- fn():
			}
		}
	}()

	return repeatStream
}

func take[T any](done <-chan struct{}, valueStream <-chan T, n int) <-chan T {
	takeStream := make(chan T)

	go func() {
		defer close(takeStream)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}
