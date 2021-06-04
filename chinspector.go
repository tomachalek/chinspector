package main

import (
	"flag"
	"log"

	"github.com/tomachalek/chinspector/config"
	"github.com/tomachalek/chinspector/logrecord"
	"github.com/tomachalek/chinspector/reader"
)

var (
	defaultTickerIntervalSecs = 10
	actionTail                = "tail"
	actionBatch               = "batch"
)

type ProcessOptions struct {
	dryRun bool
}

func main() {
	procOpts := new(ProcessOptions)
	flag.BoolVar(&procOpts.dryRun, "dry-run", false, "Do not write data to InfluxDB, just output parsed log info")
	flag.Parse()
	conf := config.Load(flag.Arg(0))
	action := flag.Arg(1)

	finishEvent := make(chan bool)
	go func() {
		switch action {
		case actionTail:
			proc, err := logrecord.NewProcessor(conf)
			if err != nil {
				log.Fatal(err)
			}
			reader.Run(conf, proc, finishEvent)
		case actionBatch:
			proc, err := logrecord.NewProcessor(conf)
			if err != nil {
				log.Fatal(err)
			}
			reader.RunBatch(conf, proc)
			finishEvent <- true
		}
	}()
	<-finishEvent
}
