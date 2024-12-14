package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	diskcollector "github.com/tiffanyfay/prometheus-macos-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":2112", "The address to listen on for HTTP requests.")

var (
	usedDisk = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "macos_df_used_bytes",
			Help: "Used disk space",
		},
		[]string{"device", "mountpoint"},
	)
)

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
				usedDisk.WithLabelValues(filesystem.Filesystem, filesystem.MountedOn).Set(float64(filesystem.Used))
				log.Printf("Disk usage recorded for filesystem %s: %d", filesystem.Filesystem, filesystem.Used)
			}
		}
	}()
}

func main() {
	flag.Parse()
	fmt.Println("Endpoint: http://localhost:2112/metrics")

	// Create non-global registry.
	// reg := prometheus.NewRegistry()

	// // Add go runtime metrics and process collectors.
	// reg.MustRegister(
	// 	collectors.NewGoCollector(),
	// 	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	// )

	recordDiskUsage()

	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
