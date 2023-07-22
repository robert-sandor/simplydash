package internal

import "time"

const Version = "dev"

const (
	DefaultDockerInterval = 10 * time.Second
	DefaultDockerTimeout  = 5 * time.Second

	DefaultEnableHealthcheck   = false
	DefaultHealthcheckInterval = 10 * time.Second
	DefaultHealthcheckTimeout  = 5 * time.Second
)
