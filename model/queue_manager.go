package model

import (
	"fmt"
	"time"
)

func RemoveTasksWithExpiredTTL(taskQueue *TaskQueue) {
	for {
		time.Sleep(time.Millisecond * 500) // Scanning frequency

		taskQueue.QueueLock.Lock()
		for i, task := range taskQueue.Tasks {
			if task.TTL > 0 && time.Since(task.FinishTime.Time).Seconds() >= task.TTL && !task.FinishTime.IsZero() {
				fmt.Printf("Removing task %d due to expired TTL!\n", task.NumberInQueue)
				fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")
				taskQueue.Tasks = append(taskQueue.Tasks[:i], taskQueue.Tasks[i+1:]...)
			}
		}
		taskQueue.QueueLock.Unlock()
	}
}
