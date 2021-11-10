package prome

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "k8s.io/klog"

	collector "github.com/shdn-ibm/promethues-test-data/pkg/collectors"
)

func RunExporter() {

	c, err := collector.NewPerfCollector()
	if err != nil {
		log.Warningf("NewFSPerfCollector fails, err:%s", err)
	}

	// Use custom registry to remove default go metrics
	r := prometheus.NewRegistry()
	r.MustRegister(c)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)

	// prometheus.MustRegister(c)
	// http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var _, _ = w.Write([]byte(`<html>
            <head><title>Promethues Exporter</title></head>
            <body>
            <h1>Promethues Exporter </h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	log.Info("Beginning to serve on port :9100")
	if err = http.ListenAndServe(":9100", nil); err != nil {
		log.Error(err, "failed to start http server")
		panic(err)
	}
}