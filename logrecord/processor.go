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

package logrecord

import (
	"log"
	"time"

	"chinspector/config"
	"chinspector/logrecord/fullnode"
	"chinspector/logrecord/harvester"
	"chinspector/logrecord/raw"
	"chinspector/notify"
	"chinspector/writer/influx"
)

type Processor struct {
	checkInterval  time.Duration
	numErrorsAlarm int
	tzOffset       string
	filePath       string
	currRecord     *raw.PreparsedRecord
	writeChannel   chan<- influx.InfluxRecord
	writer         *influx.RecordWriter
	lastCheckTS    time.Time
	currCheckTS    time.Time
	lastHealtyTS   time.Time
	errCounts      map[string][]time.Time // key = measurement (e.g. "harvester")
	mailConf       *config.EmailNotification
}

func (proc *Processor) OnCheckStart(ts time.Time) {
	proc.lastCheckTS = proc.currCheckTS
	proc.currCheckTS = ts

}

func (proc *Processor) OnCheckStop() {

}

func (proc *Processor) OnQuit() {
	close(proc.writeChannel)
}

func (proc *Processor) processFullNodeRec(rec *raw.Record) {

}

func (proc *Processor) OnLineRead(item string) {
	prec := raw.ExtractDate(item, proc.tzOffset)
	if prec.Dt.IsZero() {
		if proc.currRecord != nil {
			proc.currRecord.Text.WriteString(prec.Rest)

		} else {
			// ERROR TODO
		}

	} else {
		if proc.currRecord != nil {
			proc.parseCurrentRecord()
			proc.currRecord = nil
		}

		proc.currRecord = &raw.PreparsedRecord{
			Dt: prec.Dt,
		}
		_, err := proc.currRecord.Text.WriteString(prec.Rest)
		if err != nil {
			// TODO
		}
	}
}

func (proc *Processor) tsRangeIsInAlarmLimit(t1 time.Time, t2 time.Time) bool {
	return t2.Sub(t1) < proc.checkInterval
}

func (proc *Processor) parseCurrentRecord() {
	if time.Since(proc.lastHealtyTS) > 5*time.Minute {
		go func() {
			err := notify.SendNotification(
				proc.mailConf,
				"Chinspector ALARM - node out of sync",
				"Check the Node, Dude!!! Node is likely fucked Man!!!",
			)
			if err != nil {
				log.Print("ERROR: failed to send notification mail: ", err)
			}
		}()
	}
	rec := raw.ParseRecord(proc.currRecord)
	if rec == nil {
		log.Printf("WARNING: unknown record: %v", rec)
		return
	}

	if rec.Level == "ERROR" {
		if len(proc.errCounts[rec.Service]) == 0 {
			proc.errCounts[rec.Service] = make([]time.Time, 0, proc.numErrorsAlarm)
		}
		proc.errCounts[rec.Service] = append(proc.errCounts[rec.Service], time.Now())
		numRecs := len(proc.errCounts[rec.Service])
		if numRecs >= proc.numErrorsAlarm && proc.tsRangeIsInAlarmLimit(
			proc.errCounts[rec.Service][numRecs-1],
			proc.errCounts[rec.Service][0],
		) {
			// TODO run alarm
		}
		proc.writeChannel <- NewErrLineRec(
			rec.Service,
			rec.Datetime,
			rec.Message,
		)

	} else {
		switch rec.Service {
		case "harvester":
			harvester.ProcessRecord(rec, proc.writeChannel)
		case "full_node":
			fullnode.ProcessRecord(rec, proc.writeChannel)
			proc.lastHealtyTS = time.Now()
		case "wallet":
			// TODO
		}
	}
}

func (proc *Processor) CheckIntervalSec() int {
	return int(proc.checkInterval.Seconds())
}

func (proc *Processor) FilePath() string {
	return proc.filePath
}

func NewProcessor(conf *config.Props) (*Processor, error) {
	wch := make(chan influx.InfluxRecord)
	writer, err := influx.NewRecordWriter(&conf.InfluxDB, wch)
	if err != nil {
		return nil, err
	}

	return &Processor{
		checkInterval:  time.Duration(conf.CheckIntervalSec) * time.Second,
		numErrorsAlarm: conf.NumErrorsAlarm,
		filePath:       conf.ChiaLogPath,
		tzOffset:       conf.TimeZoneOffset,
		writeChannel:   wch,
		writer:         writer,
		errCounts:      make(map[string][]time.Time),
	}, nil
}
