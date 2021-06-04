package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

type InfluxProps struct {
	Server          string `json:"server"`
	Token           string `json:"token"`
	Organization    string `json:"organization"`
	Bucket          string `json:"bucket"`
	RetentionPolicy string `json:"retentionPolicy"`
}

// Validate tests whether the configuration is filled in
// correctly. Please note that if the function returns nil
// then IsConfigured() must return 'true'.
func (conf *InfluxProps) Validate() error {
	var err error
	if conf.Server == "" {
		return fmt.Errorf("Missing 'server' information for InfluxDB")
	}
	_, err = url.Parse(conf.Server)
	if err != nil {
		return fmt.Errorf("Invalid InfluxDB server URL: %s", conf.Server)
	}
	if conf.Token == "" {
		return fmt.Errorf("Missing 'token' information for InfluxDB")
	}
	if conf.Organization == "" {
		return fmt.Errorf("Missing 'organization' information for InfluxDB")
	}
	if conf.Bucket == "" {
		return fmt.Errorf("Missing 'bucket' information for InfluxDB")
	}
	return nil
}

type Props struct {
	ChiaLogPath string `json:"chiaLogPath"`

	// TimeZoneOffset specifies a local timezone by
	// providing hh:mm offset. E.g. +02:00, -03:00
	TimeZoneOffset        string      `json:"timeZoneOffset"`
	CheckIntervalSec      int         `json:"checkIntervalSec"`
	ErrCountTimeRangeSecs int         `json:"errCountTimeRangeSecs"`
	NumErrorsAlarm        int         `json:"numErrorsAlarm"`
	InfluxDB              InfluxProps `json:"influxDb"`
}

// Load loads main configuration
func Load(path string) *Props {
	rawData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	var conf Props
	json.Unmarshal(rawData, &conf)
	return &conf
}
