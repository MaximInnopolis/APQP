FROM golang:latest

ENV GO111MODULE=on
ENV GOPROXY=proxy.golang.org

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["./main", "-Port", "8080", "-MaxConcurrentTasks", "5"]
