// Copyright 2021 Tomas Machalek <tomas.machalek@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fullnode

import (
	"fmt"
	"time"
)

type BlockChainValidationTimeRec struct {
	ts                 time.Time
	validationTime     float64
	preValidationTime  float64
	farmerResponseTime float64
	cost               int64
}

func (e *BlockChainValidationTimeRec) Tags() map[string]string {
	return map[string]string{}
}

func (e *BlockChainValidationTimeRec) Fields() map[string]interface{} {
	return map[string]interface{}{
		"blockValidationTime": e.validationTime,
		"preValidationTime":   e.preValidationTime,
		"farmerResponseTime":  e.farmerResponseTime,
		"cost":                e.cost,
	}
}

func (e *BlockChainValidationTimeRec) Time() time.Time {
	return e.ts
}

func (e *BlockChainValidationTimeRec) Measurement() string {
	return "full_node"
}

func (e *BlockChainValidationTimeRec) String() string {
	return fmt.Sprintf("BlockChainValidationTimeRec{ts: %v, blockValidationTime: %.3f, preValidationTime: %.3f, farmerResponseTime: %.3f}",
		e.ts, e.validationTime, e.preValidationTime, e.farmerResponseTime)
}
