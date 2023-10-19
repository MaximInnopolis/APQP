package handlers

import (
	"APQP/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func EnqueueTask(taskQueue *model.TaskQueue, taskQueueCh chan *model.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read data from the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Parse JSON
		var request model.EnqueueRequest
		err = json.Unmarshal(body, &request)
		if err != nil {
			http.Error(w, "Failed to parse JSON request", http.StatusBadRequest)
			return
		}

		task := &model.Task{
			NumberInQueue: len(taskQueue.Tasks) + 1,
			Status:        "In Queue",
			N:             request.N,
			D:             request.D,
			N1:            request.N1,
			I:             request.I,
			TTL:           request.TTL,
			CurrentIter:   request.N1,
			TaskingTime:   time.Now(),
		}

		// Add the task to the queue and channel
		taskQueue.QueueLock.Lock()
		taskQueue.Tasks = append(taskQueue.Tasks, task)
		taskQueue.QueueLock.Unlock()
		taskQueueCh <- task

		// Send a successful response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Task enqueued successfully.")
		fmt.Printf("Task %d is in the queue!\n", task.NumberInQueue)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
}

func ListTasks(taskQueue *model.TaskQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskQueue.QueueLock.Lock()
		defer taskQueue.QueueLock.Unlock()

		data, err := json.Marshal(taskQueue.Tasks)
		if err != nil {
			http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
