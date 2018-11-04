package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarm/serial"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	teleinfo "github.com/j-vizcaino/goteleinfo"
)

func main() {
	//-----------------------
	registry := prometheus.NewRegistry()
	opsProcessed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	registry.Register(opsProcessed)
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(time.Second)
		}
	}()
	//-----------------------

	var serialDevice string
	flag.StringVar(&serialDevice, "device", "/dev/ttyUSB0", "Serial port to read frames from")
	flag.Parse()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	frames, err := initTeleinfo(ctx, serialDevice)
	if err != nil {
		panic(err)
	}

	go func() {
		for frame := range frames {
			log.Printf("%#v\n", frame)
		}
		log.Println("No more frames from channel (closed)")
	}()

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.ListenAndServe(":2112", nil)
}

func initTeleinfo(ctx context.Context, serialDevice string) (<-chan teleinfo.Frame, error) {
	port, err := teleinfo.OpenPort(serialDevice)
	if err != nil {
		return nil, fmt.Errorf("Error teleinfo open port: %v", err)
	}

	frames := make(chan teleinfo.Frame)

	go func(port *serial.Port) {
		defer port.Close()
		reader := teleinfo.NewReader(port)

		log.Println("Reading teleinfo serial port " + serialDevice)
		for {
			select {
			case <-ctx.Done():
				log.Println("Done with teleinfo reading")
				close(frames)
				return
			default:
			}

			frame, err := reader.ReadFrame()
			if err != nil {
				log.Printf("Error reading frame from '%s' (%s)\n", serialDevice, err)
				continue
			}
			frames <- frame
		}
	}(port)

	return frames, nil
}
