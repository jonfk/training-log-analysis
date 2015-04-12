#!/bin/bash
set -e -x

docker rmi $(docker images -q -f dangling=true) || true
docker rm grafana || true
