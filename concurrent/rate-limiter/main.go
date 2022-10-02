package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const (
		taskCount      = 50
		limitPerSecond = 25
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, taskCount)

	limiter := make(chan struct{})
	ticker := time.NewTicker(time.Second / limitPerSecond)
	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				limiter <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	wg := sync.WaitGroup{}
	wg.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go func() {
			defer wg.Done()

			<-limiter
			ch <- RPCCall()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for value := range ch {
		fmt.Println(value)
	}
}

func RPCCall() int {
	return rand.Int()
}
