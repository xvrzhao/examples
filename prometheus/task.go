package main

import (
	"math/rand"
	"time"
)

var taskTypes = [3]string{
	"send email",
	"download file",
	"close order",
}

type Task struct {
	taskType string
}

func NewTask() *Task {
	return &Task{taskType: taskTypes[rand.Intn(3)]}
}

func (t *Task) GetType() string {
	return t.taskType
}

// Process 处理任务。
func (*Task) Process() time.Duration {
	// 模拟任务耗时，随机取 [0, 5) 秒
	d := time.Duration(rand.Intn(5000)) * time.Millisecond
	time.Sleep(d)
	return d
}

func init() {
	Seed()
}
