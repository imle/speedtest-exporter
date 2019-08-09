package exporter

import (
	"fmt"
	"log"
	"sync"

	"github.com/kylegrantlucas/speedtest"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "speed_test"
	subsystem = "measured"
)

var (
	client *speedtest.Client
)

type SpeedTestExporter struct {
	mu sync.Mutex

	download *prometheus.Desc
	upload   *prometheus.Desc
	ping     *prometheus.Desc
}

func init() {
	var err error
	client, err = speedtest.NewDefaultClient()
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
}

func (s *SpeedTestExporter) Describe(ch chan<- *prometheus.Desc) {
	s.download = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "download"),
		"Network download speed as measured by device in bits/s",
		[]string{"name", "server_id"},
		nil,
	)
	ch <- s.download

	s.upload = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "upload"),
		"Network upload speed as measured by device in bits/s",
		[]string{"name", "server_id"},
		nil,
	)
	ch <- s.upload

	s.ping = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "ping"),
		"Ping as measured by device in milliseconds",
		[]string{"name", "server_id"},
		nil,
	)
	ch <- s.ping
}

func (s SpeedTestExporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Choosing server")
	server, err := client.GetServer("")
	if err != nil {
		fmt.Printf("error getting server: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(
		s.ping,
		prometheus.GaugeValue,
		server.Latency,
		server.Name, server.ID)

	log.Println("Download test")
	downBps, _ := client.Download(server)
	ch <- prometheus.MustNewConstMetric(
		s.download,
		prometheus.GaugeValue,
		downBps,
		server.Name, server.ID)

	log.Println("Upload test")
	upBps, _ := client.Download(server)
	ch <- prometheus.MustNewConstMetric(
		s.upload,
		prometheus.GaugeValue,
		upBps,
		server.Name, server.ID)

	log.Println("Testing complete")

	fmt.Printf("Ping: %3.2f ms | Download: %3.2f Mbps | Upload: %3.2f Mbps\n", server.Latency, downBps, upBps)
}
