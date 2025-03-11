package common

import (
	"os/user"

	"github.com/rs/zerolog/log"
)

var HomeFileCfg string

func init() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	HomeFileCfg = currentUser.HomeDir + "/.config/synctainer/config.toml"

}
