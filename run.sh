#!/usr/bin/env bash
action=${1}
function init() {
    if [ ! -f "./docker/etc" ];then
      mkdir .docker/etc/ -p  && cp ./api/etc/ipservice.yaml .docker/etc/ipservice.yaml
    fi
}
function start() {
    docker-compose up -d
}

function restart() {
    docker-compose restart
}

function clean() {
    docker-compose stop && docker-compose rm -f
}

function main(){
  case "$action" in
  "start")
    start
   ;;
  "restart")
   restart
  ;;
  "clean")
  clean
  ;;
  init)
  init
  ;;
  *)
  echo "run <command>"
  echo "command: "
  echo "   start : start service docker    "
  echo "   restart : restart service       "
  echo "   init  : init service config for docker "
  echo "   clean :  stop and service docker"
    ;;
  esac
}

main