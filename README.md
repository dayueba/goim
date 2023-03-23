goim v2.0
==============

goim is an im server writen in golang.

## Features
 * Supports single push, multiple push and broadcasting
 * Supports one key to multiple subscribers (Configurable maximum subscribers count)
 * Supports heartbeats (Application heartbeats, TCP, KeepAlive, HTTP long pulling)
 * Supports authentication (Unauthenticated user can't subscribe)
 * Supports multiple protocols (WebSocket，TCP，HTTP）
 * Scalable architecture (Unlimited dynamic job and logic modules)
 * Asynchronous push notification based on Kafka
- 轻量级
- 高性能
- 纯go编写


## Architecture
![arch](./docs/arch.png)

1. comet / job / logic 支持多实例部署, 注册在 bilibili/discovery 中

![](https://tsingson.github.io/tech/assets/goim-architecture-002.png)

## Quick Start

### Build
```
    make build
```

### Run
```
    make run
    make stop

    // or
    nohup target/logic -conf=target/logic.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 2>&1 > target/logic.log &
    nohup target/comet -conf=target/comet.toml -region=sh -zone=sh001 -deploy.env=dev -weight=10 -addrs=127.0.0.1 2>&1 > target/logic.log &
    nohup target/job -conf=target/job.toml -region=sh -zone=sh001 -deploy.env=dev 2>&1 > target/logic.log &

```
### Environment
```
    env:
    export REGION=sh
    export ZONE=sh001
    export DEPLOY_ENV=dev

    supervisor:
    environment=REGION=sh,ZONE=sh001,DEPLOY_ENV=dev

    go flag:
    -region=sh -zone=sh001 deploy.env=dev
```
### Configuration
You can view the comments in target/comet.toml,logic.toml,job.toml to understand the meaning of the config.

### Dependencies
[Discovery](https://github.com/bilibili/discovery)

[Kafka](https://kafka.apache.org/quickstart)

这俩中间件如何替换？抽象成interface

- logic 向 kafka 投递消息
  - internal/logic/dao/kafka.go 里面有三个方法
- job 订阅 kafka 消息
  - 

## Document
[Protocol](./docs/protocol.png)

[中文](./README_cn.md)

## Examples
Websocket: [Websocket Client Demo](https://github.com/Terry-Mao/goim/tree/master/examples/javascript)


## LICENSE
goim is is distributed under the terms of the MIT License.
