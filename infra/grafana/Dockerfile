FROM ubuntu:trusty
MAINTAINER Jonathan Fok kan <jfokkan@gmail.com>

USER root

#[Grafana] Install
COPY grafana-latest.x86_64.tar.gz /tmp/
RUN mkdir -p /opt/grafana && tar xvzf /tmp/grafana-latest.x86_64.tar.gz -C /opt/grafana && rm -f /tmp/grafana-latest.x86_64.tar.gz


#[Grafana] My config
#COPY conf /etc/jonfk

#[Grafana] Ports
EXPOSE 3000

#[Influxdb] Volumes

WORKDIR /opt/grafana/grafana-2.0.0-beta2/

#[Influxdb] run
CMD /opt/grafana/grafana-2.0.0-beta2/grafana web