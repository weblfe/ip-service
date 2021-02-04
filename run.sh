#!/usr/bin/env bash

action="${1}"
# 分支
branch="master"
# osx-x86_64, win64, linux-aarch_64, linux-ppcle_64, linux-s390x, linux-x86_64
platform="linux-x86_64"
# get https://github.com/protocolbuffers/protobuf/releases/ all version
protoc_version="3.14.0"

# protoc save dir
protoc_save="/usr/local/bin"
# protoc url
protoc_download_url="https://github.com/protocolbuffers/protobuf/releases/download/v${protoc_version}/protoc-${protoc_version}-${platform}.zip"
# network
network="bongolive"

# docker version
if [ "${2}x" == "x" ];then
  export  SERVICE_VERSION=latest
else
  export SERVICE_VERSION="${2}"
fi

# dockerfile
if [ "${3}x" == "x" ];then
  export DOCKERFILE=Dockerfile
else
  export DOCKERFILE="${3}"
fi

# 安装tools
function install() {
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/tal-tech/go-zero/tools/goctl && \
    go get -u github.com/zeromicro/goctl-swagger && \
    curl --request GET -sL \
         --url "'${protoc_download_url}'"\
         --output "'${protoc_save}/protoc-${protoc_version}-${platform}.zip'" && \
    unzip "${protoc_save}/protoc-${protoc_version}-${platform}.zip" -d "${protoc_save}/protoc_${protoc_version}" && \
    rm -f "${protoc_save}/protoc-${protoc_version}-${platform}.zip" && \
    ln -s "${protoc_save}/protoc_${protoc_version}/bin/protoc" "/usr/local/bin/protoc"
}

function init() {
    if [ ! -f "./.docker/etc" ];then
      mkdir ./.docker/etc/ -p  && cp ./api/etc/ipService.yaml .docker/etc/ipService.yaml
    fi
    if [ "`docker network inspect "${network}" | grep "No such network:" |grep -v grep`x" == "x" ];then
      docker network create "${network}"
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

function dev() {
    go run api/ipService.go -f api/etc/ipService.yaml
}

# 更新代码
function push() {
    git add . && \
    git commit -am "`date +"%Y-%m-%d%H:%M.%S"` `git status -s`" && \
    git pull && \
    git push origin "${branch}" -u
}

function main(){
  case "$action" in
  "install")
    install
  ;;
  "dev")
    dev
  ;;
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
  "push")
  push
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
  echo "   dev   : run dev for api         "
  echo "   install  : install all tools    "
  echo "   push  : git push code           "
    ;;
  esac
}

main
