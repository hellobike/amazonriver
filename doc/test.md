# 性能测试

## 测试环境

APP:

|类别|名称|
|---|---|
|OS       | centos 6.5 |
|CPU      |Common KVM CPU 8 CORE|
|RAM      |16GB|
|DISK     |100GB|

PG:

|类别|名称|
|---|---|
|pg    |9.4|
|OS       | centos 6.5 |
|CPU      |Common KVM CPU 16 CORE|
|RAM      |32GB|
|DISK     |500GB|

ES:

|类别|名称|
|---|---|
| es | 5.6.3 |
|OS       | centos 6.5 |
|CPU      |Common KVM CPU 8 CORE|
|RAM      |64GB|
|DISK     |1788G SSD|
|NODE     |3|

KAFKA:

|类别|名称|
|---|---|
| kafka | 0.9 |
|OS       | centos 6.5 |
|CPU      |Common KVM CPU 4 CORE|
|RAM      |8GB|
|DISK     | 500G|
|NODE     |2|

## 性能需求

- 测试wal没有积压的情况下，dml的tps峰值

## 测试用例

单个表约5个字段，每个事务insert单条数据，分别同步到es与kafka，在没有wal积压的情况下，测试峰值tps

## ES

1. 单事务单表同步，app cpu ≈ 20.08%，load 峰值1.03，mem ≈ 700MB；内网 net/io in ≈ 154.95Mb，net/io out ≈ 23.32Mb，ES先出现瓶颈；
2. 压测端压力峰值约 17k/s，DB tps ≈ 18k，es avg rate index ≈ 13k/s；
3. 单事务多表同步，app cpu、es index rate不会明显提升；
4. 补充4c8g服务器的资源消耗对比：cpu ≈ 28.86%，load峰值0.6，mem ≈ 48MB；内网 net/io in ≈ 114.77Mb，net/io out ≈ 18.7Mb。

## kafka

1. 单事务单表同步，cpu ≈ 28.15%，load 峰值0.4，mem ≈ 30MB；内网 net/io in ≈ 253.66Mb，net/io out ≈ 28.53Mb；压测端压力峰值约 23k/s，DB tps ≈ 23k，kafka rate ≈ 23k；
2. 单事务多表同步，cpu ≈ 35.38%，load 峰值0.09，mem ≈ 30MB；内网 net/io in ≈ 292.71Mb，net/io out ≈ 37.99Mb；压测端压力峰值约 21k/s，DB tps ≈ 40k，kafka rate ≈ 40k；
综述
1. 应用会过滤 slot配置的数据库内所有表的DML操作，消耗cpu ≈ 17%；
2. es 同步速率远逊于kafka 同步速率，前提是 应用服务器和DB服务没有瓶颈。在我们的测试环境，受限于基础配置，只能体现es和kafka部分性能差异；
3. 当日志级别为error时，同步过程中服务器的cpu、mem、net io消耗均比较低，mem占用取决于wal堆积的大小；

