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

	"github.com/jackc/pgx"
)

func TestParse(t *testing.T) {

	tcs := []struct {
		name    string
		src     *pgx.WalMessage
		want    *WalData
		wantErr bool
	}{
		{
			name:    "DELETE",
			src:     &pgx.WalMessage{WalData: []byte("table public.source: DELETE: id[integer]:0"), WalStart: 11},
			want:    &WalData{OperationType: Delete, Pos: 11, Schema: "public", Table: "source", Data: map[string]interface{}{"id": int64(0)}},
			wantErr: false,
		},
		{
			name:    "INSERT",
			src:     &pgx.WalMessage{WalData: []byte("table public.source: INSERT: id[integer]:0"), WalStart: 11},
			want:    &WalData{OperationType: Insert, Pos: 11, Schema: "public", Table: "source", Data: map[string]interface{}{"id": int64(0)}},
			wantErr: false,
		},
		{
			name:    "INSERT MULTI COLUMNS",
			src:     &pgx.WalMessage{WalData: []byte("table public.source: INSERT: id[integer]:0 name[text]:'phil'"), WalStart: 11},
			want:    &WalData{OperationType: Insert, Pos: 11, Schema: "public", Table: "source", Data: map[string]interface{}{"id": int64(0), "name": "phil"}},
			wantErr: false,
		},

		{
			name:    "UPDATE",
			src:     &pgx.WalMessage{WalData: []byte("table public.source: UPDATE: id[integer]:0"), WalStart: 11},
			want:    &WalData{OperationType: Update, Pos: 11, Schema: "public", Table: "source", Data: map[string]interface{}{"id": int64(0)}},
			wantErr: false,
		},
		{
			name:    "UPDATE MULTI COLUMNS",
			src:     &pgx.WalMessage{WalData: []byte("table public.source: UPDATE: id[integer]:0 name[text]:'phil'"), WalStart: 11},
			want:    &WalData{OperationType: Update, Pos: 11, Schema: "public", Table: "source", Data: map[string]interface{}{"id": int64(0), "name": "phil"}},
			wantErr: false,
		},
		{
			name:    "BEGIN",
			src:     &pgx.WalMessage{WalData: []byte("BEGIN 1024"), WalStart: 11},
			want:    &WalData{OperationType: Begin, Pos: 11},
			wantErr: false,
		},
		{
			name:    "COMMIT",
			src:     &pgx.WalMessage{WalData: []byte("COMMIT 1024"), WalStart: 11},
			want:    &WalData{OperationType: Commit, Pos: 11},
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.src)
			if (err != nil) != tc.wantErr {
				t.Errorf("test %s wantErr %v, got %v", tc.name, tc.wantErr, err)
				return
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("test %s want %#v, got %#v", tc.name, tc.want, got)
				return
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	src := &pgx.WalMessage{WalData: []byte("table public.source: UPDATE: id[integer]:0 name[text]:'phil'"), WalStart: 11}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		got, err := Parse(src)
		if err != nil {
			b.Error(err)
			return
		}
		_ = got
	}
}
