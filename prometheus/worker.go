package main

import (
	"examples/prometheus/prom"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type Worker struct {
	id int
}

func NewWorker(ID int) *Worker {
	return &Worker{ID}
}

func (w *Worker) GetID() int {
	return w.id
}

func (w *Worker) GetStrID() string {
	return strconv.Itoa(w.GetID())
}

// Work Worker w 处理任务。
func (w *Worker) Work(taskQueue TaskQueue) {
	for task := range taskQueue {
		d := task.Process()
		fmt.Printf("Worker-%d handled task: %s, takes %d milliseconds.\n",
			w.GetID(), task.GetType(), d.Milliseconds())

		// 收集指标：处理任务数、任务耗时
		prom.HandledTasksCounter.With(prometheus.Labels{"worker_id": w.GetStrID(), "task_type": task.GetType()}).Inc()
		prom.TaskSecondHistogram.Observe(float64(d.Milliseconds()))
		prom.TaskSecondSummary.Observe(float64(d.Milliseconds()))
	}
}
