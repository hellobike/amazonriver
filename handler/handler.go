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

	"github.com/hellobike/amazonriver/conf"
	"github.com/hellobike/amazonriver/handler/output"
	"github.com/hellobike/amazonriver/model"
)

// Handler handle dump data
type Handler interface {
	// Handle row data
	Handle(wal ...*model.WalData) error
	Stop()
}

// PosCallback for handler
type PosCallback func(uint64)

// NewHandler create wal handler with subscribe config
func NewHandler(sub *conf.Subscribe, callback PosCallback) Handler {
	ret := &handlerWrapper{
		dataCh:    make(chan []*model.WalData, 20480),
		callback:  callback,
		rules:     sub.Rules,
		sub:       sub,
		ruleCache: map[string]*conf.Rule{},
		skipCache: map[string]struct{}{},
		done:      make(chan struct{}),
	}

	ret.output = output.NewOutput(sub)

	ctx, cancel := context.WithCancel(context.Background())
	ret.cancel = cancel
	go ret.runloop(ctx)
	return ret
}
