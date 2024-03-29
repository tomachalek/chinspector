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

type BlockChainHeightRec struct {
	height int64
	ts     time.Time
}

func (e *BlockChainHeightRec) Tags() map[string]string {
	return map[string]string{}
}

func (e *BlockChainHeightRec) Fields() map[string]interface{} {
	return map[string]interface{}{
		"height": e.height,
	}
}

func (e *BlockChainHeightRec) Time() time.Time {
	return e.ts
}

func (e *BlockChainHeightRec) Measurement() string {
	return "full_node"
}

func (e *BlockChainHeightRec) String() string {
	return fmt.Sprintf("BlockChainHeightRec{ts: %v, height: %d}",
		e.ts, e.height)
}
