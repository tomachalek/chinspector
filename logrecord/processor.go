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

	"github.com/tomachalek/chinspector/config"
	"github.com/tomachalek/chinspector/logrecord/fullnode"
	"github.com/tomachalek/chinspector/logrecord/harvester"
	"github.com/tomachalek/chinspector/logrecord/raw"
	"github.com/tomachalek/chinspector/writer/influx"
)

type Processor struct {
	checkIntervalSec int
	tzOffset         string
	filePath         string
	currRecord       *raw.PreparsedRecord
	writeChannel     chan<- influx.InfluxRecord
	writer           *influx.RecordWriter
}

func (proc *Processor) OnCheckStart() {

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

func (proc *Processor) parseCurrentRecord() {
	rec := raw.ParseRecord(proc.currRecord)
	if rec == nil {
		log.Printf("WARNING: unknown record: %v", rec)
		return
	}
	switch rec.RecType {
	case "harvester":
		harvester.ProcessRecord(rec, proc.writeChannel)
	case "full_node":
		fullnode.ProcessRecord(rec, proc.writeChannel)
	case "wallet":
		// TODO
	}
}

func (proc *Processor) CheckIntervalSec() int {
	return proc.checkIntervalSec
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
		checkIntervalSec: conf.CheckIntervalSec,
		filePath:         conf.ChiaLogPath,
		tzOffset:         conf.TimeZoneOffset,
		writeChannel:     wch,
		writer:           writer,
	}, nil
}
