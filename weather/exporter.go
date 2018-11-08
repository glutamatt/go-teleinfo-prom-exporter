package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const weatherURL = "http://api.openweathermap.org/data/2.5/weather?units=metric&id=2972444&appid="

func main() {
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		panic(fmt.Errorf("No weather api key provided: run with 'WEATHER_API_KEY=xxx ...'"))
	}
	registry := prometheus.NewRegistry()
	temperature := &MetricCollector{desc: prometheus.NewDesc("weather_temperature", "Température Instantanée, en degrés Celsius", nil, nil), valType: prometheus.GaugeValue}
	humidity := &MetricCollector{desc: prometheus.NewDesc("weather_humidity", "Humidité Instantanée, en %", nil, nil), valType: prometheus.GaugeValue}
	registry.MustRegister(temperature)
	registry.MustRegister(humidity)

	go func() {
		type Resp struct {
			Main struct {
				Temp     float64
				Humidity float64
			}
		}
		errInc := 0
		for {
			resp, err := http.Get(weatherURL + weatherAPIKey)
			if err == nil {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					var r Resp
					if err := json.Unmarshal(body, &r); err == nil {
						temperature.val, humidity.val = r.Main.Temp, r.Main.Humidity
						errInc = 0
						fmt.Printf("\rSet last values %.1f°C , %.0f%% humidité at %s", temperature.val, humidity.val, time.Now().Format("02/01 15:04:05"))
					} else {
						errInc++
						log.Printf("Error on json.Unmarshal: %v\n", err)
					}
				} else {
					errInc++
					log.Printf("Error on ioutil.ReadAll: %v\n", err)
				}
			} else {
				log.Printf("Error on http.Get: %v\n", err)
				errInc++
			}
			time.Sleep(time.Duration((errInc+1)*10) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Println("Exposing metrics on http://0.0.0.0:2113/metrics ...\n")
	http.ListenAndServe(":2113", nil)
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
