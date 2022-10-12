package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for val := range orDone(ctx.Done(), repeat(ctx.Done(), []int{1})) {
		time.Sleep(time.Second / 10)
		fmt.Println(val)
	}
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
