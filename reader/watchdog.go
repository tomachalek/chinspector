// Copyright 2019 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2019 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
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

package reader

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"chinspector/config"
	"chinspector/logrecord"
)

const (
	defaultTickerIntervalSecs = 30
)

func Run(conf *config.Props, recProcessor *logrecord.Processor, finishEvent chan<- bool) {
	tickerInterval := time.Duration(conf.CheckIntervalSec)
	if tickerInterval == 0 {
		log.Printf("WARNING: intervalSecs for tail mode not set, using default %ds", defaultTickerIntervalSecs)
		tickerInterval = time.Duration(defaultTickerIntervalSecs)

	} else {
		log.Printf("INFO: configured to check for file changes every %d second(s)", tickerInterval)
	}
	ticker := time.NewTicker(tickerInterval * time.Second)
	log.Printf("Check interval: %s", tickerInterval*time.Second)
	quitChan := make(chan bool, 10)
	syscallChan := make(chan os.Signal, 10)
	signal.Notify(syscallChan, os.Interrupt)
	signal.Notify(syscallChan, syscall.SIGTERM)

	worklog := NewWorklog()
	err := worklog.Init()
	var rdr *FileTailReader

	if err != nil {
		log.Print("ERROR: ", err)
		quitChan <- true

	} else {
		wlItem := worklog.GetData()
		rdr, err = NewReader(recProcessor, wlItem.Inode, wlItem.Seek)
		if err != nil {
			log.Print("ERROR: ", err)
			quitChan <- true
		}
	}

	for {
		select {
		case ts := <-ticker.C:
			var wg sync.WaitGroup
			wg.Add(1)
			go func(rdr *FileTailReader) {
				rdr.Processor().OnCheckStart(ts)
				rdr.ApplyNewContent(
					func(v string) {
						rdr.Processor().OnLineRead(v)
					},
					func(inode int64, seek int64) {
						worklog.UpdateFileInfo(inode, seek)
					},
				)
				rdr.Processor().OnCheckStop()
				wg.Done()
			}(rdr)
			wg.Wait()
		case quit := <-quitChan:
			if quit {
				ticker.Stop()
				recProcessor.OnQuit()
				worklog.Close()
				finishEvent <- true
			}
		case <-syscallChan:
			log.Print("INFO: Caught signal, exiting...")
			ticker.Stop()
			rdr.Processor().OnQuit()
			worklog.Close()
			finishEvent <- true
		}
	}
}
