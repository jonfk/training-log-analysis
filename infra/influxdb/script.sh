#!/bin/bash
set -e -x

if [ ! -e "influxdb_latest_amd64.deb" ]; then
    #wget http://get.influxdb.org/influxdb_0.9.0-rc18_amd64.deb
    wget https://s3.amazonaws.com/influxdb/influxdb_latest_amd64.deb
fi

sudo docker build -t jonfk/influxdb .

sudo docker run -d -p 8086:8086 -p 8083:8083 -v /var/influxdb/jonfk/db:/var/influxdb/jonfk/db -v /var/log/jonfk/:/var/log/jonfk/influxdb --name influxdb jonfk/influxdb

sleep 5

curl -G --user root:root 'http://localhost:8086/query' --data-urlencode "q=CREATE DATABASE traininglog"
curl -G --user root:root 'http://localhost:8086/query' --data-urlencode "q=CREATE USER jonfk WITH PASSWORD 'password'"
curl -G --user root:root 'http://localhost:8086/query' --data-urlencode "q=GRANT ALL ON traininglog TO jonfk"
