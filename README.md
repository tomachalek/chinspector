# Chinspector 

Chinspector is a simple log information extractor and agent for a Chia blockchain node/harvester.


🚧 Status note: this is a "raw and dirty" work in progress project 🚧


Chinspector is an agent you can install on the same server as your Chia services (node, harvester)
where it reads the Chia's `debug.log`, exctracts some useful information from it and stores it into
the InfluxDB database. The data can be reviewed either via Influx's own web dashboard of via
Grafana.

## Data storage and presentation

The easiest way how to run the whole data stack is to use Docker for both InfluxDB and Grafana. In the
Chinspector directory, run:

```
docker-compose up
```

Log into the InfluxDB container:

```
docker exec -it chinspector_influxdb_1 bash
```

Create an InfluxDB user:

```
influx auth create -o chinspector
```

Create an auth token for your apps (for simplicity, we stick with a single one for both reading
and writing data):

```
influx auth create -o chinspector --write-bucket status
```

Now open Grafana web page at http://localhost:3000, use default `admin:admin` credentials and
set something of your own for better security.

In Grafana, configure a new data source ("configuration" -> "data sources" -> "Add data source").
Select "InfluxDB".

Select "Flux" as the "query language".

In the *HTTP* section, fill in:

URL: http://chinspector_influxdb_1:8086

In *InfluxDB Details*:

Organization: chinspector
Token: the_token_you_have_obtained_when_creating_influxdb_credentials

Then "Save and test"

## Chinspector log agent

1. [download](https://golang.org/dl/) and [install](https://golang.org/doc/install) Go.
1. clone the Chinspector repository https://github.com/tomachalek/chinspector to the same
    server as your Chia installation (or to a server which is able to read the log file via a
    mounted filesystem)
1. `go build`
1. edit config.json

```json
{
    "chiaLogPath": "./data/debug.log",
    "timeZoneOffset": "+02:00",
    "checkIntervalSec": 5,
    "influxDb": {
        "server": "http://localhost:8096",
        "organization": "chinspector",
        "bucket": "status",
        "token": "Fu3vgKPMsMb_x0vqrlwV4YrNdtLa1nY_37PFRk68l2fnS0AE04S22fZNS6iK4cDcL7xaWQsNDd1B12_8C-rYA4w==",
        "retentionPolicy": "midterm"
    }
}
```
Please note the port 8096 - it is mapped from InfluxDB "native" 8086 by docker-compose.yml. In case
you want/need a different one, just edit docker-compose.yml and update the `config.json` accordingly.

Start the agent:

```
./chinspector ./config.json tail
```

To run the agent as a systemd service:

```
[Unit]
Description=A custom agent for collecting UCNK apps logs
After=network.target

[Service]
Type=simple
ExecStart=/path/to/chinspector /path/to/chinspector/config.json tail
User=your_user
Group=your_group

[Install]
WantedBy=multi-user.target
```
