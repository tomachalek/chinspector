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

	"chinspector/logrecord/raw"
	"chinspector/writer/influx"
)

var (
	// Block validation time: 2.41 seconds, pre_validation time: 0.21 seconds, cost: 1477974703, percent full: 13.436%
	blockValidTimeRgOld = regexp.MustCompile(
		`Block validation time: (\d*\.)(\d+) seconds, pre_validation time: (\d*\.)(\d+) seconds, cost: (\d+)|(None), percent full: (\d*\.)(\d+)%`)
	// Added unfinished_block 7d926f, not farmed by us, SP: 26 farmer response time: 4.5475, Pool pk xch1, validation time: 0.0147 seconds, pre_validation time 0.0148, cost: 198166855, percent full: 1.802%
	blockValidTimeRg   = regexp.MustCompile(`Added unfinished_block [0-9a-f]+, .+ farmer response time: (\d*\.)(\d+), .+, validation time: (\d*\.)(\d+) seconds, pre_validation time (\d*\.)(\d+), cost: (\d+)|(None), percent full: (\d*\.)(\d+)%`)
	blockchainHeightRg = regexp.MustCompile(
		`Updated peak to height (\d+), weight (\d+), hh`)
)

// ----

func ProcessRecord(rec *raw.Record, ch chan<- influx.InfluxRecord) error {
	if rec.Service != "full_node" {
		return fmt.Errorf("invalid record type used to process full_node record: %s", rec.Service)
	}

	ans := blockValidTimeRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		farmRespTime, err := strconv.ParseFloat(ans[1]+ans[2], 64)
		if err != nil {
			return err

		}
		blockValTime, err := strconv.ParseFloat(ans[3]+ans[4], 64)
		if err != nil {
			return err

		}
		preValTime, err := strconv.ParseFloat(ans[5]+ans[6], 64)
		if err != nil {
			return err
		}
		cost, err := strconv.ParseInt(ans[7], 10, 64)
		if err != nil {
			return err
		}
		ch <- &BlockChainValidationTimeRec{
			ts:                 rec.Datetime,
			validationTime:     blockValTime,
			preValidationTime:  preValTime,
			farmerResponseTime: farmRespTime,
			cost:               cost,
		}
		return nil
	}

	ans = blockchainHeightRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		height, err := strconv.ParseInt(ans[1], 10, 64)
		if err != nil {
			return err
		}
		ch <- &BlockChainHeightRec{
			ts:     rec.Datetime,
			height: height,
		}

	}

	return nil
}
