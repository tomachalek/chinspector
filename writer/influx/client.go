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

package influx

import (
	"log"
	"time"

	"chinspector/config"

	influxv2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxRecord represents a record suitable to be inserted
// into an InfluxDB database
type InfluxRecord interface {
	Tags() map[string]string
	Fields() map[string]interface{}
	Time() time.Time

	// Measurements specifies a logical group of results; in case of Chinspector,
	// it is typically: harvester, blockchain, wallet
	Measurement() string

	// String is just for convenience - i.e. there is no rule how to resulting string
	// should look like
	String() string
}

// RecordWriter allows writing InfluxRecord instances into
// an InfluxDb v2 instance via its recent API
type RecordWriter struct {
	conn         influxv2.Client
	writeAPI     api.WriteAPI
	address      string
	organization string
	bucket       string
	incomingData <-chan InfluxRecord
}

func (c *RecordWriter) addRecord(rec InfluxRecord) {
	point := influxv2.NewPointWithMeasurement(rec.Measurement())
	for k, v := range rec.Tags() {
		point.AddTag(k, v)
	}
	for k, v := range rec.Fields() {
		point.AddField(k, v)
	}
	point.SetTime(rec.Time())
	c.writeAPI.WritePoint(point)
}

// Finish ensures that the current operation is fully
// processed and all the data are written to InfluxDB.
func (c *RecordWriter) Finish() {
	c.writeAPI.Flush()
	c.conn.Close()
}

// NewRecordWriter is a factory function for RecordWriter
func NewRecordWriter(conf *config.InfluxProps, incoming <-chan InfluxRecord) (*RecordWriter, error) {
	conn := influxv2.NewClientWithOptions(
		conf.Server,
		conf.Token,
		influxv2.DefaultOptions().SetBatchSize(20),
	)
	writeAPI := conn.WriteAPI(conf.Organization, conf.Bucket)
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Printf("ERROR: write error: %s\n", err.Error())
		}
	}()
	w := &RecordWriter{
		conn:         conn,
		writeAPI:     writeAPI,
		address:      conf.Server,
		organization: conf.Organization,
		bucket:       conf.Bucket,
		incomingData: incoming,
	}
	go func() {
		for item := range incoming {
			//log.Print("DEBUG: sending >>>> ", item)
			w.addRecord(item)
		}
		w.Finish()
	}()
	return w, nil
}
