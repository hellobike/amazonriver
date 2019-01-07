# amazonriver [![CircleCI](https://circleci.com/gh/hellobike/amazonriver.svg?style=svg)](https://circleci.com/gh/hellobike/amazonriver)

[![Go Report Card](https://goreportcard.com/badge/github.com/hellobike/amazonriver)](https://goreportcard.com/report/github.com/hellobike/amazonriver)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/hellobike/amazonriver?status.svg)](https://godoc.org/github.com/hellobike/amazonriver)
[![GitHub release](https://img.shields.io/github/release/hellobike/amazonriver.svg)](https://github.com/hellobike/amazonriver/releases)

amazonriver utilize the postgresql logical replication to sync realtime data into elasticsearch or kafka...

## architecture

![architecture](./doc/arch.png)

## Required

- Postgresql 9.4 or later
- Kafka 0.8 or later
- ElasticSearch 6.x

## Howto

### Install

```shell
$git clone https://github.com/hellobike/amazonriver
$cd amazonriver
$glide install
$go install
```

### Use amazonriver

    amazonriver -config config.json

## PG config

[pg config](./doc/pg.md)

## amazonriver config

### Monitor

amazonriver has built-in support for prometheus [How to config](./doc/prometheus.md)

### Elasticsearch

[sync to elastic search](./doc/es.md)

### Kafka

[sync to kafka](./doc/kafka.md)

## License

amazonriver released under Apache License 2