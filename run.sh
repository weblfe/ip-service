#!/usr/bin/env bash
action=${1}

function init() {
    if [ ! -f "./.docker/etc" ];then
      mkdir ./.docker/etc/ -p  && cp ./api/etc/ipService.yaml .docker/etc/ipService.yaml
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

function docs() {
    if [ -f serviceGetIp.json ];then
        rm -f serviceGetIp.json
    fi
    goctl api plugin -plugin goctl-swagger="swagger -filename serviceGetIp.json" -api api/ipService.api -dir .
}

function update() {
      goctl api go  -api api/ipService.api -dir api/
}

function apiDocs() {
    docker run --rm -p 8081:8080 -e SWAGGER_JSON=/app/ipService.json -v $PWD:/usr/share/nginx/html/app swaggerapi/swagger-ui
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
  "init")
  init
  ;;
  "docs")
  docs
  ;;
  "apiDocs")
  apiDocs
  ;;
  "update")
  update
  ;;
  *)
  echo "run <command>"
  echo "command: "
  echo "   start : start service docker    "
  echo "   restart : restart service       "
  echo "   docs  : gen swagger docs        "
  echo "   apiDocs  : see swagger ui       "
  echo "   init  : init service config for docker "
  echo "   clean :  stop and service docker"
  echo "   update : update code for api    "
    ;;
  esac
}

main
