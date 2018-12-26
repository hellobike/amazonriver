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
	"net/http"
	"sync"

	"github.com/hellobike/amazonriver/conf"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// this is a static check
var _ Interface = (*river)(nil)

// Interface of river
type Interface interface {
	Start() error
	Stop()
	Update(config *conf.Conf)
}

type river struct {
	conf   *conf.Conf
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

// New create river from conf
func New(conf *conf.Conf) Interface {
	return &river{conf: conf}
}

// Start flow the river
func (r *river) Start() error {
	r.wg = new(sync.WaitGroup)
	r.ctx, r.cancel = context.WithCancel(context.Background())
	if r.conf != nil {
		for _, sub := range r.conf.Subscribes {
			r.wg.Add(1)
			stream := newStream(r.conf.PgDumpExec, sub)
			go stream.start(r.ctx, r.wg)
		}
	}

	if r.conf.PrometheusAddress != "" {
		// prometheus exporter
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe(r.conf.PrometheusAddress, nil)
	}

	return nil
}

func (r *river) Update(config *conf.Conf) {
	// stop running streams
	r.Stop()

	r.conf = config
	r.wg = new(sync.WaitGroup)
	r.ctx, r.cancel = context.WithCancel(context.Background())
	for _, sub := range config.Subscribes {
		r.wg.Add(1)
		stream := newStream(r.conf.PgDumpExec, sub)
		go stream.start(r.ctx, r.wg)
	}
}

func (r *river) Stop() {
	r.cancel()
	r.wg.Wait()
}
