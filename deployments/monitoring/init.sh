#!/bin/bash

mkdir -p $(pwd)/data/prometheus/{data,etc}
mkdir -p $(pwd)/data/otel-collector/etc
mkdir -p $(pwd)/data/cassandra/{db,log}
mkdir -p $(pwd)/data/loki/{etc,log,data}
mkdir -p $(pwd)/data/promtail/{etc,log}
mkdir -p $(pwd)/data/grafana/{data,log,crypto_data,provisioning}
mkdir -p $(pwd)/data/grafana/provisioning/{dashboards,datasources}
mkdir -p $(pwd)/data/otel-collector/etc

chmod -R 777 $(pwd)/data


# otel-collector-config.yaml
cat > $(pwd)/data/otel-collector/etc/otel-collector-config.yaml << EOF
receivers:
  otlp:
    protocols:
      grpc:

exporters:
  logging:

  jaeger:
    endpoint: jaeger-collector:14250
    tls:
      insecure: true

processors:
  batch:

extensions:
  health_check:
  pprof:
    endpoint: :1888
  zpages:
    endpoint: :55679

service:
  extensions: [pprof, zpages, health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [logging]
EOF

cat > $(pwd)/data/loki/etc/local-config.yaml << EOF
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9095
  grpc_server_max_recv_msg_size: 1073741824 #grpc最大接收消息值,默认4m
  grpc_server_max_send_msg_size: 1073741824 #grpc最大发送消息值,默认4m

ingester:
  lifecycler:
    address: loki
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    final_sleep: 0s
  chunk_idle_period: 5m
  chunk_retain_period: 30s
  max_transfer_retries: 0
  max_chunk_age: 20m #一个timeseries块在内存中的最大持续时间。如果timeseries运行的时间超过此时间,则当前块将刷新到存储并创建一个新块

schema_config:
  configs:
    - from: 2021-01-01
      store: boltdb
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 168h

storage_config:
  boltdb:
    directory: /opt/loki/index #存储索引地址
  filesystem:
    directory: /opt/loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  ingestion_rate_mb: 30 #修改每用户摄入速率限制,即每秒样本量,默认值为4M
  ingestion_burst_size_mb: 20 #修改每用户摄入速率限制,即每秒样本量,默认值为6M

chunk_store_config:
  max_look_back_period: 168h   #回看日志行的最大时间，只适用于即时日志

table_manager:
  retention_deletes_enabled: true #日志保留周期开关,默认为false
  retention_period: 168h #日志保留周期

EOF

cat > $(pwd)/data/promtail/etc/promtail-config.yaml << EOF
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: king-collector
    pipeline_stages:
    - json:
       expressions:
         level: level
         time: time
         caller: caller
         message: message
         error: error

    static_configs:
      - targets:
          - localhost
        labels:
          app: king-collector
          type: data 
          __path__: /log/king-collector/data.log

      - targets:
          - localhost
        labels:
          app: king-collector
          type: access 
          __path__: /log/king-collector/access.log

      - targets:
          - localhost
        labels:
          app: king-collector
          type: panic
          __path__: /log/king-collector/panic.log

  - job_name: king-email
    pipeline_stages:
    - json:
       expressions:
         level: level
         time: time
         caller: caller
         message: message
         error: error

    static_configs:
      - targets:
          - localhost
        labels:
          app: king-email
          type: data 
          __path__: /log/king-email/data.log

      - targets:
          - localhost
        labels:
          app: king-email
          type: access 
          __path__: /log/king-email/access.log

      - targets:
          - localhost
        labels:
          app: king-email
          type: panic 
          __path__: /log/king-email/panic.log

  - job_name: king-repository
    pipeline_stages:
    - json:
       expressions:
         level: level
         time: time
         caller: caller
         message: message
         error: error
    static_configs:
      - targets:
          - localhost
        labels:
          app: king-repository
          type: data 
          __path__: /log/king-repository/data.log

      - targets:
          - localhost
        labels:
          app: king-repository
          type: access 
          __path__: /log/king-repository/access.log

      - targets:
          - localhost
        labels:
          app: king-repository
          type: panic 
          __path__: /log/king-repository/panic.log

EOF


# prometheus.yaml
cat > $(pwd)/data/prometheus/etc/prometheus.yml << EOF
# my global config
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # By default, scrape targets every 15 seconds.
  # scrape_timeout is set to the global default (10s).

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
      monitor: 'king'

# Load and evaluate rules in this file every 'evaluation_interval' seconds.
rule_files:
  - 'alert.rules'

# alert
# alerting:
#   alertmanagers:
#   - scheme: http
#     static_configs:
#     - targets:
#       - "alertmanager:9093"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label \`job=<job_name>\` to any timeseries scraped from this config.

  - job_name: 'prometheus'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 15s

    static_configs:
         - targets: ['localhost:9090']

  - job_name: 'node-exporter'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 15s
  
    static_configs:
      - targets: ['node-exporter:9100']
EOF

cat > $(pwd)/data/prometheus/etc/alert.rules << EOF
groups:
- name: example
  rules:

  # Alert for any instance that is unreachable for >2 minutes.
  - alert: service_down
    expr: up == 0
    for: 2m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} down"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 2 minutes."

  - alert: high_load
    expr: node_load1 > 1.0
    for: 2m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} under high load"
      description: "{{ $labels.instance }} of job {{ $labels.job }} is under high load."
EOF

cat > $(pwd)/data/grafana/provisioning/datasources/datasource.yml << EOF
# config file version
apiVersion: 1

# list of datasources that should be deleted from the database
deleteDatasources:
  - name: Prometheus
    orgId: 1

# list of datasources to insert/update depending
# whats available in the database
datasources:
  # <string, required> name of the datasource. Required
- name: Prometheus
  type: prometheus
  access: proxy
  orgId: 1
  url: http://prometheus:9090
  password:
  user:
  database:
  basicAuth: false
  basicAuthUser:
  basicAuthPassword:
  withCredentials:
  isDefault: true
  jsonData:
     graphiteVersion: "1.1"
     tlsAuth: false
     tlsAuthWithCACert: false
  secureJsonData:
    tlsCACert: "..."
    tlsClientCert: "..."
    tlsClientKey: "..."
  version: 1
  editable: true

- name: Loki
  type: loki
  access: proxy 
  orgId: 1
  url: http://loki:3100
  basicAuth: false
  isDefault: false
  version: 1
  editable: true

EOF


cat > $(pwd)/data/grafana/provisioning/dashboards/dashboards.yml << EOF
apiVersion: 1

providers:
- name: 'Prometheus'
  orgId: 1
  folder: ''
  type: file
  disableDeletion: false
  editable: true
  options:
    path: /etc/grafana/provisioning/dashboards

EOF

cat > $(pwd)/data/grafana/provisioning/dashboards/dashboards.yml << EOF
apiVersion: 1

providers:
- name: 'Prometheus'
  orgId: 1
  folder: ''
  type: file
  disableDeletion: false
  editable: true
  options:
    path: /etc/grafana/provisioning/dashboards

EOF

touch $(pwd)/data/.ready