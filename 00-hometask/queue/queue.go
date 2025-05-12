package queue

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var config = map[string]int{
	"queue1":  1,
	"queue10": 10,
	"queue20": 20,
}

const maxSubscribers = 3

var (
	once  = make(map[string]*sync.Once, len(config))
	input = make(map[string]chan string, len(config))
	out   = make(map[string][]chan string)
)

func init() {
	for k, size := range config {
		inp := make(chan string, size)
		input[k] = inp
		o := sync.Once{}
		once[k] = &o
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	queueName := strings.Split(r.URL.Path, "/")[3]
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	select {
	case input[queueName] <- string(msg):
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("queue " + queueName + " is full"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("added to queue successfully"))
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	queueName := strings.Split(r.URL.Path, "/")[3]
	if _, ok := config[queueName]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("queue not found"))
		return
	} else if len(out[queueName]) >= maxSubscribers {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("max subscribers reached"))
		return
	} else {
		o := make(chan string)
		out[queueName] = append(out[queueName], o)
		go func(o <-chan string) {
			for {
				select {
				case msg := <-o:
					fmt.Println("consumer got message", msg)
				}
			}
		}(o)
		handlersOnce(queueName)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("subscribed successfully"))
}

func handlersOnce(queueName string) {
	once[queueName].Do(func() {
		go func(k string, inp <-chan string) {
			for {
				msg := <-inp
				for i := range out[k] {
					out[k][i] <- msg
				}
			}
		}(queueName, input[queueName])
	})
}
