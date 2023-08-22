#!/bin/bash

# exist_network=$(docker network ls -f name=net-kinga -q)

# if [ ! -n ${exist_network} ]; then
#     docker network create net-king
# fi

# ./setup.sh ${name} up
# ./setup.sh ${name} down

function help(){
    echo -e "eg. \r\n ./setup.sh [name] up \r\n ./setup.sh [name] down"
    exit 1
}

function invalid(){
    echo "Error: wrong name '${args_1}'. Support: [apps、db、monitoring] !"
    exit 1
}

function docker_compose(){
    cd ${root_dir}
    if [ -d ${1} ]; then
        cd ${1}
        if [ ! -d "./data" ]; then
            ./init.sh
        fi
        docker compose ${2} ${3}
    else
        echo "Error: wrong file_path ${1} !"
        exit 1
    fi
}

root_dir=$(pwd)
name=""
action=""

case $# in
0)
    help;
    ;;
1)
    args_1=$1
    if [ ${args_1} != "up" ]&&[ ${args_1} != "down" ]; then
       help;
    fi
    action=${args_1}
    ;;
2)
    args_1=$1
    args_2=$2
    if [ ${args_2} != "up" ]&&[ ${args_2} != "down" ]; then
       help;
    fi
    if [ ${args_1} != "apps" ]&&[ ${args_1} != "db" ]&&[ ${args_1} != "monitoring" ]; then
        invalid ${args_1}
        exit 1
    fi
    name=${args_1}
    action=${args_2}
    ;;
*)
    help;
    ;;
esac

case ${name} in
'apps')
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./apps" ${action} ${is_d}
    ;;
'db')
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./componets/db" ${action} ${is_d}
    ;;
'monitoring')
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./componets/monitoring" ${action} ${is_d}
    ;;
'')
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./componets/db" ${action} ${is_d}
    docker_compose "./componets/monitoring" ${action} ${is_d}
    docker_compose "./apps" ${action} ${is_d}
    ;;
*)
    invalid ${args_1}
    ;;
esac


