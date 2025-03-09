package main

import (
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

var cfg struct {
	FcCallback string `env:"FC_CALLBACK"`
	RunId      string `env:"RUN_ID"`

	Image    string `env:"IMAGE"`
	Platform string `env:"PLATFORM" default:"linux/amd64"`
	TaskId   string `env:"TASK_ID"`
	Registry string `env:"REGISTRY"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

func CfgLoad() {
	kong.Parse(&cfg, kong.Configuration(kongtoml.Loader))
	if cfg.Registry == "" {
		log.Fatal().Msg("Registry is empty")
	}
	if cfg.Image == "" {
		log.Fatal().Msg("Image is empty")
	}
	if cfg.Username == "" {
		log.Fatal().Msg("Username is empty")
	}
	if cfg.Password == "" {
		log.Fatal().Msg("Password is empty")
	}

	// must not print password
	log.Info().
		Str("FcCallback", cfg.FcCallback).
		Str("Platform", cfg.Platform).
		Str("Image", cfg.Image).
		Str("TaskId", cfg.TaskId).
		Str("Registry", cfg.Registry).
		Str("Username", cfg.Username).
		Str("RunId", cfg.RunId).
		Send()
}
