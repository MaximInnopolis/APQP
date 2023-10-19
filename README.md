# APQP
Arithmetic Progression Queue Processor

APQP is a service for processing tasks in a queue using arithmetic progression. 
It provides a mechanism for task management, processing and status tracking. 
The service is developed in Go language and uses HTTP API for interaction with clients.

### Requirements and dependencies
You will need the following to start the project:

Go: Make sure you have Go installed on your computer. 
If not, you can download and install it from the official Go website.

### Install and run
Clone the repository with the project:

```
git clone https://github.com/MaximInnopolis/APQP.git
```

Navigate to the project catalog:

```
cd APQP
```

Install dependencies for the project using the following command:

```
go get
```

Run service using following command:

```
go run main.go -Port 8000 -MaxConcurrentTasks 5
```

To find out what parameters can be entered use following command:

```
go run main.go -help
``` 

### Usage Examples
After you start the service, you can perform the following actions:

Start a task: To add a task to the queue, use an HTTP **POST** request to the **/enqueue** endpoint. Example request:

Viewing the list of tasks: To get a list of current tasks, send a **GET** request to the **/tasks** endpoint.
Example request:
``` 
curl http://localhost:8000/tasks
``` 

### APIs and endpoints

The service provides the following HTTP endpoints:

- **POST /enqueue**: Adds a task to the queue. 
The request body must contain the task parameters in JSON format. Example:

``` 
{
    "N": 10,
    "D": 2,
    "N1": 1,
    "I": 3,
    "TTL": 3600
}
``` 

- **GET /tasks**: Returns a list of current tasks in the queue. 
The response contains data in JSON format representing the tasks.