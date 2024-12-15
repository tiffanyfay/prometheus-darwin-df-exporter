package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	diskcollector "github.com/tiffanyfay/prometheus-darwin-df-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":2112", "The address to listen on for HTTP requests.")

var (
	size = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_size_bytes",
			Help: "Total disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	used = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_used_bytes",
			Help: "Used disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	avail = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_avail_bytes",
			Help: "Filesystem used disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	capacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_capacity_percent",
			Help: "Filesystem capacity(%) used of disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	iused = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_iused",
			Help: "Filesystem iused disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	ifree = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_ifree",
			Help: "Filesystem ifree disk space",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
	piused = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darwin_df_iused_percent",
			Help: "Filesystem iused percent",
		},
		[]string{"filesystem", "device", "mountpoint"},
	)
)

func recordDiskUsage() {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			// Get disk free/usage for each filesystem
			diskUsages, err := diskcollector.GetDiskUsages()
			if err != nil {
				log.Printf("Error getting disk usage: %v", err)
			}
			for _, fs := range diskUsages {
				used.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.Used))
				size.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.Size))
				avail.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.Available))
				capacity.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.Capacity))
				iused.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.IUsed))
				ifree.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.IFree))
				piused.WithLabelValues(fs.Filesystem, fs.Filesystem, fs.MountedOn).Set(float64(fs.PIUsed))
				log.Printf("Disk info recorded for filesystem %s", fs.Filesystem)
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
