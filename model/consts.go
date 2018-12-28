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

// Operation type
type Operation uint8

const (
	// Insert operation
	Insert Operation = iota
	// Delete operation
	Delete
	// Update operation
	Update
	// Begin transaction
	Begin
	// Commit transaction
	Commit
	// Unknow operation
	Unknow
)

func (o Operation) String() string {
	switch o {
	case Insert:
		return "INSERT"
	case Delete:
		return "DELETE"
	case Update:
		return "UPDATE"
	case Begin:
		return "BEGIN"
	case Commit:
		return "COMMIT"
	}

	return "UNKNOW"
}
