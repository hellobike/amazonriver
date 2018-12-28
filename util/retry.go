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

// WithRetry retry util job succeeded or retry count limit exceed
func WithRetry(retry int, job func() error) error {
	var retrys int
	for {
		retrys++
		if err := job(); err != nil {
			if retry == -1 || (retry > 0 && retrys <= retry) {
				continue
			}
			return err
		}
		return nil
	}
}
