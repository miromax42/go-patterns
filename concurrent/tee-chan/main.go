package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	out1, out2 := teeChan(ctx.Done(), repeat(ctx.Done(), []int{1, 2}))

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		for v := range out1 {
			time.Sleep(time.Second / 10)
			fmt.Println(v)
		}
	}()

	go func() {
		defer wg.Done()

		for v := range out2 {
			time.Sleep(time.Second / 10)
			fmt.Println(v)
		}
	}()

	wg.Wait()
}

func teeChan[T any](done <-chan struct{}, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
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
