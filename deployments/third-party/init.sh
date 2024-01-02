#!/bin/bash

mkdir -p ${DATA_HOME}/ntfy/{cache,lib,etc}

chmod -R 777 ${DATA_HOME}/ntfy

touch ${DATA_HOME}/.third-party.ready