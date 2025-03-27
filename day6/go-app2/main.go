package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "goapp_http_requests_total",
		Help: "Number of incoming HTTP requests",
	},
	[]string{"path", "status"},
)

func init() {
	prometheus.MustRegister(httpRequests)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "Hello from go-app2")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "OK")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "User endpoint reached")
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "Order endpoint reached")
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "404").Inc()
	// Simulate a 4XX error
	w.WriteHeader(http.StatusNotFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "Login endpoint reached")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path).Inc()
	httpRequests.WithLabelValues(r.URL.Path, "500").Inc()

	//Simulate a 5XX error
	w.WriteHeader(http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/product", productHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("go-app2 listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
