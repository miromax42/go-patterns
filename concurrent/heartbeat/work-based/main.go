package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	workStream, heartbeat := Work(ctx.Done())

	for {
		select {
		case <-workStream:
			fmt.Println("work done")
		case <-heartbeat:
			fmt.Println("heartbeat")
		case <-ctx.Done():
			return
		}
	}
}

func Work(done <-chan struct{}) (<-chan time.Time, <-chan struct{}) {
	workStream := make(chan time.Time)
	pulseStream := make(chan struct{}, 1)

	go func() {
		defer close(workStream)
		defer close(pulseStream)

		work := time.Tick(time.Second)

		for {
			select {
			case pulseStream <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case workStream <- <-work:
			}
		}

	}()

	return workStream, pulseStream
}
