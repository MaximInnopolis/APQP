package worker

import (
	"APQP/model"
	"APQP/utils"
	"fmt"
	"sync"
	"time"
)

func ProcessTask(task *model.Task) {
	task.Status = "In progress"
	task.StartTime = utils.CustomTime{Time: time.Now()}

	for task.CurrentIter != (task.N1 + float64(task.N-1)*task.D) {
		time.Sleep(time.Duration(task.I) * time.Second)
		task.CurrentIter = model.CalculateCurrentIter(task)
		fmt.Printf("Task %d\n", task.NumberInQueue)
		fmt.Println("Current iter:", task.CurrentIter)
		fmt.Println("Status:", task.Status)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}

	task.Status = "Completed"
	task.FinishTime = utils.CustomTime{Time: time.Now()}
	fmt.Printf("Task %d completed!\n", task.NumberInQueue)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")

}

func Worker(wg *sync.WaitGroup, taskQueueCh chan *model.Task) {
	defer wg.Done()
	for {
		task, ok := <-taskQueueCh
		if !ok {
			return
		}

		ProcessTask(task)

		// After completing one task, attempt to take the next one from the queue
		select {
		case newTask, ok := <-taskQueueCh:
			if !ok {
				return
			}
			fmt.Println("Received new task!", newTask.NumberInQueue)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")
			ProcessTask(newTask)
		default:
			// If there are no tasks in the queue, simply proceed to the next worker iteration
			continue
		}
	}
}
