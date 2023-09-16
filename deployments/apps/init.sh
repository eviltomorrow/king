#!/bin/bash

export king_collector_version="latest"
export king_email_version="latest"
export king_storage_version="latest"
export DATA_HOME="/home/shepard/data/king/data"

mkdir -p ${DATA_HOME}/king_collector/log
mkdir -p ${DATA_HOME}/king_email/log
mkdir -p ${DATA_HOME}/king_storage/log

chmod -R 777 ${DATA_HOME}/king_collector
chmod -R 777 ${DATA_HOME}/king_email
chmod -R 777 ${DATA_HOME}/king_storage

touch ${DATA_HOME}/.apps.ready