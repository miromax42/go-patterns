package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	for v := range take(ctx.Done(), repeat(ctx.Done(), []int{1, 2, 3}), 10000) {
		time.Sleep(time.Second / 10)
		fmt.Println(v)
	}

	fmt.Printf("goroutine count when exit: %d\n", runtime.NumGoroutine())
}

func repeat[T any](done <-chan struct{}, values []T) <-chan T {
	repeatStream := make(chan T)

	go func() {
		defer close(repeatStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case repeatStream <- v:
				}
			}
		}
	}()

	return repeatStream
}

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
