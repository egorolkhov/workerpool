package main

import (
	"log"
	"time"
	"workerpool/internal/config"
	"workerpool/internal/workerpool"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err)
	}

	wp := workerpool.NewWorkerPool(cfg.WorkerCount)
	wp.Run()

	for i, str := range cfg.Tasks {
		if i == 5 {
			wp.Resize(2)
		}
		if i == 15 {
			wp.Resize(10)
		}
		if i == 25 {
			wp.Resize(5)
		}
		wp.Tasks <- str
		time.Sleep(100 * time.Millisecond)
	}

	wp.Stop()
}
