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

package output

import (
	"bytes"
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/hellobike/amazonriver/conf"
	"github.com/hellobike/amazonriver/model"
	"github.com/hellobike/amazonriver/util"

	"github.com/Shopify/sarama"
	"github.com/json-iterator/go"
)

type kafkaHandler struct {
	sub   *conf.Subscribe
	kafka sarama.SyncProducer
}

type kafkaMessage struct {
	Schema      string                 `json:"schema"`
	Table       string                 `json:"table"`
	Operation   string                 `json:"operation"`
	Data        map[string]interface{} `json:"data"`
	OperateTime int64                  `json:"operateTime"`
}

func newKafkaOutput(sub *conf.Subscribe) Output {
	config := sarama.NewConfig()
	config.Producer.Compression = sarama.CompressionSnappy  // Compress messages
	config.Producer.Flush.Frequency = 10 * time.Millisecond // Flush batches every 10ms
	config.Producer.Flush.MaxMessages = 1 << 29
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Return.Successes = true
	sarama.Logger = stdlog.New(os.Stdout, "[KAFKA]", stdlog.LstdFlags)
	config.Producer.Retry.Max = 3
	if sub.Retry > 0 {
		config.Producer.Retry.Max = sub.Retry
	}
	config.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(sub.KafkaConf.Addrs, config)
	if err != nil {
		panic(err)
	}

	ret := &kafkaHandler{
		kafka: producer,
		sub:   sub,
	}

	return ret
}

func fromData(data *model.WalData) (key []byte, topic string, msg *kafkaMessage) {
	defer model.PutWalData(data)

	if len(data.Data) == 0 {
		return nil, "", nil
	}

	var keyb bytes.Buffer
	var connector string
	for _, f := range data.Rule.PKs {
		fmt.Fprint(&keyb, connector)
		fmt.Fprintf(&keyb, "%v", data.Data[f])
		connector = "-"
	}

	return keyb.Bytes(), data.Rule.Topic, &kafkaMessage{
		Schema:      data.Schema,
		Table:       data.Table,
		Data:        data.Data,
		Operation:   data.OperationType.String(),
		OperateTime: data.Timestamp,
	}
}

func (k *kafkaHandler) Write(datas ...*model.WalData) error {

	if len(datas) == 0 {
		return nil
	}

	var msgs = make([]*sarama.ProducerMessage, 0, len(datas))
	for _, data := range datas {
		key, topic, msg := fromData(data)
		if msg == nil {
			continue
		}
		bts, err := jsoniter.Marshal(msg)
		if err != nil {
			// TODO: print err
			continue
		}

		m := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(bts),
		}

		msgs = append(msgs, m)
	}

	if len(msgs) == 0 {
		return nil
	}

	return util.WithRetry(k.sub.Retry, func() error {
		if err := k.kafka.SendMessages(msgs); err != nil {
			// TODO: print err
			return err
		}
		// TODO: metric succeed

		return nil
	})
}

func (k *kafkaHandler) Close() {
	k.kafka.Close()
}
