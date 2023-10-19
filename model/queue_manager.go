package model

import (
	"log"
	"time"
)

func RemoveTasksWithExpiredTTL(taskQueue *TaskQueue, logger *log.Logger) {
	for {
		time.Sleep(time.Millisecond * 500) // Scanning frequency

		taskQueue.QueueLock.Lock()
		for i, task := range taskQueue.Tasks {
			if task.TTL > 0 && time.Since(task.FinishTime.Time).Seconds() >= task.TTL && !task.FinishTime.IsZero() {
				logger.Printf("Removing task %d due to expired TTL\n\n", task.NumberInQueue)
				taskQueue.Tasks = append(taskQueue.Tasks[:i], taskQueue.Tasks[i+1:]...)
			}
		}
		taskQueue.QueueLock.Unlock()
	}
}
