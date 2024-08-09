package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Algorithm string
type ActionType string

const (
	ServeAction   ActionType = "serve"
	ForwardAction ActionType = "forward"
)

const (
	WRR Algorithm = "WRR"
)

type ServerConfig struct {
	URI      string   `toml:"uri"`
	NAME     string   `toml:"name"`
	LISTEN   []string `toml:"listen"`
	MAXCONN  int16    `toml:"max_connections"`
	LOGFILE  string   `toml:"logfile"`
	LOGLEVEL string   `toml:"loglevel"`
	LOGNAME  string
}

type Backend struct {
	Address string `toml:"address"`
	Weight  int    `toml:"weight"`
}

type Pattern struct {
	URI    string `toml:"uri"`
	Action Action `toml:",inline"`
}

type Forward struct {
	Backends  []Backend `toml:"forward"`
	Algorithm Algorithm `toml:"algorithm"`
}

type Action struct {
	Type    ActionType `toml:"-"`
	Forward *Forward   `toml:"forward,omitempty"`
	Serve   *string    `toml:"serve,omitempty"`
}

type Config struct {
	Server  ServerConfig `toml:"server"`
	Pattern []Pattern    `toml:"match"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(filename string) (*Config, error) {

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	log.Printf("INFO: %v", c)

	return c, nil
}

func (c *Config) Get() *Config {
	return c
}
