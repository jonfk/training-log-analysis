FROM ubuntu:trusty
MAINTAINER Jonathan Fok kan <jfokkan@gmail.com>

USER root

#[Influxdb] Install
COPY influxdb_latest_amd64.deb /tmp/
RUN dpkg -i /tmp/influxdb_latest_amd64.deb && rm -f /tmp/influxdb_latest_amd64.deb

#[Influxdb] My configs
COPY conf /etc/jonfk

#[Influxdb] Ports
EXPOSE 8083
EXPOSE 8086

#[Influxdb] Volumes
VOLUME /var/influxdb/jonfk/db
VOLUME /var/log/jonfk/

WORKDIR /var/log/jonfk

#[Influxdb] Environment
ENV STDOUT /dev/null
ENV STDERR /var/log/jonfk/influx.log

#[Influxdb] run
CMD /opt/influxdb/influxd -pidfile /var/run/influxdb/influxd.pid -config /etc/opt/influxdb/influxdb.conf >>$STDOUT 2>>$STDERR
