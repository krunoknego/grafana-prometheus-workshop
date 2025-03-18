package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from go-app1")
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("go-app1 listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
