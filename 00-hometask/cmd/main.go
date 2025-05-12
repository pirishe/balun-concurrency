package main

import (
	"concurrency/queue"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /v1/queues/{queuename}/messages", queue.CreateMessage)
	http.HandleFunc("POST /v1/queues/{queuename}/subscriptions", queue.Subscribe)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
