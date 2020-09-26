package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Panicln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "hello"); err != nil {
		log.Println(err)
	}
}
