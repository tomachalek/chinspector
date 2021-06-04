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
	"regexp"
	"strconv"

	"github.com/tomachalek/chinspector/logrecord/raw"
	"github.com/tomachalek/chinspector/writer/influx"
)

var (
	blockValidTimeRg         = regexp.MustCompile("Block validation time: (\\d*\\.)(\\d+), cost: (\\d+), percent full: (\\d*\\.)(\\d+)%")
	blockValidTimeNoneCostRg = regexp.MustCompile("Block validation time: (\\d*\\.)(\\d+), cost: None")
	blockchainHeightRg       = regexp.MustCompile("Updated peak to height (\\d+), weight (\\d+), hh")
)

// ----

func ProcessRecord(rec *raw.Record, ch chan<- influx.InfluxRecord) error {
	if rec.RecType != "full_node" {
		return fmt.Errorf("Invalid record type used to process full_node record: %s", rec.RecType)
	}

	ans := blockValidTimeRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		blockValTime, err := strconv.ParseFloat(ans[1]+ans[2], 64)
		if err != nil {
			// TODO
		}
		ch <- &BlockChainValidationTimeRec{
			ts:   rec.Datetime,
			time: blockValTime,
		}
		return nil
	}

	ans = blockValidTimeNoneCostRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		blockValTime, err := strconv.ParseFloat(ans[1]+ans[2], 64)
		if err != nil {
			// TODO
		}
		ch <- &BlockChainValidationTimeRec{
			ts:   rec.Datetime,
			time: blockValTime,
		}
		return nil
	}
	ans = blockchainHeightRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		height, err := strconv.ParseInt(ans[1], 10, 64)
		if err != nil {
			// TODO
		}
		ch <- &BlockChainHeightRec{
			ts:     rec.Datetime,
			height: height,
		}

	}

	return nil
}
