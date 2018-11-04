# go-teleinfo-prom-exporter

```
env GOOS=linux GOARCH=arm GOARM=7 go build -o go-teleinfo-prom-exporter main.go
scp go-teleinfo-prom-exporter   pi@192.168.0.18:.
ssh -t  pi@192.168.0.18 ./go-teleinfo-prom-exporter -device /dev/serial0
```
