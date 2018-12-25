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
	"fmt"

	"github.com/hellobike-com/amazonriver/conf"
	"github.com/hellobike-com/amazonriver/model"
)

// newFakeHandler create a Handler print all data
func newFakeHandler(_ *conf.Subscribe) Output {
	return &fakeHandler{}
}

type fakeHandler struct{}

func (l *fakeHandler) Write(datas ...*model.WalData) error {
	for _, data := range datas {
		// TODO: update print
		fmt.Printf("TYPE:%s SCHEME:%s TABLE:%s DATA:%#v\n", data.OperationType.String(), data.Schema, data.Table, data.Data)
	}
	return nil
}

func (l *fakeHandler) Close() {}
