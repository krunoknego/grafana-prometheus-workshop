package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

var counter int64

func handler(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddInt64(&counter, 1)

	data := map[string]interface{}{
		"unique_id":     id,
		"random_number": rand.Intn(100),
		"timestamp":     time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
