#!/bin/bash

export DATA_HOME=""

function check0(){
    if [ "0" != "${1}" ]; then
        echo -e "\033[34m=> Execute Failure\033[0m: ${2}."
        exit 1
    fi
}

function help(){
    echo -e "eg. \r\n ./setup.sh [name] up \r\n ./setup.sh [name] down"
    exit 1
}

function support_check(){
    echo "Error: wrong name '${args_1}'. Support: [apps、db、monitoring] !"
    exit 1
}

function docker_compose_pull(){
    cd ${root_dir}
    check0 $? "cd ${root_dir}"   
    if [ -d ${1} ]; then
        cd ${1}
        check0 $? "cd ${1}"
        docker compose pull
        check0 $? "docker compose pull"
    else
        echo "Error: wrong path ${1} !"
        exit 1
    fi
}

function docker_compose_action(){
    cd ${root_dir}
    check0 $? "cd ${root_dir}"
    if [ -d ${1} ]; then
        cd ${1}
        check0 $? "cd ${1}"

        if [ ! -f "${DATA_HOME}/.${1}.ready" ]; then
            if [ "${2}"="up" ]; then
                ./init.sh
            fi
        fi
        docker compose ${2} ${3}
        check0 $? "docker compose ${2} ${3}"
    else
        echo "Error: wrong path ${1} !"
        exit 1
    fi
}

if [ ! -n "${DATA_HOME}" ]; then
    echo "=> No set DATA_HOME, please set value"
    exit 1
else
    if [ ! -d ${DATA_HOME} ]; then
        mkdir -p ${DATA_HOME}
        check0 $? "mkdir -p ${DATA_HOME}"
    fi
fi

exist_network=$(docker network ls -f name=net-king -q)
if [ ! -n "${exist_network}" ]; then
    docker network create net-king > /dev/null
    check0 $? "docker network create net-king"
fi

root_dir=$(pwd)

supported_name=("db" "apps" "monitoring")
supported_action=("up" "down" "pull")
up_ordering=("db" "apps" "monitoring")
down_ordering=("apps" "db" "monitoring") 

# version
export opentelemetry_collector_contrib_version="0.85.0"
export jaeger_version="1.49"
export prometheus_version="2.47.0"
export node_exporter_version="1.6.1"
export grafana_version="10.0.5"
export mongo_version="7.0.1"
export mysql_version="8.0.34"
export etcd_version="3.5.9"
export king_collector_version="latest"
export king_email_version="latest"
export king_storage_version="latest"

name=""
action=""

case $# in
0)
    help;
    ;;

1)
    args_1=$1
    for v in ${supported_action[@]}
    do
        if [ "${args_1}" = "${v}" ]; then
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
    for v in ${supported_name[@]}
    do
        if [ "${args_1}" = "${v}" ]; then
            name=${args_1}
            break;
        fi
    done
    if [ ! -n "${name}" ]; then
        support_check ${args_1}
    fi

    for v in ${supported_action[@]}
    do
        if [ "${args_2}" = "${v}" ]; then
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

case ${action} in 
'pull')
    read -r -p "=> Will pull image, Are you sure? [Y/n]: " input

    case ${input} in 
        [yY][eE][sS]|[yY])
            if [ -n "${name}" ]; then
                docker_compose_pull "${name}"
            else
                for v in ${supported_name[@]}
                do
                    docker_compose_pull "${v}"
                done
            fi
            ;;
        [nN][oO]|[nN])
            echo "abort"
            exit 1
            ;;
        *)
            echo "Error: invalid input"
            exit 1
            ;;
    esac
    ;;
'up')
    if [ -n "${name}" ]; then
        docker_compose_action "${name}" "up" "-d"
    else
        for v in ${up_ordering[@]}
        do
            docker_compose_action "${v}" "up" "-d"
        done
    fi
    ;;

'down')
    if [ -n "${name}" ]; then
        docker_compose_action "${name}" "down"
    else
        for v in ${down_ordering[@]}
        do
            docker_compose_action "${v}" "down"
        done
    fi
    ;;

*)
    ;;

esac
