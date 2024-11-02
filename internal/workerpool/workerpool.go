package workerpool

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	wg          sync.WaitGroup
	Tasks       chan string
	StopChan    chan struct{}
	WorkerCount int
	mu          sync.Mutex
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		wg:          sync.WaitGroup{},
		Tasks:       make(chan string),
		StopChan:    make(chan struct{}),
		WorkerCount: workerCount,
		mu:          sync.Mutex{},
	}
}

func (wp *WorkerPool) Run() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for i := 0; i < wp.WorkerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.Tasks)
	wp.wg.Wait()
}

func (wp *WorkerPool) Resize(newWorkerCount int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	oldWorkerCount := wp.WorkerCount
	delta := newWorkerCount - wp.WorkerCount
	if newWorkerCount > wp.WorkerCount {
		fmt.Println(oldWorkerCount)
		wp.WorkerCount = newWorkerCount
		for i := 0; i < delta; i++ {
			wp.wg.Add(1)
			go wp.worker(oldWorkerCount + i + 1)
		}
	} else {
		fmt.Println(oldWorkerCount)
		for i := 0; i < -delta; i++ {
			wp.StopChan <- struct{}{}
		}
		wp.WorkerCount = newWorkerCount
	}
}

func (wp *WorkerPool) worker(num int) {
	defer wp.wg.Done()
	for {
		select {
		case res, ok := <-wp.Tasks:
			if !ok {
				return
			}
			fmt.Println("Worker number", num, "Result", res)
		case <-wp.StopChan:
			fmt.Println("Worker number", num, "closed")
			return
		}
	}
}
