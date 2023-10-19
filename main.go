package main

import (
	"APQP/model"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var Taskqueue model.TaskQueue
var taskQueueCh chan *model.Task

func calculateCurrentIter(task *model.Task) float64 {
	return task.CurrentIter + task.D
}

func processTask(task *model.Task) {
	task.Status = "In progress"
	task.StartTime = model.CustomTime{Time: time.Now()}

	for task.CurrentIter != (task.N1 + float64(task.N-1)*task.D) {
		time.Sleep(time.Duration(task.I) * time.Second)
		task.CurrentIter = calculateCurrentIter(task)
		fmt.Printf("Task %d\n", task.ID)
		fmt.Println("Current iter:", task.CurrentIter)
		fmt.Println("Status:", task.Status)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~pro")
	}

	task.Status = "Completed"
	task.FinishTime = model.CustomTime{Time: time.Now()}
	fmt.Printf("Task %d completed!\n", task.ID)
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-taskQueueCh
		if !ok {
			return
		}

		processTask(task)

		// После завершения одной задачи, попытаемся взять следующую из очереди
		select {
		case newTask, ok := <-taskQueueCh:
			if !ok {
				return
			}
			fmt.Println("Received new task", newTask.ID)
			processTask(newTask)
		default:
			// Если нет задач в очереди, просто переходим к следующей итерации воркера
			continue
		}
	}
}

func removeTasksWithExpiredTTL() {
	for {
		time.Sleep(time.Millisecond * 500) // Периодичность сканирования

		Taskqueue.QueueLock.Lock()
		for i, task := range Taskqueue.Tasks {
			if task.TTL > 0 && time.Since(task.FinishTime.Time).Seconds() >= task.TTL && !task.FinishTime.IsZero() {
				fmt.Printf("Deleting task %d with expired TTL\n", task.ID)
				Taskqueue.Tasks = append(Taskqueue.Tasks[:i], Taskqueue.Tasks[i+1:]...)
			}
		}
		Taskqueue.QueueLock.Unlock()
	}
}

func enqueueTask(w http.ResponseWriter, r *http.Request) {
	// Чтение данных из тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Парсинг JSON
	var request model.EnqueueRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Failed to parse JSON request", http.StatusBadRequest)
		return
	}

	task := &model.Task{
		ID:          len(Taskqueue.Tasks) + 1,
		Status:      "In Queue",
		N:           request.N,
		D:           request.D,
		N1:          request.N1,
		I:           request.I,
		TTL:         request.TTL,
		CurrentIter: request.N1,
		TaskingTime: time.Now(),
	}

	// Добавляем задачу в очередь и канал
	Taskqueue.QueueLock.Lock()
	Taskqueue.Tasks = append(Taskqueue.Tasks, task)
	Taskqueue.QueueLock.Unlock()
	taskQueueCh <- task

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task enqueued successfully.")
	fmt.Printf("Task %d is in the queue!!!!\n", task.ID)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~en")
}

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

	// Проверка на нулевое или отрицательное значение maxConcurrentTasks
	if *maxConcurrentTasksPtr <= 0 {
		fmt.Println("Number of concurrent tasks should be a positive integer.")
		return
	}

	fmt.Println("Number of concurrent tasks:", int32(*maxConcurrentTasksPtr))
	fmt.Println("Port:", *portPtr)

	//Инициализация очереди задач
	Taskqueue = model.TaskQueue{
		Tasks: []*model.Task{},
	}

	taskQueueCh = make(chan *model.Task, *maxConcurrentTasksPtr)

	var wg sync.WaitGroup

	for i := 0; i < *maxConcurrentTasksPtr; i++ {
		wg.Add(1)
		go worker(&wg)
	}

	// Routing
	http.HandleFunc("/enqueue", enqueueTask)
	http.HandleFunc("/tasks", listTasks)

	go func() {
		wg.Wait()
		close(taskQueueCh)
	}()

	// Запускаем горутину для периодической проверки TTL
	//go removeTasksWithExpiredTTL()

	fmt.Printf("The server is running on port %s\n", *portPtr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func listTasks(w http.ResponseWriter, r *http.Request) {

	Taskqueue.QueueLock.Lock()
	defer Taskqueue.QueueLock.Unlock()

	data, err := json.Marshal(Taskqueue.Tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
