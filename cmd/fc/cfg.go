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
	GithubToken string
	FcCallback  string `env:"FC_CALLBACK"`

	TablestoreAk       string
	TablestoreSk       string
	TablestoreEndpoint string
	TablestoreName     string
}

func cfgLoad() {
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
	log.Info().Interface("cfg", cfg).Send()
}
