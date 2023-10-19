package model

import (
	"sync"
	"time"
)

// Task структура для представления задачи
type Task struct {
	ID          int
	Status      string
	N           int // Кол-во элементов в прогрессии
	D           float64
	N1          float64 // Первый элемент в прогрессии
	I           float64 // В секундах (интервал)
	TTL         float64 // В секундах время хранения результата
	CurrentIter float64 // Текущий элемент прогрессии
	TaskingTime time.Time
	StartTime   CustomTime
	FinishTime  CustomTime
}

// TaskQueue структура для управления задачами
type TaskQueue struct {
	Tasks     []*Task
	QueueLock sync.Mutex
}
