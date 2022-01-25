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

package harvester

import (
	"fmt"
	"time"
)

type EligiblePlotsRec struct {
	numEligible int
	foundProofs int
	procTime    float64
	ts          time.Time
}

func (e *EligiblePlotsRec) Tags() map[string]string {
	return map[string]string{
		"type": "harvester",
	}
}

func (e *EligiblePlotsRec) Fields() map[string]interface{} {
	return map[string]interface{}{
		"eligiblePlots":    e.numEligible,
		"eligibleProcTime": e.procTime,
		"foundProofs":      e.foundProofs,
	}
}

func (e *EligiblePlotsRec) Time() time.Time {
	return e.ts
}

func (e *EligiblePlotsRec) Measurement() string {
	return "harvester"
}

func (e *EligiblePlotsRec) String() string {
	return fmt.Sprintf("EligiblePlotsRec{ts: %v, eligiblePlots: %d, foundProofs: %d, eligibleProcTime: %v}",
		e.ts, e.numEligible, e.foundProofs, e.procTime)
}
