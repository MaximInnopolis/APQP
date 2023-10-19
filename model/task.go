package model

import (
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
	Tasks     []*Task
	QueueLock sync.Mutex
}
