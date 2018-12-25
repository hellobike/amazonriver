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

package util

import (
	"reflect"
	"testing"
)

func TestString2Bytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{s: "test1"},
			want: []byte("test1"),
		},
		{
			name: "test2",
			args: args{s: "The quick brown fox jumps over the lazy dog"},
			want: []byte("The quick brown fox jumps over the lazy dog"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String2Bytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String2Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes2String(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{b: []byte("test1")},
			want: "test1",
		},
		{
			name: "test2",
			args: args{b: []byte("The quick brown fox jumps over the lazy dog")},
			want: "The quick brown fox jumps over the lazy dog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes2String(tt.args.b); got != tt.want {
				t.Errorf("Bytes2String() = %v, want %v", got, tt.want)
			}
		})
	}
}
