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
	"context"
	"fmt"
	"strings"

	"github.com/hellobike-com/amazonriver/conf"
	"github.com/hellobike-com/amazonriver/model"
	"github.com/hellobike-com/amazonriver/util"

	"github.com/olivere/elastic"
)

type esHandler struct {
	client *elastic.Client
	sub    *conf.Subscribe
}

// newESOutput create handler write data to es
func newESOutput(sub *conf.Subscribe) Output {
	if sub.ESConf == nil {
		panic("es conf is nil")
	}

	client, err := elastic.NewClient(elastic.SetURL(strings.Split(sub.ESConf.Addrs, ",")...), elastic.SetBasicAuth(sub.ESConf.User, sub.ESConf.Password))
	if err != nil {
		panic(err)
	}

	handler := &esHandler{
		client: client,
		sub:    sub,
	}

	return handler
}

func (e *esHandler) Write(datas ...*model.WalData) error {
	numReqs := len(datas)
	if numReqs == 0 {
		return nil
	}

	var reqs = make([]elastic.BulkableRequest, 0, numReqs)
	for _, data := range datas {
		if req := e.makeRequest(data); req != nil {
			reqs = append(reqs, req)
		}
	}

	if len(reqs) == 0 {
		return nil
	}

	bulk := e.client.Bulk().Add(reqs...).Refresh("true")

	return util.WithRetry(e.sub.Retry, func() error {
		if _, err := bulk.Do(context.Background()); err != nil {
			// TODO: metric err
			return err
		}
		// TODO: metric succeed
		return nil
	})
}

func (e *esHandler) Close() {
	e.client.Stop()
}

func (e *esHandler) makeRequest(data *model.WalData) elastic.BulkableRequest {
	defer model.PutWalData(data)

	matchedRule := data.Rule

	var idBuf bytes.Buffer

	for _, field := range matchedRule.ESID {
		idBuf.WriteString(fmt.Sprint(data.Data[field]))
	}

	id := idBuf.String()
	if id == "" {
		return nil
	}

	switch data.OperationType {
	case model.Insert, model.Update:
		return elastic.NewBulkIndexRequest().
			Index(matchedRule.Index).
			Type(matchedRule.Type).
			Id(id).
			Doc(data.Data)
	case model.Delete:
		return elastic.NewBulkDeleteRequest().
			Index(matchedRule.Index).
			Type(matchedRule.Type).
			Id(id)
	}
	return nil
}
