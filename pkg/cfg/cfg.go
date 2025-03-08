package cfg

import (
	"os"
	"os/user"

	"github.com/117503445/goutils"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

var Cfg struct {
	GithubToken string
}

func cfgCheck() {
	if Cfg.GithubToken == "" {
		log.Warn().Msg("GithubToken is empty")
	}
}

func CfgLoad() {
	log.Info().Str("ver", os.Getenv("VER")).Send()

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

	kong.Parse(&Cfg, kong.Configuration(kongtoml.Loader, "/workspace/config.toml", fileHomeCfg))

	cfgCheck()

	// log.Info().
}
