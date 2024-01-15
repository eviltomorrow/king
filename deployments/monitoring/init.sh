#!/bin/bash

mkdir -p ${DATA_HOME}/prometheus/data
mkdir -p ${DATA_HOME}/otel-collector
mkdir -p ${DATA_HOME}/cassandra/{db,log}
mkdir -p ${DATA_HOME}/loki/{log,data}
mkdir -p ${DATA_HOME}/promtail/{conf,log}
mkdir -p ${DATA_HOME}/grafana/{data,log,provisioning}
mkdir -p ${DATA_HOME}/grafana/provisioning/{dashboards,datasources}

chmod -R 777 ${DATA_HOME}/prometheus
chmod -R 777 ${DATA_HOME}/otel-collector
chmod -R 777 ${DATA_HOME}/cassandra
chmod -R 777 ${DATA_HOME}/loki
chmod -R 777 ${DATA_HOME}/promtail
chmod -R 777 ${DATA_HOME}/grafana

touch ${DATA_HOME}/.monitoring.ready