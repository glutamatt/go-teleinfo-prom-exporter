# go-teleinfo-prom-exporter

#### pour regler le port serie correctement sur le fichier /dev/serial0 (se fait programmaticly en go)

`stty -F /dev/serial0 1200 sane evenp parenb cs7 -crtscts`

```
env GOOS=linux GOARCH=arm GOARM=7 go build -o go-teleinfo-prom-exporter main.go
scp go-teleinfo-prom-exporter   pi@192.168.0.18:.
ssh -t  pi@192.168.0.18 ./go-teleinfo-prom-exporter -device /dev/serial0
```


## Notes sur le PI

```
# Docker install (useless)
https://medium.freecodecamp.org/the-easy-way-to-set-up-docker-on-a-raspberry-pi-7d24ced073ef
`docker run -d -p 3000:3000 --name grafana fg2it/grafana-armhf:v5.1.4`
```

## Grafana
```
wget https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana-5.3.2.linux-armv7.tar.gz 
tar -zxvf grafana-5.3.2.linux-armv7.tar.gz
cd grafana-5.3.2
./bin/grafana-server         
```

## Prometheus

#### some dirs as ram partitions

> `sudo mkdir -p /var/prometheus/data/wal`
> `sudo chmod -R 777 /var/prometheus`

then add at the end of `/etc/fstab` file 

```
tmpfs /tmp tmpfs defaults,noatime,nosuid,size=10m 0 0
tmpfs /var/tmp tmpfs defaults,noatime,nosuid,size=10m 0 0
tmpfs /var/prometheus/data/wal tmpfs defaults,noatime,nosuid,size=20m 0 0
tmpfs /var/log tmpfs defaults,noatime,nosuid,mode=0755,size=10m 0 0
```

> `sudo mount -a`

#### Prometheus config copy

> `scp prometheus.yml   pi@192.168.0.18:.`


#### Prometheus Binary

```
curl -sSLO https://github.com/prometheus/prometheus/releases/download/v2.4.3/prometheus-2.4.3.linux-armv7.tar.gz
tar -xvf prometheus-2.4.3.linux-armv7.tar.gz
cd prometheus-2.4.3.linux-armv7
./prometheus --config.file=/home/pi/prometheus.yml --storage.tsdb.path=/var/prometheus/data --storage.tsdb.retention=30d --web.enable-lifecycle --web.console.libraries=console_libraries --web.console.templates=consoles
```

## Rclone

### config 

```
curl -sSLO https://downloads.rclone.org/v1.44/rclone-v1.44-linux-arm.zip
unzip rclone-v1.44-linux-arm.zip
cd rclone-v1.44-linux-arm
```

### config 

```
pi@raspberrypi:~/.config/rclone $ cat rclone.conf
[prombackup]
type = dropbox
token = {"access_token":"*****","token_type":"bearer","expiry":"0001-01-01T00:00:00Z"}
```

> token in google keep archive note

### Run rclone

> `rclone copy /var/prometheus prombackup:prometheus-backup --exclude data/wal/**`

for quick and dirty repeat : `watch --interval 120`
