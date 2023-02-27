package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var messageChan chan string

func addSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeData(w http.ResponseWriter) (int, error) {
	data := fmt.Sprintf(`{"time": "%s"}`, <-messageChan)

	return fmt.Fprintf(w, "data: %s\n\n", data)
}

func sseStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addSSEHeaders(w)

		messageChan = make(chan string)

		defer func() {
			close(messageChan)
			messageChan = nil
		}()

		flusher, _ := w.(http.Flusher)
		for {
			_, err := writeData(w)
			if err != nil {
				log.Println(err)
			}
			flusher.Flush()
		}
	}
}

func main() {
	http.HandleFunc("/stream", sseStream())

	go func() {
		for {
			if messageChan != nil {
				now := time.Now().Format(time.DateTime)
				messageChan <- now
			}

			time.Sleep(1 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":5000", nil))
}
