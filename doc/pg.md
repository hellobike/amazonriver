# PostgreSQL 配置

## 参数修改

```
wal_level = 'logical';
max_replication_slots = 5; #该值要大于1
```

**修改后需要重启才能生效**

## 创建有replication权限的用户

```sql
CREATE ROLE test_rep LOGIN  ENCRYPTED PASSWORD 'xxxx' REPLICATION;
GRANT CONNECT ON DATABASE test_database to test_rep;
```

## 修改白名单配置

在 pg_hba.conf 中增加配置:   ```host replication test_rep all md5```

**修改后需要reload才能生效**
