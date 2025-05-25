# wcollect

Tool to collect weather data from [FMI](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje) and save it to Clickhouse database, for display in Grafana dashboard.

## Building

```
make wcollect
```

## Install

```
sudo make install
```

This installs the binary to /usr/local/bin

## Usage

Running the program will read the configuration file from /etc/wcollect/config.json by default or the path can be changed with `-config` flag.

Then just add the binary to crontab 

```
0 */2 * * * /usr/local/bin/wcollect
```