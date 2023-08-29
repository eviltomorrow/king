#!/bin/bash

mkdir -p $(pwd)/data/prometheus/data
mkdir -p $(pwd)/data/otel-collector
mkdir -p $(pwd)/data/cassandra/{db,log}
mkdir -p $(pwd)/data/loki/{log,data}
mkdir -p $(pwd)/data/promtail/log
mkdir -p $(pwd)/data/grafana/{data,log,provisioning}
mkdir -p $(pwd)/data/grafana/provisioning/{dashboards,datasources}

chmod -R 777 $(pwd)/data

touch $(pwd)/data/.ready