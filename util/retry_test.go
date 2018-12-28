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
	"errors"
	"testing"
)

func TestWithRetry(t *testing.T) {
	if err := WithRetry(1, func() error {
		return nil
	}); err != nil {
		t.Error(err)
	}

	var i int
	if err := WithRetry(3, func() error {
		i++
		if i < 3 {
			return errors.New("less than 3")
		}
		return nil
	}); err != nil {
		t.Error(err)
	}

	var j int
	if err := WithRetry(3, func() error {
		j++
		if j < 5 {
			return errors.New("less than 5")
		}
		return nil
	}); err == nil {
		t.Errorf("errro should not be nil")
	}

}
