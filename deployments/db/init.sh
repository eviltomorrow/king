#!/bin/bash

mkdir -p ${DATA_HOME}/mongo/{db,log}
mkdir -p ${DATA_HOME}/mysql/{db,log,run}
mkdir -p ${DATA_HOME}/mysql/run/{mysqld}
mkdir -p ${DATA_HOME}/etcd
mkdir -p ${DATA_HOME}/redis

chmod -R 777 ${DATA_HOME}/mongo
chmod -R 777 ${DATA_HOME}/mysql
chmod -R 777 ${DATA_HOME}/etcd
chmod -R 777 ${DATA_HOME}/redis

touch ${DATA_HOME}/.db.ready