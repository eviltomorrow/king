#!/bin/bash

exist_network=$(docker network ls -f name=net-king -q)

if [ ! -n "${exist_network}" ]; then
    docker network create net-king
fi

# ./setup.sh ${name} up
# ./setup.sh ${name} down
# ./setup.sh ${name} clear

function help(){
    echo -e "eg. \r\n ./setup.sh [name] up \r\n ./setup.sh [name] down"
    exit 1
}

function invalid(){
    echo "Error: wrong name '${args_1}'. Support: [apps、db、monitoring] !"
    exit 1
}

function clear_data(){
    
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

support_name=("apps" "db" "monitoring")
support_action=("up" "down" "clear")


case $# in
0)
    help;
    ;;

1)
    args_1=$1
    for v in ${support_action[@]}
    do
        if [ ${args_1} = ${v} ]; then
            action=${args_1}
            break;
        fi
    done
    if [ ! -n "${action}" ]; then
        help;
    fi
    ;;

2)
    args_1=${1}
    args_2=${2}
    for v in ${support_name[@]}
    do
        if [ ${args_1} = ${v} ]; then
            name=${args_1}
            break;
        fi
    done
    if [ ! -n "${name}" ]; then
        invalid ${args_1}
    fi

    for v in ${support_action[@]}
    do
        if [ ${args_2} = ${v} ]; then
            action=${args_2}
            break;
        fi
    done
    if [ ! -n "${action}" ]; then
        help;
    fi

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
    docker_compose "./components/db" ${action} ${is_d}
    ;;

'monitoring')
    export APPS_HOME=${root_dir}/apps
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./components/monitoring" ${action} ${is_d}
    ;;

'')
    export APPS_HOME=${root_dir}/apps
    is_d=""
    if [ ${action} = "up" ]; then
        is_d="-d"
    fi
    docker_compose "./components/db" ${action} ${is_d}
    docker_compose "./components/monitoring" ${action} ${is_d}
    docker_compose "./apps" ${action} ${is_d}
    ;;

*)
    invalid ${args_1}
    ;;
esac


