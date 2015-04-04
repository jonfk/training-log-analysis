FROM ubuntu:trusty
MAINTAINER Jonathan Fok kan <jfokkan@gmail.com>

USER root

#[Monit]
#RUN apt-get update
#RUN apt-get install -y -qq monit
#RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

#[Influxdb] Install
COPY influxdb_0.9.0-rc18_amd64.deb /tmp/
RUN dpkg -i /tmp/influxdb_0.9.0-rc18_amd64.deb && rm -f /tmp/influxdb_0.9.0-rc18_amd64.deb

#RUN mkdir -p /var/influxdb/jonfk;chown -R influxdb:influxdb /var/influxdb/jonfk/

#USER influxdb

#[Influxdb] My configs
COPY conf /etc/jonfk

#[Influxdb] Ports
EXPOSE 8083
EXPOSE 8086
EXPOSE 8087
#EXPOSE 8090
#EXPOSE 8099

#[Influxdb] Volumes
VOLUME /var/influxdb/jonfk/raft
VOLUME /var/influxdb/jonfk/db
VOLUME /var/influxdb/jonfk/state

#[Influxdb] run
#CMD /opt/influxdb/influxd -pidfile /var/opt/influxdb/run/influxd.pid -config /etc/opt/influxdb/influxdb.conf
CMD ["/opt/influxdb/influxd", "-pidfile", "/var/opt/influxdb/run/influxd.pid", "-config", "/etc/jonfk/influxdb.conf"]