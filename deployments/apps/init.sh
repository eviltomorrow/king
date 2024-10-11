#!/bin/bash

mkdir -p ${DATA_HOME}/king-collector/log
mkdir -p ${DATA_HOME}/king-notification/log
mkdir -p ${DATA_HOME}/king-storage/log
mkdir -p ${DATA_HOME}/king-cron/log

chmod -R 777 ${DATA_HOME}/king-collector
chmod -R 777 ${DATA_HOME}/king-notification
chmod -R 777 ${DATA_HOME}/king-storage
chmod -R 777 ${DATA_HOME}/king-cron

touch ${DATA_HOME}/.apps.ready