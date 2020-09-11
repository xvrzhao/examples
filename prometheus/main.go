package main

import (
	"examples/prometheus/prom"
	"math/rand"
	"time"
)

type TaskQueue chan *Task

// produceTasks 生产任务。
func produceTasks() TaskQueue {
	q := make(TaskQueue, 1000)

	go func(queue chan<- *Task) {
		// 每过 [0, 5) 秒向队列推送 3 个任务
		for ; true; time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond) {
			// 收集指标：队列剩余任务
			prom.RemainTasksInQueueGauge.Set(float64(len(queue)))

			for i := 0; i < 3; i++ {
				queue <- NewTask()
			}
		}
	}(q)

	return q
}

// handleTasks 消费任务。
func handleTasks(taskQueue TaskQueue, workerNum int) {
	for i := 0; i < workerNum; i++ {
		go NewWorker(i).Work(taskQueue)
	}
}

func init() {
	Seed()
}

func main() {
	handleTasks(produceTasks(), 3)
	select {}
}
