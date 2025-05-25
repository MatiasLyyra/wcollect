package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"githuv.com/MatiasLyyra/wcollect"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "/etc/wcollect/config.json", "config path for collector")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Printf("[ERROR] collect: %v", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := parseConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse config file %v: %v", configPath, err)
	}
	conn, err := createClient(config)
	if err != nil {
		return fmt.Errorf("failed to create clickhouse client to database %v: %v", config.Clickhouse.Database, err)
	}
	defer conn.Close()
	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database %v: %v", config.Clickhouse.Database, err)
	}
	now := time.Now()
	start := now.Add(-168 * time.Hour)
	forecastEnd := now.Add(168 * time.Hour)
	for _, loc := range config.Locations {
		observationData, err := wcollect.FetchObservations(loc, start, now)
		if err != nil {
			return fmt.Errorf("failed to collect observations: %v", err)
		}
		if err := wcollect.InsertData(context.Background(), conn, observationData); err != nil {
			return fmt.Errorf("failed to write data to database %v: %v", config.Clickhouse.Database, err)
		}

		forecastData, err := wcollect.FetchForecast(loc, now, forecastEnd)
		if err != nil {
			return fmt.Errorf("failed to collect forecasts: %v", err)
		}
		if err := wcollect.InsertData(context.Background(), conn, forecastData); err != nil {
			return fmt.Errorf("failed to write data to database %v: %v", config.Clickhouse.Database, err)
		}
	}
	return nil
}

func createClient(conf Config) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: conf.Clickhouse.Addr,
		Auth: clickhouse.Auth{
			Database: conf.Clickhouse.Database,
			Username: conf.Clickhouse.User,
			Password: conf.Clickhouse.Password,
		},
	})
	return conn, err
}

func parseConfig(path string) (Config, error) {
	var c Config
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}
