package core

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrInvalidChainID    = errors.New("invalid chain id")
	ErrInvalidERPCWSS    = errors.New("invalid erpc wss")
	ErrInvalidERPCHTTPS  = errors.New("invalid erpc https")
	ErrInvalidProtocolID = errors.New("invalid protocol id")
	ErrInvalidInstanceID = errors.New("invalid instance id")
)

type Config struct {
	ChainId    uint64 `env:"CHAIN_ID" env-default:"11155111"`
	ErpcWss    string `env:"ERPC_WSS"`
	ErpcHttps  string `env:"ERPC_HTTPS"`
	ProtocolID string `env:"PROTOCOL_ID"`
	InstanceID string `env:"INSTANCE_ID"`
}

func NewConfig() (Config, error) {
	conf := Config{}
	err := cleanenv.ReadEnv(&conf)
	conf.Validate()
	return conf, err
}

func (cfg *Config) Validate() error {
	if cfg.ChainId == 0 {
		return ErrInvalidChainID
	}
	if cfg.ErpcWss == "" {
		return ErrInvalidERPCWSS
	}
	if cfg.ErpcHttps == "" {
		return ErrInvalidERPCHTTPS
	}
	if cfg.ProtocolID == "" {
		return ErrInvalidProtocolID
	}
	if cfg.InstanceID == "" {
		return ErrInvalidInstanceID
	}
	return nil
}
