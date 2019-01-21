# 同步到es配置

如下面的配置文件，配置了一个把 test.student_name 中的实时数据同步到es中student_name的索引中

```json
{
    # pg_dump 可执行文件path，如pg_dump在 $PATH 路径下面，则不需配置
    "pg_dump_path": "",
    "subscribes": [{
        # 是否dump 历史数据，如只需要实时数据，可以不配或配置为false，默认false
        "dump": false,
        # 逻辑复制槽名称，确保唯一
        "slotName": "slot_for_es",
        # pg 连接配置
        "pgConnConf": {
            "host": "127.0.0.1",
            "port": 5432,
            "database": "test",
            "user": "postgres",
            "password": "admin"
        },
        # 同步规则配置
        "rules": [
            {
                # 表名匹配，支持通配符
                "table": "student_name",
                # 表的主键配置
                "pks": ["id"],
                # es的id值配置，如配置 "id" 则会把表中的id字段作为es的_id
                "esid": ["id"],
                # es的索引配置
                "index": "student_name",
                # es的type配置
                "type": "logs"
            }
        ],
        # es 连接配置
        "esConf": {
            "addrs": "http://localhost:9200",
            "user": "",
            "password": ""
        },
        # 错误重试配置,0为不重试,-1会一直重试直到成功
        "retry": 0
    }],
    # 监控抓取地址配置
    "prometheus_address": ":8080"
}
```