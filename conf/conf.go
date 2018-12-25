/*
 * Copyright 2018 Shanghai Junzheng Network Technology Co.,Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conf

// Conf ...
type Conf struct {
	// pg_dump 可执行文件路径
	PgDumpExec string `json:"pg_dump_path"`
	// Subs 订阅规则
	Subscribes []*Subscribe `json:"subscribes"`
}

// ESConf es 配置
type ESConf struct {
	// es 地址
	Addrs string `json:"addrs"`
	// User : es username
	User string `json:"user"`
	// Password : es password
	Password string `json:"password"`
}

// PGConnConf of pg
type PGConnConf struct {
	// 地址
	Host string `json:"host"`
	// 端口
	Port uint16 `json:"port"`
	// database
	Database string `json:"database"`
	// schema
	Schema string `json:"schema"`
	// user
	User string `json:"user"`
	// password
	Password string `json:"password"`
}

// Subscribe 订阅一个数据库中的表的wal，根据规则保存到es里相应的index，type中
type Subscribe struct {
	// Dump 创建复制槽成功后，是否 dump 历史数据
	Dump bool `json:"dump"`
	// 逻辑复制槽
	SlotName string `json:"slotName"`
	// pg 连接配置
	PGConnConf *PGConnConf `json:"pgConnConf"`
	// 订阅规则
	Rules []*Rule `json:"rules"`
	// ES 配置
	ESConf *ESConf `json:"esConf"`
	// KafkaConf kafka 配置
	KafkaConf *KafkaConf `json:"kafkaConf"`
	// Retry 重试次数 -1:无限重试
	Retry int `json:"retry"`
}

// KafkaConf kafka 配置
type KafkaConf struct {
	// Addrs kafka集群地址
	Addrs []string `json:"addrs"`
}

// Rule 同步规则
type Rule struct {
	// 订阅数据表，支持 ?* 匹配
	Table string `json:"table"`
	// 主键
	PKs []string `json:"pks"`

	// 下面几项同步到es中时需配置
	// 用作es中id的字段，多个字段内容会连在一起
	ESID []string `json:"esid"`
	// Index es中的idex
	Index string `json:"index"`
	// Type es中的type
	Type string `json:"type"`

	// 下面几项同步到kafka中时需要配置
	// Topic ...
	Topic string `json:"topic"`
}
