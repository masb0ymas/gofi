package config

import "time"

type Config struct {
	App    ConfigApp
	DB     ConfigDB
	Redis  ConfigRedis
	Resend ConfigResend
	Google ConfigGoogle
}

type ConfigApp struct {
	Env       string
	Debug     bool
	Port      int
	MachineID uint16
	Name      string
	Secret    string
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

type ConfigRedis struct {
	Addr     string
	Password string
	DB       int
}

type ConfigResend struct {
	ApiKey       string
	FromEmail    string
	DebugToEmail string
}

type ConfigGoogle struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
