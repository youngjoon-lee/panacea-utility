package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel             LogLevel `envconfig:"LOG_LEVEL" default:"info"`
	NodeRPCAddr          string   `envconfig:"NODE_RPC_ADDR" required:"true"`
	StartHeight          uint64   `envconfig:"START_HEIGHT" required:"true"`
	EndHeight            uint64   `envconfig:"END_HEIGHT" required:"true"`
	SafeBlockIntervalSec float64  `envconfig:"SAFE_BLOCK_INTERVAL_SEC" required:"true"`
	EpochHour            uint8    `envconfig:"EPOCH_HOUR" required:"true"`
	EpochMinute          uint8    `envconfig:"EPOCH_MINUTE" required:"true"`
	TargetAccAddr        string   `envconfig:"TARGET_ACCOUNT" required:"true"`
}

func mustLoadConfig() Config {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		log.Panic("failed to load env vars: %v", err)
	}

	return conf
}

type LogLevel log.Level

func (l *LogLevel) Decode(value string) error {
	level, err := log.ParseLevel(value)
	if err != nil {
		return err
	}

	*l = LogLevel(level)
	return nil
}
