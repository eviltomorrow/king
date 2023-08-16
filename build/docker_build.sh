#!/bin/bash

root_dir=$(pwd)
app_dir=${root_dir}/app

function check0(){
    if [ "0" != ${1} ]; then
        echo -e "\033[34m=> Build Failure\033[0m"
        exit 1
    fi
}

registry="registry.cn-beijing.aliyuncs.com"
account="eviltomorrow"
for name in $(ls ${app_dir}); do
    echo -e "\033[32m=> Building docker image(${name})...\033[0m"
    docker build --target prod -t ${registry}/${account}/${name} . --build-arg APPNAME=${name} --build-arg MAINVERSION=${1} --build-arg GITSHA=${2} --build-arg BUILDTIME=${3}
    check0 ${?}
    echo -e "\033[32m=> Build Success\033[0m"
done

