# amazonriver [![CircleCI](https://circleci.com/gh/hellobike/amazonriver.svg?style=svg)](https://circleci.com/gh/hellobike/amazonriver)

[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/hellobike/amazonriver?status.svg)](https://godoc.org/github.com/hellobike/amazonriver)
![GitHub release](https://img.shields.io/github/release/hellobike/amazonriver.svg)

amazonriver 是一个将postgresql的实时数据同步到es或kafka的服务

## 版本支持

- Postgresql 9.4 or later
- Kafka 0.8 or later
- ElasticSearch 6.x

## 原理

amazonriver 利用pg内部的逻辑复制功能,通过在pg创建逻辑复制槽,接收数据库的逻辑变更,通过解析test_decoding特定格式的消息,得到逻辑数据

## 安装使用

### 安装

```shell
$git clone https://github.com/hellobike/amazonriver
$cd amazonriver
$glide install
$go install
```

### 使用

    amazonriver -configs config.json

## PG 配置

PG数据库需要预先开启逻辑复制[pg配置](./doc/pg.md)

## amazonriver 配置

### 监控

amazonriver支持使用prometheus来监控同步数据状态,[配置Grafana监控](./doc/prometheus.md)

### 同步到 elasticsearch

[同步到elasticsearch](./doc/es.md)

### 同步到 kafka

[同步到kafka](./doc/kafka.md)

## 许可

amazonriver 使用 Apache License 2 许可