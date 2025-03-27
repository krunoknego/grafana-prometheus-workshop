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
	fmt.Fprintln(w, "Hello from go-app1")
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
	httpRequests.WithLabelValues(r.URL.Path, "500").Inc()
	// Simulate 5XX error
	w.WriteHeader(http.StatusInternalServerError)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "Product endpoint reached")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "307").Inc()
	// Simulate a 3XX error
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	httpRequests.WithLabelValues(r.URL.Path, "200").Inc()
	fmt.Fprintln(w, "Logout endpoint reached")
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

	log.Println("go-app1 listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
