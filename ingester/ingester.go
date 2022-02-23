package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func fastHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.WriteHeader(http.StatusOK)
}

func main() {
	kind := os.Getenv("PYROSCOPE_AGENT_BENCHMARK_INGESTER_TYPE")
	var handler http.HandlerFunc
	switch kind {
	case "slow":
		handler = slowHandler
	default:
		handler = fastHandler
	}
	http.HandleFunc("/ingest", handler)
	log.Fatal(http.ListenAndServe(":4040", nil))
}
