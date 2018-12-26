# amazonriver

amazonriver 是一个将postgresql的实时数据同步到es或kafka的服务

## 原理

amazonriver 利用pg内部的逻辑复制功能,通过在pg创建逻辑复制槽,接收数据库的逻辑变更,通过解析test_decoding特定格式的消息,得到逻辑数据.

## PG 配置

PG的配置[pg配置](./doc/pg.md)

## amazonriver 配置

## 监控

amazonriver支持使用prometheus来监控同步数据状态

## 性能测试

blabla

## 许可

amazonriver 使用 Apache License 2 许可