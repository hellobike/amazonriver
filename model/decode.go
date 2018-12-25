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
	"strconv"
	"strings"

	"github.com/hellobike-com/amazonriver/util"
	"github.com/jackc/pgx"
	"github.com/nickelser/parselogical"
)

// Parse test_decoding format wal to WalData
func Parse(msg *pgx.WalMessage) (*WalData, error) {
	result := parselogical.NewParseResult(util.Bytes2String(msg.WalData))
	if err := result.Parse(); err != nil {
		return nil, err
	}
	var ret = NewWalData()

	var schema, table string
	if result.Relation != "" {
		i := strings.IndexByte(result.Relation, '.')
		if i < 0 {
			table = result.Relation
		} else {
			schema = result.Relation[:i]
			table = result.Relation[i+1:]
		}

		ret.Schema = schema
		ret.Table = table
	}
	ret.Pos = msg.WalStart
	switch result.Operation {
	case "INSERT":
		ret.OperationType = Insert
	case "UPDATE":
		ret.OperationType = Update
	case "DELETE":
		ret.OperationType = Delete
	case "BEGIN":
		ret.OperationType = Begin
	case "COMMIT":
		ret.OperationType = Commit
	}

	if len(result.Columns) > 0 {
		ret.Data = make(map[string]interface{}, len(result.Columns))
	}
	for key, column := range result.Columns {
		if column.Quoted {
			ret.Data[key] = column.Value
			continue
		}

		if column.Value == "null" {
			ret.Data[key] = nil
			continue
		}

		if val, err := strconv.ParseInt(column.Value, 10, 64); err == nil {
			ret.Data[key] = val
			continue
		}
		if val, err := strconv.ParseFloat(column.Value, 64); err == nil {
			ret.Data[key] = val
			continue
		}
		ret.Data[key] = column.Value
	}

	return ret, nil
}
