# Ip 地理解析服务

> 数据来源 (开源项目)

https://github.com/flyaction/ipdatabase

> 服务使用go-zero 搭建

[官网go-zero](http://zero.gocn.vip/zero)

> 服务使用

```http

curl  http://micro.word-server.com/service/ip?ip=xxx.xxx.xxx.xxx

success response :

{
    "msg": "OK",
    "code": 0,
    "data": {
        "ip": "219.136.95.223",
        "city_id": 2140,
        "country": "中国",
        "district": "0",
        "province": "广东省",
        "city": "广州市",
        "ISP": "电信"
    }
}

error response :
{
    "msg": "ip format error",
    "code": 500,
    "data": null
}
```
> 快速使用

1. 安装 docker

2. 安装 docker-compose

3. 运行
```bash
   run.sh init  && run.sh start
```
4. 请求

```http
curl localhost:28000?ip=xxx.xxx.xxx.xxx
```

> 演示站点接口

```http request
https://api.word-server.com/service/ip
```

> 开源协议 MIT

欢迎 大家 fork