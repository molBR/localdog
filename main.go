package main

import (
	"log"
	"net/http"

	"openlocaldog/api"
	"openlocaldog/statsd"
)

func main() {
	go statsd.StartUDPListener(8125) // Porta padrão do DogStatsD

	http.HandleFunc("/metrics", api.HandleGetMetrics)
	http.HandleFunc("/reset", api.HandleResetMetrics)
	http.HandleFunc("/dashboard", api.HandleDashboard)
	http.HandleFunc("/cardinality", api.HandleGetCardinality)

	log.Println("LocalDog API server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
