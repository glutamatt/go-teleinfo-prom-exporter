# go-teleinfo-prom-exporter

```
env GOOS=linux GOARCH=arm GOARM=7 go build -o go-teleinfo-prom-exporter main.go
scp go-teleinfo-prom-exporter   pi@192.168.0.18:.
ssh -t  pi@192.168.0.18 ./go-teleinfo-prom-exporter -device /dev/serial0
```


## Notes sur le PI

```

# Docker install

https://medium.freecodecamp.org/the-easy-way-to-set-up-docker-on-a-raspberry-pi-7d24ced073ef

-----------------------------------
# pour regler le port serie correctement sur le fichier /dev/serial0
# se fait programmaticly en go

stty -F /dev/serial0 1200 sane evenp parenb cs7 -crtscts


----------------------------------
# temp dir as ram partitions : file /etc/fstab
# sudo mkdir /var/prometheus

tmpfs /tmp tmpfs defaults,noatime,nosuid,size=10m 0 0
tmpfs /var/tmp tmpfs defaults,noatime,nosuid,size=10m 0 0
tmpfs /var/prometheus tmpfs defaults,noatime,nosuid,size=200m 0 0
tmpfs /var/log tmpfs defaults,noatime,nosuid,mode=0755,size=10m 0 0
```

## Prometheus config copy

`scp prometheus.yml   pi@192.168.0.18:.`


## Prometheus sur le PI

```
curl -sSLO https://github.com/prometheus/prometheus/releases/download/v2.4.3/prometheus-2.4.3.linux-armv7.tar.gz
tar -xvf prometheus-2.4.3.linux-armv7.tar.gz
cd prometheus-2.4.3.linux-armv7
./prometheus --config.file=/home/pi/prometheus.yml --storage.tsdb.path=/var/prometheus/data --storage.tsdb.retention=30d --web.enable-lifecycle --web.console.libraries=console_libraries --web.console.templates=consoles
```


