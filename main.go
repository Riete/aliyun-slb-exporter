package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log"
	"net/http"

	"github.com/riete/aliyun-slb-exporter/exporter"
)

const ListenPort string = "10002"

func main() {
	slb := exporter.SlbExporter{}
	slb.InitGauge()
	registry := prometheus.NewRegistry()
	registry.MustRegister(slb)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ListenPort), nil))
}
