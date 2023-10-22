package model

import (
	"APQP/logger"
	"APQP/utils"
	"sync"
	"time"
)

// Task is a structure for representing a task
type Task struct {
	NumberInQueue int
	Status        string
	N             int // Number of elements in the progression
	D             float64
	N1            float64 // First element in the progression
	I             float64 // In seconds (interval)
	TTL           float64 // Time in seconds for storing the result
	CurrentIter   float64 // Current element in the progression
	TaskingTime   time.Time
	StartTime     utils.CustomTime
	FinishTime    utils.CustomTime
}

// TaskQueue is a structure for managing tasks
type TaskQueue struct {
	tasks     []*Task
	ptr       int
	QueueLock sync.Mutex
}

func (t *TaskQueue) GetTasks() []*Task {
	return t.tasks
}

func (t *TaskQueue) DeQueue() (tempTask *Task, isEmpty bool) {
	t.QueueLock.Lock()
	defer t.QueueLock.Unlock()
	if len(t.tasks)-t.ptr <= 0 {
		return nil, true
	}
	tempTask = t.tasks[t.ptr]
	t.ptr += 1
	return tempTask, false
}

func (t *TaskQueue) EnQueue(task *Task) {
	t.QueueLock.Lock()
	t.tasks = append(t.tasks, task)
	t.QueueLock.Unlock()
}

func (t *TaskQueue) RemoveTasksWithExpiredTTL() {
	for {
		time.Sleep(time.Millisecond * 500) // Scanning frequency

		t.QueueLock.Lock()
		for i, task := range t.tasks {
			if task.TTL > 0 && time.Since(task.FinishTime.Time).Seconds() >= task.TTL && !task.FinishTime.IsZero() {
				logger.Printf("Removing task %d due to expired TTL\n\n", task.NumberInQueue)
				t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
				if i <= t.ptr {
					t.ptr -= 1
				}
			}
		}
		t.QueueLock.Unlock()
	}
}
