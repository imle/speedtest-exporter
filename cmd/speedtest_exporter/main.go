package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"speedtest_exporter/internal/pkg/exporter"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.MustRegister(&exporter.SpeedTestExporter{})

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Fatal(err)
	}
}
