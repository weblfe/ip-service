#!/usr/bin/env bash

action=${1}

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
  *)
  echo "run <command>"
  echo "command: "
  echo "   start : start service docker    "
  echo "   restart : restart service       "
  echo "   clean :  stop and service docker"
    ;;
  esac
}

