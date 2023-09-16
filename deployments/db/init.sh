#!/bin/bash

export mongo_version="7.0.1"
export mysql_version="8.0.34"
export etcd_version="3.5.9"
export DATA_HOME="/home/shepard/data/king/data"

mkdir -p ${DATA_HOME}/mongo/{db,log}
mkdir -p ${DATA_HOME}/mysql/{db,log,run}
mkdir -p ${DATA_HOME}/etcd

chmod -R 777 ${DATA_HOME}/mongo
chmod -R 777 ${DATA_HOME}/mysql
chmod -R 777 ${DATA_HOME}/etcd

touch ${DATA_HOME}/.db.ready