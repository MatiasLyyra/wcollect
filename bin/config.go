package main

type Config struct {
	Clickhouse Clickhouse `json:"clickhouse"`
	Locations  []string   `json:"locations"`
}

type Clickhouse struct {
	Addr     []string `json:"address"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Database string   `json:"database"`
}
