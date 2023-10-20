package main

import (
	"APQP/handlers"
	"APQP/model"
	"APQP/worker"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	logger *log.Logger
)

func main() {
	// Create a directory for log files
	os.MkdirAll("logs", os.ModePerm)
	// Create a file for log messages
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Initialize the logger
	logger = log.New(logFile, "APQP: ", log.Ldate|log.Ltime|log.Lshortfile)

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

	taskQueue := model.TaskQueue{
		Tasks: []*model.Task{},
	}

	taskQueueCh := make(chan *model.Task, *maxConcurrentTasksPtr)

	var wg sync.WaitGroup

	for i := 0; i < *maxConcurrentTasksPtr; i++ {
		wg.Add(1)
		go worker.Worker(&wg, taskQueueCh, logger)
	}

	// Routing
	http.HandleFunc("/enqueue", handlers.EnqueueTask(&taskQueue, taskQueueCh, logger))
	http.HandleFunc("/tasks", handlers.ListTasks(&taskQueue, logger))

	go func() {
		wg.Wait()
		close(taskQueueCh)
	}()

	// Start a goroutine for periodic TTL checking
	go model.RemoveTasksWithExpiredTTL(&taskQueue, logger)

	logger.Printf("The server is running on port %s\n\n", *portPtr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
