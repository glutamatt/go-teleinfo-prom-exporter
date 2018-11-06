package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tarm/serial"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	teleinfo "github.com/j-vizcaino/goteleinfo"
)

var registry *prometheus.Registry
var labels []string //metric label names

func main() {
	labels = []string{"OPTARIF", "HHPHC", "PTEC"}
	metrics := map[string]*MetricCollector{
		"BASE":  &MetricCollector{desc: prometheus.NewDesc("teleinfo_base", "Index option, en Watt heure", labels, nil), valType: prometheus.CounterValue},
		"PAPP":  &MetricCollector{desc: prometheus.NewDesc("teleinfo_ppap", "Puissance apparente, en VA (Volt Ampères)", labels, nil), valType: prometheus.GaugeValue},
		"IINST": &MetricCollector{desc: prometheus.NewDesc("teleinfo_iinst", "Intensité Instantanée, Ampères", labels, nil), valType: prometheus.GaugeValue},
	}

	var serialDevice string
	flag.StringVar(&serialDevice, "device", "/dev/serial0", "Serial port to read frames from")
	flag.Parse()
	HandleMetrics(metrics, initTeleinfo(serialDevice))
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":2112", nil)
}

func HandleMetrics(metrics map[string]*MetricCollector, frames <-chan teleinfo.Frame) {
	registry = prometheus.NewRegistry()
	for _, m := range metrics {
		m.labels = []string{}
		registry.MustRegister(m)
	}

	go func() {
		for frame := range frames {
			log.Printf("%#v\n", frame)

			labelVals := make([]string, len(labels))
			for i, l := range labels {
				if val, ok := frame.GetStringField(l); ok {
					labelVals[i] = val
				}
			}

			for f, m := range metrics {
				m.labels = labelVals
				if val, ok := frame.GetUIntField(f); ok {
					m.val = val
				}
			}
		}

		log.Println("No more frames from channel (closed)")
	}()
}

func initTeleinfo(serialDevice string) <-chan teleinfo.Frame {
	port, err := teleinfo.OpenPort(serialDevice)
	if err != nil {
		panic(fmt.Errorf("Error teleinfo open port: %v", err))
	}

	frames := make(chan teleinfo.Frame)

	go func(port *serial.Port) {
		defer port.Close()
		reader := teleinfo.NewReader(port)

		log.Println("Reading teleinfo serial port " + serialDevice)
		for {
			frame, err := reader.ReadFrame()
			if err != nil {
				log.Printf("Error reading frame from '%s' (%s)\n", serialDevice, err)
				continue
			}
			frames <- frame
		}
	}(port)

	return frames
}

type MetricCollector struct {
	val     uint
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
