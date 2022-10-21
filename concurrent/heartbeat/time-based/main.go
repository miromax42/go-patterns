package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	workStream, heartbeat := Work(ctx.Done(), time.Second)

	wg := sync.WaitGroup{}
	wg.Add(2)

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

func Work(done <-chan struct{}, pulseInterval time.Duration) (<-chan time.Time, <-chan struct{}) {
	workStream := make(chan time.Time)
	pulseStream := make(chan struct{})

	go func() {
		defer close(workStream)
		defer close(pulseStream)

		pulse := time.Tick(pulseInterval)
		work := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case pulseStream <- struct{}{}:
			default:
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case workStream <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case r := <-work:
				sendResult(r)
			case <-pulse:
				sendPulse()
			}
		}

	}()

	return workStream, pulseStream
}
