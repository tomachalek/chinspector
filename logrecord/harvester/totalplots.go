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

type TotalPlotsRec struct {
	totalPlots int
	sizeTiB    int
	procTime   float64
	ts         time.Time
}

func (e *TotalPlotsRec) Tags() map[string]string {
	return map[string]string{
		"type": "harvester",
	}
}

func (e *TotalPlotsRec) Fields() map[string]interface{} {
	return map[string]interface{}{
		"totalPlots":       e.totalPlots,
		"plotSizeTiB":      e.sizeTiB,
		"plotSizeProcTime": e.procTime,
	}
}

func (e *TotalPlotsRec) Time() time.Time {
	return e.ts
}

func (e *TotalPlotsRec) Measurement() string {
	return "harvester"
}

func (e *TotalPlotsRec) String() string {
	return fmt.Sprintf("TotalPlotsRec{ts: %v, numPlots: %d, sizeTiB: %d, procTime: %03.f}",
		e.ts, e.totalPlots, e.sizeTiB, e.procTime)
}
