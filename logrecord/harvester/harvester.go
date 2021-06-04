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
	"log"
	"regexp"
	"strconv"

	"github.com/tomachalek/chinspector/logrecord/raw"
	"github.com/tomachalek/chinspector/writer/influx"
)

var (
	loadedPlotsRg   = regexp.MustCompile("Loaded a total of (\\d+) plots of size (\\d*\\.)(\\d+) (TiB|PiB|EiB), in (\\d*\\.)(\\d+) seconds")
	eligiblePlotsRg = regexp.MustCompile("(\\d+) plots were eligible for farming [0-9a-f]+\\.\\.\\. Found (\\d+) proofs. Time: (\\d*\\.)(\\d+) s. Total (\\d+) plots")
)

func ProcessRecord(rec *raw.Record, ch chan<- influx.InfluxRecord) error {
	ans := loadedPlotsRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		var err error
		numPlots, err := strconv.Atoi(ans[1])
		if err != nil {
			log.Print("ERROR: failed to parse number of plots")
		}
		plotSize, err := strconv.ParseFloat(ans[2]+ans[3], 64)
		if err != nil {
			log.Print("ERROR: Failed to parse plots size")
		}
		// normalize size to TiB
		switch ans[4] {
		case "GiB":
			plotSize /= 1 << 10
		case "PiB":
			plotSize *= 1 << 10
		case "EiB":
			plotSize *= 1 << 20
		}
		time, err := strconv.ParseFloat(ans[5]+ans[6], 64)
		if err != nil {
			log.Print("ERROR: failed to parse time information")
		}
		ch <- &TotalPlotsRec{
			ts:         rec.Datetime,
			totalPlots: numPlots,
			sizeTiB:    int(plotSize),
			procTime:   time,
		}
		return nil
	}

	ans = eligiblePlotsRg.FindStringSubmatch(rec.Message)
	if len(ans) > 0 {
		numEligible, err := strconv.Atoi(ans[1])
		if err != nil {
			log.Print("ERROR: ", err)
			// TODO
		}
		procTime, err := strconv.ParseFloat(ans[3]+ans[4], 64)
		if err != nil {
			log.Print("ERROR: ", err)
			// TODO
		}
		rec := &EligiblePlotsRec{
			numPlots: numEligible,
			procTime: procTime,
			ts:       rec.Datetime,
		}
		fmt.Println("    ", rec)
		ch <- rec
		return nil
	}
	return nil

}
