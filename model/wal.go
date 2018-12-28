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

package model

import (
	"sync"

	"github.com/hellobike/amazonriver/conf"
)

// WalData represent parsed wal log data
type WalData struct {
	OperationType Operation
	Schema        string
	Table         string
	Data          map[string]interface{}
	Timestamp     int64
	Pos           uint64
	Rule          *conf.Rule
}

// Reset for reuse
func (d *WalData) Reset() {
	d.OperationType = Unknow
	d.Schema = ""
	d.Table = ""
	d.Data = nil
	d.Timestamp = 0
	d.Pos = 0
	d.Rule = nil
}

var waldatapool = sync.Pool{New: func() interface{} { return &WalData{} }}

// NewWalData get data from pool
func NewWalData() *WalData {
	data := waldatapool.Get().(*WalData)
	data.Reset()
	return data
}

// PutWalData putback data to pool
func PutWalData(data *WalData) {
	waldatapool.Put(data)
}
