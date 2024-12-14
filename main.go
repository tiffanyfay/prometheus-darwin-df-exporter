package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tiffanyfay/prometheus-macos-exporter/diskcollector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var addr = flag.String("listen-address", ":2112", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	fmt.Println("Endpoint: http://localhost:2112/metrics")

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	diskcollector.GetDiskUsage()

	// http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	// log.Fatal(http.ListenAndServe(*addr, nil))
}

func recordDiskUsage(reg prometheus.Registerer) {
	// usedDisk := promauto.With(reg).NewCounter(prometheus.CounterOpts{
	// 	Name: "disk_used_bytes",
	// 	Help: "Used disk space",
	// })

	go func() {
		for {
			time.Sleep(2 * time.Second)
			diskcollector.GetDiskUsage()
			log.Println("Disk usage recorded")
		}
	}()
}

/* node_filesystem_avail_bytes{
	device="/dev/disk1s1",
	device_error="",
	fstype="apfs",
	mountpoint="/System/Volumes/Data"
}
2.3709648896e+11
*/
