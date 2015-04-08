#!/bin/bash
set -e -x

#rm influxdb_0.9.0-rc18_amd64.deb || true

docker rmi $(docker images -q -f dangling=true) || true
docker rm influxdb || true

sudo rm -rf /var/influxdb/
