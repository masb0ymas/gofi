package config

import "time"

type Config struct {
	App    ConfigApp
	DB     ConfigDB
	Resend ConfigResend
}

type ConfigApp struct {
	Env       string
	Debug     bool
	Port      int
	MachineID uint16
	Name      string
	JWTSecret string
	ClientURL string
	ServerURL string
}

type ConfigDB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type ConfigResend struct {
	ApiKey       string
	FromEmail    string
	DebugToEmail string
}
