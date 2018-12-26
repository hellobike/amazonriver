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

package handler

import (
	"context"
	"time"

	"github.com/hellobike/amazonriver/conf"
	"github.com/hellobike/amazonriver/handler/output"
	"github.com/hellobike/amazonriver/model"
	"github.com/hellobike/amazonriver/util"
	"github.com/prometheus/client_golang/prometheus"
)

type handlerWrapper struct {
	output    output.Output
	dataCh    chan []*model.WalData
	datas     []*model.WalData
	maxPos    uint64
	callback  PosCallback
	sub       *conf.Subscribe
	rules     []*conf.Rule
	ruleCache map[string]*conf.Rule
	skipCache map[string]struct{}
	cancel    context.CancelFunc
	done      chan struct{}

	successcounter prometheus.Counter
	errcounter     prometheus.Counter
}

func (h *handlerWrapper) runloop(ctx context.Context) {
	defer close(h.done)

	timer := time.NewTimer(time.Second)
	for {
		var needflush bool

		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			needflush = true
		case datas := <-h.dataCh:
			for _, data := range datas {
				if data.Pos > h.maxPos {
					h.maxPos = data.Pos
				}

				if rule, matched := h.filterData(data); matched {
					data.Rule = rule
					h.datas = append(h.datas, data)
				}

			}
			needflush = len(h.datas) >= 20000
		}

		if needflush {
			h.flush()
			resetTimer(timer, time.Second)
		}
	}
}

func resetTimer(t *time.Timer, d time.Duration) {
	// reset timer
	select {
	case <-t.C:
	default:
	}
	t.Reset(d)
}

func (h *handlerWrapper) flush() (err error) {
	defer func() {
		if len(h.datas) > 0 {
			if err != nil {
				h.errcounter.Inc()
			} else {
				h.successcounter.Add(float64(len(h.datas)))
			}
		}

		h.callback(h.maxPos)
		h.datas = nil
	}()

	if len(h.datas) == 0 {
		return nil
	}

	if err := h.output.Write(h.datas...); err != nil {
		return err
	}
	return nil
}

func (h *handlerWrapper) filterData(data *model.WalData) (matchedRule *conf.Rule, matched bool) {
	if len(data.Data) == 0 {
		return
	}

	if _, skip := h.skipCache[data.Table]; skip {
		return
	}

	matchedRule, matched = h.ruleCache[data.Table]
	if !matched {
		for _, rule := range h.rules {
			if util.MatchSimple(rule.Table, data.Table) {
				matched = true
				matchedRule = rule
				break
			}
		}

		if !matched {
			h.skipCache[data.Table] = struct{}{}
			return
		}
		h.ruleCache[data.Table] = matchedRule
		return
	}
	return
}

func (h *handlerWrapper) Handle(datas ...*model.WalData) error {

	h.dataCh <- datas
	return nil
}

func (h *handlerWrapper) Stop() {
	h.cancel()
	<-h.done
	h.output.Close()
}
