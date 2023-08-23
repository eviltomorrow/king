#!/bin/bash

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

function invalid(){
    echo "Error: wrong name '${args_1}'. Support: [apps、db、monitoring] !"
    exit 1
}

function clear_data(){
    sudo rm -rf ${1}
    check0 $? "rm -rf ${1}"
}

function docker_compose(){
    cd ${root_dir}
    check0 $? "cd ${root_dir}"
    if [ -d ${1} ]; then
        cd ${1}
        check0 $? "cd ${1}"
        if [ ! -f "./data/.ready" ]; then
            ./init.sh
        fi
        docker compose ${2} ${3}
        check0 $? "docker compose ${2} ${3}"
    else
        echo "Error: wrong file_path ${1} !"
        exit 1
    fi
}


exist_network=$(docker network ls -f name=net-king -q)

if [ ! -n "${exist_network}" ]; then
    docker network create net-king > /dev/null
    check0 $? "docker network create net-king"
fi

root_dir=$(pwd)
name=""
action=""

support_name=("db" "apps" "monitoring")
support_action=("up" "down" "clear")


case $# in
0)
    help;
    ;;

1)
    args_1=$1
    for v in ${support_action[@]}
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
    for v in ${support_name[@]}
    do
        if [ "${args_1}" = "${v}" ]; then
            name=${args_1}
            break;
        fi
    done
    if [ ! -n "${name}" ]; then
        invalid ${args_1}
    fi

    for v in ${support_action[@]}
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

export APPS_HOME=${root_dir}/apps
case ${action} in 
'clear')
    if [ -n "${name}" ]; then
        prefix=""
        if [ ${name} != "apps" ]; then
            prefix="components/"
        fi
        clear_data "./${prefix}${name}/data" 
    else
        for v in ${support_name[@]}
        do
            prefix=""
            if [ ${v} != "apps" ]; then
                prefix="components/"
            fi
            clear_data "./${prefix}${v}/data" 
        done
    fi
    ;;

'up')
    if [ -n "${name}" ]; then
        prefix=""
        if [ "${name}" != "apps" ]; then
            prefix="components/"
        fi
        docker_compose "./${prefix}${name}" "up" "-d"
    else
        for v in ${support_name[@]}
        do
            prefix=""
            if [ "${v}" != "apps" ]; then
                prefix="components/"
            fi
            docker_compose "./${prefix}${v}" "up" "-d"
        done
    fi
    ;;

'down')
    if [ -n "${name}" ]; then
        prefix=""
        if [ "${name}" != "apps" ]; then
            prefix="components/"
        fi
        docker_compose "./${prefix}${name}" "down"
    else
        for v in ${support_name[@]}
        do
            prefix=""
            if [ "${v}" != "apps" ]; then
                prefix="components/"
            fi
            docker_compose "./${prefix}${v}" "down"
        done
    fi
    ;;

*)
    ;;

esac
