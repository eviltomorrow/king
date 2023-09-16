#!/bin/bash

# version
export opentelemetry_collector_contrib_version="0.85.0"
export jaeger_version="1.49"
export prometheus_version="2.47.0"
export node_exporter_version="1.6.1"
export grafana_version="10.0.5"
export DATA_HOME="/home/shepard/data/king/data"

mkdir -p ${DATA_HOME}/prometheus/data
mkdir -p ${DATA_HOME}/otel-collector
mkdir -p ${DATA_HOME}/cassandra/{db,log}
mkdir -p ${DATA_HOME}/loki/{log,data}
mkdir -p ${DATA_HOME}/promtail/log
mkdir -p ${DATA_HOME}/grafana/{data,log,provisioning}
mkdir -p ${DATA_HOME}/grafana/provisioning/{dashboards,datasources}

chmod -R 777 ${DATA_HOME}/prometheus
chmod -R 777 ${DATA_HOME}/otel-collector
chmod -R 777 ${DATA_HOME}/cassandra
chmod -R 777 ${DATA_HOME}/loki
chmod -R 777 ${DATA_HOME}/promtail
chmod -R 777 ${DATA_HOME}/grafana

touch ${DATA_HOME}/.monitoring.ready