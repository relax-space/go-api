[![Build Status](https://travis-ci.org/relax-space/go-api.svg?branch=master)](https://travis-ci.org/relax-space/go-api)
[![codecov](https://codecov.io/gh/relax-space/go-api/branch/master/graph/badge.svg)](https://codecov.io/gh/relax-space/go-api)

# go-api template

Quickly create an echo-based api project

## Getting Started

### Create a new project
>Rename go-api to your project name

### Run
```bash
$ docker-compose -f .\example\docker-compose.yml up -d
$ go run .
```

Visit           http://127.0.0.1:8080/ping

Visit swagger   http://127.0.0.1:8080/docs

## View logs
1. download [kafka.apache](https://kafka.apache.org/downloads)
2. start consumer `./kafka-console-consumer.bat --bootstrap-server localhost:9092 --from-beginning --topic behaviorlog`
3. request a api,like: http://localhost:8080/v1/fruits
4. consumer will accept info

```json
{"action_id":"d14b340b-7b68-11ea-b68b-847beb364db6","service":"go-api","timestamp":"2020-04-11T04:21:10.657333+08:00","remote_ip":"::1","host":"localhost:8080","uri":"/v1/fruits","method":"GET","path":"/v1/fruits","referer":"http://localhost:8080/docs","user_agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36","status":200,"latency":29999000,"bytes_sent":67,"hostname":"SS-CN-XIAOXINMI","controller":"controllers.FruitAPIController","action":"GetAll"}
```

## Extensions

- auth: https://gitlab.p2shop.cn:8443/sample/go-api-auth 
- validator:[github](https://github.com/relax-space/go-api/go-api-validator) | https://gitlab.p2shop.cn:8443/sample/go-api-validator 
- windows: [github](https://github.com/relax-space/go-api/go-api-windows) | https://gitlab.p2shop.cn:8443/sample/go-api-windows.git 
- producer: [github](https://github.com/relax-space/go-api/go-api-producer) | https://gitlab.p2shop.cn:8443/sample/go-api-producer 
- consumer: [github](https://github.com/relax-space/go-api/go-api-consumer) | https://gitlab.p2shop.cn:8443/sample/go-api-consumer 
- twodb: [github](https://github.com/relax-space/go-api/go-api-twodb) | https://gitlab.p2shop.cn:8443/sample/go-api-twodb 
- sqlserver: [github](https://github.com/relax-space/go-api/go-api-sqlserver) | https://gitlab.p2shop.cn:8443/sample/go-api-sqlserver 
- postgres: [github](https://github.com/elandcloud/go-api-postgres) | https://gitlab.p2shop.cn:8443/sample/go-api-postgres


## References

- echosample: [echosample](https://github.com/pangpanglabs/echosample)


