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
	FcCallback  string  `env:"FC_CALLBACK"`

	Image    string `env:"IMAGE"`
	Platform string `env:"PLATFORM" default:"linux/amd64"`
	TaskId   string `env:"TASK_ID"`
	Registry string `env:"REGISTRY"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
}

func cfgCheck(role string) {
	if role == "fc" {
		if Cfg.GithubToken == "" {
			log.Warn().Msg("GithubToken is empty")
		}
	} else if role == "gh" {
		if Cfg.Registry == "" {
			log.Fatal().Msg("Registry is empty")
		}
		if Cfg.Image == "" {
			log.Fatal().Msg("Image is empty")
		}
		if Cfg.Username == "" {
			log.Fatal().Msg("Username is empty")
		}
		if Cfg.Password == "" {
			log.Fatal().Msg("Password is empty")
		}
	} else {
		log.Fatal().Str("role", role).Msg("unknown role")
	}
}

func CfgLoad(role string) {
	// log.Info().Str("ver", os.Getenv("VER")).Send()
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

	cfgCheck(role)

	log.Info().
		Str("FcCallback", Cfg.FcCallback).
		Str("Platform", Cfg.Platform).
		Str("Image", Cfg.Image).
		Str("TaskId", Cfg.TaskId).
		Str("Registry", Cfg.Registry).
		Str("Username", Cfg.Username).
		Send()
}
