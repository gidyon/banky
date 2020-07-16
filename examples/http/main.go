package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := serveMux()

	log.Println("server started on port :9090")
	http.ListenAndServe(":9090", mux)
}

func serveMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/timev2", myHandler)

	return mux
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().String()
	fmt.Fprintln(w, now)
}
