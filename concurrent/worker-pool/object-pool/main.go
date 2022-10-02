package main

import (
	"fmt"
	"sync"
	"time"
)

type Doer interface {
	Do(i int)
}

type Worker struct{}

func (_ Worker) Do(i int) {
	time.Sleep(time.Second)
	fmt.Println("task: ", i)
}

func main() {
	const (
		taskCount   = 10
		workerCount = 3
	)

	task := make([]int, taskCount)
	for i := range task {
		task[i] = i
	}

	var worker Doer = NewWorkerPool(workerCount)

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		arg := i
		go func() {
			defer wg.Done()
			worker.Do(arg)
		}()
	}

	wg.Wait()
}

type WorkerPool struct {
	workers chan *Worker
}

func NewWorkerPool(count int) *WorkerPool {
	ch := make(chan *Worker, count)

	for i := 0; i < count; i++ {
		ch <- &Worker{}
	}

	return &WorkerPool{workers: ch}
}

func (p *WorkerPool) Do(i int) {
	w := <-p.workers
	defer func() {
		p.workers <- w
	}()

	w.Do(i)
}
