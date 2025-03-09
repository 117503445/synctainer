package main

import (
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
	kong.Parse(&cfg, kong.Configuration(kongtoml.Loader, "/workspace/config.toml"))
	log.Info().Interface("cfg", cfg).Send()
}
