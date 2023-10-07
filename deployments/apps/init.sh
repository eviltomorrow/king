#!/bin/bash

mkdir -p ${DATA_HOME}/king-collector/log
mkdir -p ${DATA_HOME}/king-email/log
mkdir -p ${DATA_HOME}/king-storage/log

chmod -R 777 ${DATA_HOME}/king-collector
chmod -R 777 ${DATA_HOME}/king-email
chmod -R 777 ${DATA_HOME}/king-storage

touch ${DATA_HOME}/.apps.ready