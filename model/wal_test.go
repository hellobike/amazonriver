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
	"reflect"
	"testing"

	"github.com/hellobike/amazonriver/conf"
)

func TestWalData_Reset(t *testing.T) {

	var empty = &WalData{}
	empty.Reset()

	type fields struct {
		OperationType Operation
		Schema        string
		Table         string
		Data          map[string]interface{}
		Timestamp     int64
		Pos           uint64
		Rule          *conf.Rule
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1",
			fields: fields{
				OperationType: Insert,
				Schema:        "test1",
				Table:         "test1_table",
				Data:          map[string]interface{}{"id": 1, "data": "test1"},
				Timestamp:     100,
				Pos:           100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &WalData{
				OperationType: tt.fields.OperationType,
				Schema:        tt.fields.Schema,
				Table:         tt.fields.Table,
				Data:          tt.fields.Data,
				Timestamp:     tt.fields.Timestamp,
				Pos:           tt.fields.Pos,
				Rule:          tt.fields.Rule,
			}
			d.Reset()

			if !reflect.DeepEqual(d, empty) {
				t.Errorf("after d.Reset, should equal to emtyp, got: %v", d)
			}
		})
	}
}

func TestNewWalData(t *testing.T) {
	var empty = &WalData{}
	empty.Reset()

	tests := []struct {
		name string
		want *WalData
	}{
		{
			name: "test1",
			want: empty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWalData()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWalData() = %v, want %v", got, tt.want)
			}
			PutWalData(got)
		})
	}
}

func BenchmarkNewWalData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := NewWalData()
		PutWalData(data)
	}
}
