package worker

import (
	"APQP/model"
	"APQP/utils"
	"log"
	"sync"
	"time"
)

func ProcessTask(task *model.Task, logger *log.Logger) {
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

func Worker(wg *sync.WaitGroup, taskQueueCh chan *model.Task, logger *log.Logger) {
	defer wg.Done()
	for {
		task, ok := <-taskQueueCh
		if !ok {
			return
		}

		ProcessTask(task, logger)

		// After completing one task, attempt to take the next one from the queue
		select {
		case newTask, ok := <-taskQueueCh:
			if !ok {
				return
			}
			logger.Printf("Received new task %d\n\n", newTask.NumberInQueue)
			ProcessTask(newTask, logger)
		default:
			// If there are no tasks in the queue, simply proceed to the next worker iteration
			continue
		}
	}
}
