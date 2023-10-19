package main

import (
	"fmt"
)

type request struct {
	action   string
	key      string
	value    string
	response chan string
}

var requests chan request = make(chan request, 100)
var done chan struct{} = make(chan struct{})
var keyStore = make(map[string]string)

func requestMonitor() {
	keyStore["test"] = "works"

	for request := range requests {
		fmt.Println("Actor processing " + request.action + " " + request.key)
		switch request.action {
		case "put":
			keyStore[request.key] = request.value
			fmt.Printf("Put value in %s\n", request.key)
		case "delete":
			delete(keyStore, request.key)
		case "read":
			request.response <- keyStore[request.key]
		case "shutdown":
			fmt.Println("Shutting down")
			close(requests)
		}
	}

	fmt.Println("All requests processed")
	for k, v := range keyStore {
		fmt.Println(k, "value is", v)
	}
	close(done)
}

func Put(key string, value string) {
	req := request{
		action: "put",
		key:    key,
		value:  value,
	}
	requests <- req
}

func Delete(key string) {
	req := request{
		action: "delete",
		key:    key,
	}
	requests <- req
}

func Get(key string) string {
	answer := make(chan string)
	req := request{
		action:   "read",
		key:      key,
		response: answer,
	}
	requests <- req

	return <-answer
}

func Start() {
	go requestMonitor()
}

func Stop() {
	shutdown := request{action: "shutdown"}
	requests <- shutdown
	<-done
}

func main() {
	Start()
	defer Stop()

	go Put("this", "that")
	go Put("these", "those")
	go Delete("this")
	go fmt.Printf("%v\n", Get("test"))

}
