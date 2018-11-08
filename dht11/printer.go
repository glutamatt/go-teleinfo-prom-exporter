package main

import (
	"fmt"
	"log"
	"time"

	"github.com/d2r2/go-dht"
	"github.com/d2r2/go-logger"
)

/*
TO COMPILE DIRECTLY ON RASPBERRY because cross compiling cgo is a mess on arch :/
*/
func main() {
	logger.ChangePackageLogLevel("dht", 0)
	for {
		temperature, humidity, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, 4, false, 10)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("{\"temperature\":%f, \"humidity\":%f}\n", temperature, humidity)
		time.Sleep(3 * time.Second)
	}
}
