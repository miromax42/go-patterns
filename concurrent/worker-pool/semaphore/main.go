package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const (
		taskCount   = 10
		workerCount = 3
	)

	args := make([]int, taskCount)
	for i := range args {
		args[i] = i
	}

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	sem := make(chan struct{}, workerCount)
	for i := 0; i < taskCount; i++ {
		arg := i
		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() {
				<-sem
			}()

			Do(args[arg])
		}()
	}

	wg.Wait()
}

func Do(i int) {
	time.Sleep(time.Second)
	fmt.Println("task: ", i)
}
