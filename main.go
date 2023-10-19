package main

import (
	"APQP/handlers"
	"APQP/model"
	"APQP/worker"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	//Flags initialization
	maxConcurrentTasksPtr := flag.Int("MaxConcurrentTasks", 2, "Number of concurrent tasks")
	portPtr := flag.String("Port", "8080", "Port to listen on")
	listenAddr := ":" + *portPtr

	helpPtr := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		return
	}

	// Check for a zero or negative value of maxConcurrentTasks
	if *maxConcurrentTasksPtr <= 0 {
		fmt.Println("Number of concurrent tasks should be a positive integer.")
		return
	}

	fmt.Println("Number of concurrent tasks:", int32(*maxConcurrentTasksPtr))

	taskQueue := model.TaskQueue{
		Tasks: []*model.Task{},
	}

	taskQueueCh := make(chan *model.Task, *maxConcurrentTasksPtr)

	var wg sync.WaitGroup

	for i := 0; i < *maxConcurrentTasksPtr; i++ {
		wg.Add(1)
		go worker.Worker(&wg, taskQueueCh)
	}

	// Routing
	http.HandleFunc("/enqueue", handlers.EnqueueTask(&taskQueue, taskQueueCh))
	http.HandleFunc("/tasks", handlers.ListTasks(&taskQueue))

	go func() {
		wg.Wait()
		close(taskQueueCh)
	}()

	// Start a goroutine for periodic TTL checking
	go model.RemoveTasksWithExpiredTTL(&taskQueue)

	fmt.Printf("The server is running on port %s\n\n", *portPtr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
