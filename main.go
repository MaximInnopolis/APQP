package main

import (
	"APQP/handlers"
	"APQP/logger"
	"APQP/model"
	"APQP/worker"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	defer logger.Close()

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

	logger.Printf("Number of concurrent tasks: %d\n", *maxConcurrentTasksPtr)

	taskQueue := model.TaskQueue{}

	var wg sync.WaitGroup

	for i := 0; i < *maxConcurrentTasksPtr; i++ {
		wg.Add(1)
		go worker.Worker(&wg, &taskQueue)
	}

	// Routing
	http.HandleFunc("/enqueue", handlers.EnqueueTask(&taskQueue))
	http.HandleFunc("/tasks", handlers.ListTasks(&taskQueue))

	go func() {
		wg.Wait()
	}()

	// Start a goroutine for periodic TTL checking
	go taskQueue.RemoveTasksWithExpiredTTL()

	logger.Printf("The server is running on port %s\n\n", *portPtr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
