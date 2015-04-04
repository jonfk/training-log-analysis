#!/bin/bash
set -e -x

if [ ! -e "influxdb_0.9.0-rc18_amd64.deb" ]; then
    wget http://get.influxdb.org/influxdb_0.9.0-rc18_amd64.deb
fi

sudo docker build -t jonfk/influxdb .

sudo docker run -d -p 8086:8086 -p 8083:8083 -p 8087:8087 -v /var/influxdb/jonfk/raft:/var/influxdb/jonfk/raft -v /var/influxdb/jonfk/db:/var/influxdb/jonfk/db -v /var/influxdb/jonfk/state:/var/influxdb/jonfk/state --name influxdb jonfk/influxdb
