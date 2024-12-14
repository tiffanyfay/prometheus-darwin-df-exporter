package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	diskcollector "github.com/tiffanyfay/prometheus-macos-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":2112", "The address to listen on for HTTP requests.")

var (
	usedDisk = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_used_bytes",
			Help: "Used disk space",
		},
		[]string{"line"}, // Define a label named "line"
	)
)

func init() {
	// Register the gauge with the Prometheus registry
	prometheus.MustRegister(usedDisk)
}

func main() {
	flag.Parse()
	fmt.Println("Endpoint: http://localhost:2112/metrics")

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// // Add go runtime metrics and process collectors.
	// reg.MustRegister(
	// 	collectors.NewGoCollector(),
	// 	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	// )

	// Register custom metrics
	reg.MustRegister(usedDisk)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func recordDiskUsage() {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			// Get disk usage for each line
			diskUsages, err := diskcollector.GetDiskUsages()
			if err != nil {
				log.Printf("Error getting disk usage: %v", err)
			}
			for _, filesystem := range diskUsages {
				usedDisk.WithLabelValues(filesystem.Filesystem).Set(float64(filesystem.Used))
				log.Printf("Disk usage recorded for line %s: %f", filesystem.Filesystem, filesystem.Used)
			}
		}
	}()
}
