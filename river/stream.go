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

package river

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hellobike/amazonriver/log"

	"github.com/hellobike/amazonriver/conf"
	"github.com/hellobike/amazonriver/dump"
	"github.com/hellobike/amazonriver/handler"
	"github.com/hellobike/amazonriver/model"
	"github.com/jackc/pgx"
)

type stream struct {
	// pg_dump 可执行文件的path
	pgDump string
	// dump 是否dump历史数据
	dump bool
	// 订阅配置
	sub *conf.Subscribe
	// 当前 wal 位置
	receivedWal uint64
	flushWal    uint64
	// 复制连接
	replicationConn *pgx.ReplicationConn
	// 消息处理
	handler handler.Handler
	// 取消
	cancel context.CancelFunc
	// ack 锁
	sendStatusLock sync.Mutex
	// buffered data
	datas []*model.WalData
}

func (s *stream) getReceivedWal() uint64 {
	return atomic.LoadUint64(&s.receivedWal)
}

func (s *stream) setReceivedWal(val uint64) {
	atomic.StoreUint64(&s.receivedWal, val)
}

func (s *stream) getFlushWal() uint64 {
	return atomic.LoadUint64(&s.flushWal)
}

func (s *stream) setFlushWal(val uint64) {
	atomic.StoreUint64(&s.flushWal, val)
}

func (s *stream) getStatus() (*pgx.StandbyStatus, error) {
	return pgx.NewStandbyStatus(s.getReceivedWal(), s.getFlushWal(), s.getFlushWal())
}

func newStream(pgDump string, sub *conf.Subscribe) *stream {
	var ret = &stream{pgDump: pgDump, sub: sub, dump: sub.Dump}
	ret.handler = handler.NewHandler(sub, ret.setFlushWal)

	return ret
}

func (s *stream) start(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	log.Logger.Infof("start stream for %s", s.sub.SlotName)
	ctx, s.cancel = context.WithCancel(ctx)

	config := pgx.ConnConfig{Host: s.sub.PGConnConf.Host, Port: s.sub.PGConnConf.Port, Database: s.sub.PGConnConf.Database, User: s.sub.PGConnConf.User, Password: s.sub.PGConnConf.Password}
	conn, err := pgx.ReplicationConnect(config)
	if err != nil {
		log.Logger.Errorf("create replication connection err: %v", err)
		return err
	}

	s.replicationConn = conn

	slotname := s.sub.SlotName

	_, snapshotID, err := conn.CreateReplicationSlotEx(slotname, "test_decoding")
	if err != nil {
		// 42710 means replication slot already exists
		if pgerr, ok := err.(pgx.PgError); !ok || pgerr.Code != "42710" {
			log.Logger.Errorf("create replication slot err: %v", err)
			return fmt.Errorf("failed to create replication slot: %s", err)
		}
	}

	_ = s.sendStatus()

	// Handle old data from db
	if err := s.exportSnapshot(snapshotID); err != nil {
		log.Logger.Errorf("export snapshot %s err: %v", snapshotID, err)
		return fmt.Errorf("slot name %s, err export snapshot: %v", s.sub.SlotName, err)
	}

	if err := conn.StartReplication(slotname, 0, -1); err != nil {
		log.Logger.Errorf("start replication err: %v", err)
		return err
	}

	return s.runloop(ctx)
}

func (s *stream) stop() error {
	s.cancel()
	s.handler.Stop()
	return s.replicationConn.Close()
}

func (s *stream) exportSnapshot(snapshotID string) error {
	// replication slot already exists
	if snapshotID == "" || !s.dump {
		return nil
	}
	dumper := dump.New(s.pgDump, s.sub)
	return dumper.Dump(snapshotID, s.handler)
}

func (s *stream) runloop(ctx context.Context) error {
	defer s.stop()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				_ = s.sendStatus()
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		msg, err := s.replicationConn.WaitForReplicationMessage(ctx)
		if err != nil {
			if err == ctx.Err() {
				return err
			}
			if err := s.checkAndResetConn(); err != nil {
				log.Logger.Errorf("reset replication connection err: %v", err)
			}
			continue
		}

		if msg == nil {
			continue
		}

		if err := s.replicationMsgHandle(msg); err != nil {
			log.Logger.Errorf("handle replication msg err: %v", err)
			continue
		}
	}
}

func (s *stream) checkAndResetConn() error {
	if s.replicationConn != nil && s.replicationConn.IsAlive() {
		return nil
	}

	time.Sleep(time.Second * 10)

	config := pgx.ConnConfig{
		Host:     s.sub.PGConnConf.Host,
		Port:     s.sub.PGConnConf.Port,
		Database: s.sub.PGConnConf.Database,
		User:     s.sub.PGConnConf.User,
		Password: s.sub.PGConnConf.Password,
	}
	conn, err := pgx.ReplicationConnect(config)
	if err != nil {
		return err
	}

	if _, _, err := conn.CreateReplicationSlotEx(s.sub.SlotName, "test_decoding"); err != nil {
		if pgerr, ok := err.(pgx.PgError); !ok || pgerr.Code != "42710" {
			return fmt.Errorf("failed to create replication slot: %s", err)
		}
	}

	if err := conn.StartReplication(s.sub.SlotName, 0, -1); err != nil {
		_ = conn.Close()
		return err
	}

	s.replicationConn = conn

	return nil
}

// ReplicationMsgHandle handle replication msg
func (s *stream) replicationMsgHandle(msg *pgx.ReplicationMessage) error {

	// 回复心跳
	if msg.ServerHeartbeat != nil {

		if msg.ServerHeartbeat.ServerWalEnd > s.getReceivedWal() {
			s.setReceivedWal(msg.ServerHeartbeat.ServerWalEnd)
		}
		if msg.ServerHeartbeat.ReplyRequested == 1 {
			_ = s.sendStatus()
		}
	}

	if msg.WalMessage != nil {

		logmsg, err := model.Parse(msg.WalMessage)
		if err != nil {
			return fmt.Errorf("invalid pgoutput msg: %s", err)
		}

		logmsg.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
		if err := s.handleMessage(logmsg); err != nil {
			return err
		}
	}

	return nil
}

func (s *stream) handleMessage(data *model.WalData) (err error) {
	log.Logger.Infof("handle wal data: %v", data)
	var needFlush bool
	switch data.OperationType {

	// 事务开始
	case model.Begin:
	// 	事务结束
	case model.Commit:
		needFlush = true
	default:
		s.datas = append(s.datas, data)
		// 防止大事务耗尽内存
		needFlush = len(s.datas) > 1000
	}

	if needFlush {
		_ = s.flush()
	}

	return nil
}

func (s *stream) flush() error {
	if len(s.datas) > 0 {
		_ = s.handler.Handle(s.datas...)
		s.datas = nil
	}
	return nil
}

// 发送心跳
func (s *stream) sendStatus() error {
	s.sendStatusLock.Lock()
	defer s.sendStatusLock.Unlock()

	log.Logger.Debug("send heartbeat")
	status, err := s.getStatus()
	if err != nil {
		return err
	}
	return s.replicationConn.SendStandbyStatus(status)
}
