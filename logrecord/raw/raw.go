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

package raw

import (
	"fmt"
	"strings"
	"time"
)

type PreparsedLine struct {
	Dt   time.Time
	Rest string
}

func minVal(v1 int, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func (pl *PreparsedLine) String() string {
	return fmt.Sprintf("PreparsedLine{dt: %v, rest: %s}", pl.Dt, pl.Rest[:minVal(20, len(pl.Rest))]+"...")
}

type PreparsedRecord struct {
	Dt   time.Time
	Text strings.Builder
}

func ExtractDate(line string, tz string) *PreparsedLine {
	items := strings.SplitN(line, " ", 2)
	t, err := time.Parse(time.RFC3339Nano, items[0]+tz)
	if err != nil {
		return &PreparsedLine{time.Time{}, items[0]}
	}
	return &PreparsedLine{t, items[1]}
}

// Record is a preparsed log record with known date,
// service, log level. The message itself remains unparsed.
type Record struct {
	Datetime time.Time
	Service  string
	Module   string
	Level    string
	Message  string
}

func (rec Record) String() string {
	return fmt.Sprintf(
		"Record{Datetime: %s, Service: %s, Module: %s, Level: %s, Message: %s",
		rec.Datetime, rec.Service, rec.Module, rec.Level, rec.Message,
	)
}

func ParseRecord(rec *PreparsedRecord) *Record {
	items := strings.Split(rec.Text.String(), " ")
	if len(items) >= 5 {
		return &Record{
			Datetime: rec.Dt,
			Service:  items[0],
			Module:   items[1],
			Level:    items[2],
			Message:  strings.Join(items[3:], " "),
		}
	}
	return nil
}
