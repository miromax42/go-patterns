package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer fmt.Println("gen exited")
			defer close(chanStream)
			for i := 0; i < 100; i++ {
				stream := make(chan interface{}, 1)
				time.Sleep(time.Second / 20)
				stream <- i
				close(stream)
				chanStream <- stream
			}

		}()
		return chanStream
	}

	for v := range bridge(ctx.Done(), genVals()) {
		fmt.Printf("%v \n", v)
	}

	time.Sleep(time.Second)

	fmt.Printf("goroutines:%d\n", runtime.NumGoroutine())
}

func bridge[T any](done <-chan struct{}, chanStream <-chan <-chan T) <-chan T {
	outputStream := make(chan T)

	go func() {
		defer close(outputStream)

		for {
			for ch := range orDone(done, chanStream) {
				for v := range orDone(done, ch) {
					select {
					case <-done:
						return
					case outputStream <- v:
					}
				}
			}
		}

	}()

	return outputStream
}

func orDone[T any](done <-chan struct{}, input <-chan T) <-chan T {
	outStream := make(chan T)

	go func() {
		defer close(outStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-input:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case outStream <- val:
				}
			}
		}
	}()

	return outStream
}
