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

### Running the Application with Docker

If you prefer to run the application in a Docker container, follow these steps:

1. Build the Docker image. Open a terminal and navigate to the project directory:

``` 
cd /path/to/APQP
``` 

2. Build the Docker image using the provided Dockerfile

``` 
docker build -t apqp-app .  
``` 

3. Run the Docker container:

``` 
docker run -d -p 8000:8080 apqp-app
``` 

Now your application is running in a Docker container, and it's accessible at http://localhost:8000

### Usage Examples
After you start the service, you can perform the following actions:

Start a task: To add a task to the queue, use an HTTP **POST** request to the **/enqueue** endpoint.

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

![](https://i.imgur.com/MrIOn4N.png)

- **GET /tasks**: Returns a list of current tasks in the queue. 
The response contains data in JSON format representing the tasks.

![](https://i.imgur.com/rMKDbKr.png)

### Logging

This project utilizes logging to capture various events and messages. Logging allows us to track, analyze, and manage the actions and state of the application.

#### Logging Library

I use the standard Go `log` library to perform logging. This enables us to record messages in the log file instead of printing them to the console.

Example of logging usage:

```
log.Printf("This message will be logged.")
log.Println("This will also be logged.")
```

### Where to Find Logs
Logs are located in the logs/app.log file. You can review this file to learn more about the application's operation and track events and errors.

If you used docker to run the application,
you can see the logs using the following command:

```
docker exec <CONTAINER_ID> cat logs/app.log
```

Container ID can be found by using the command:
```
docker ps -a
```

In my case:

```
CONTAINER ID   IMAGE      COMMAND                  CREATED          STATUS          PORTS                    NAMES
40109414bcf2   apqp-app   "./main -Port 8080 -…"   13 minutes ago   Up 13 minutes   0.0.0.0:8080->8080/tcp   pedantic_wright
```

![](https://i.imgur.com/VFuNwc7.png)
