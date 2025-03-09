package main

import (
	"os"
	"os/user"

	"github.com/117503445/goutils"
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
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fileHomeCfg := currentUser.HomeDir + "/.config/synctainer/config.toml"

	envCfg := os.Getenv("SYNCTAINER_CONFIG")
	if envCfg != "" {
		err := goutils.WriteText(fileHomeCfg, envCfg)
		if err != nil {
			log.Fatal().Err(err).Msg("WriteText")
		}
		log.Info().Str("envCfg", envCfg).Send()
	}

	kong.Parse(&cfg, kong.Configuration(kongtoml.Loader, "/workspace/config.toml", fileHomeCfg))
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
