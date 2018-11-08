package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	registry := prometheus.NewRegistry()
	temperature := &MetricCollector{desc: prometheus.NewDesc("dht11_temperature", "Température Instantanée, en degrés Celsius", nil, nil), valType: prometheus.GaugeValue}
	humidity := &MetricCollector{desc: prometheus.NewDesc("dht11_humidity", "Humidité Instantanée, en %", nil, nil), valType: prometheus.GaugeValue}
	registry.MustRegister(temperature)
	registry.MustRegister(humidity)

	go func() {
		//{"temperature":20.000000, "humidity":63.000000}
		jsonDecoder := json.NewDecoder(os.Stdin)
		line := struct {
			Temperature float64
			Humidity    float64
		}{}
		for {
			if err := jsonDecoder.Decode(&line); err != nil {
				log.Printf("Error on jsonDecoder.Decode: %v\n", err)
				continue
			}

			temperature.val, humidity.val = line.Temperature, line.Humidity
			fmt.Printf("\rSet last values %.1f°C , %.0f%% humidité at %s", temperature.val, humidity.val, time.Now().Format("02/01 15:04:05"))
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Println("Exposing metrics on http://0.0.0.0:2114/metrics ...\n")
	http.ListenAndServe(":2114", nil)
}

type MetricCollector struct {
	val     float64
	valType prometheus.ValueType
	labels  []string
	desc    *prometheus.Desc
}

func (l *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- l.desc
}
func (l *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	m, err := prometheus.NewConstMetric(l.desc, l.valType, float64(l.val), l.labels...)
	if err != nil {
		log.Printf("MetricCollector Collect error on NewConstMetric: %v", err)
		return
	}
	ch <- m
}
