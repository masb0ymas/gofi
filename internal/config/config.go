package config

import "time"

type Config struct {
	App ConfigApp
	DB  ConfigDB
}

type ConfigApp struct {
	Env       string
	Debug     bool
	Port      int
	MachineID uint16
}

type ConfigDB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}
