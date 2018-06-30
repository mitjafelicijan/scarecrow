package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}

	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println("--->", os.Args[1], req.URL.String())

		response := fmt.Sprintf("%s -> %s", os.Args[1], req.URL.String())
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), nil)
}
