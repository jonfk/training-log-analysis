#!/bin/bash
set -e -x

if [ ! -e "grafana-latest.x86_64.tar.gz" ]; then
    wget https://grafanarel.s3.amazonaws.com/builds/grafana-latest.x86_64.tar.gz
fi

sudo docker build -t jonfk/grafana .

sudo docker run -d -p 3000:3000 --link influxdb:influxdb --name grafana jonfk/grafana
