# go-teleinfo-prom-exporter

```
env GOOS=linux GOARCH=arm GOARM=7 go build -o go-teleinfo-prom-exporter main.go
scp go-teleinfo-prom-exporter   pi@192.168.0.18:.
ssh -t  pi@192.168.0.18 ./go-teleinfo-prom-exporter -device /dev/serial0
```


## Notes sur le PI

```
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