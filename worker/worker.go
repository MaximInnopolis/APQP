package worker

import (
	"APQP/logger"
	"APQP/model"
	"APQP/utils"
	"sync"
	"time"
)

func processTask(task *model.Task) {
	task.Status = "In progress"
	task.StartTime = utils.CustomTime{Time: time.Now()}

	for task.CurrentIter != (task.N1 + float64(task.N-1)*task.D) {
		time.Sleep(time.Duration(task.I) * time.Second)
		task.CurrentIter = model.CalculateCurrentIter(task)
		logger.Printf("Task %d\n", task.NumberInQueue)
		logger.Println("Current iter:", task.CurrentIter)
		logger.Printf("Status: %s\n\n", task.Status)
	}

	task.Status = "Completed"
	task.FinishTime = utils.CustomTime{Time: time.Now()}
	logger.Printf("Task %d completed\n\n", task.NumberInQueue)
}

func Worker(wg *sync.WaitGroup, taskQueue *model.TaskQueue) {
	defer wg.Done()
	for {

		if task, ok := taskQueue.DeQueue(); ok {
			time.Sleep(time.Millisecond * 100)
		} else {
			processTask(task)
		}
	}
}
