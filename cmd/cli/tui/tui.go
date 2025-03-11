package tui

import (
	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/cmd/cli/common"
	"github.com/rs/zerolog/log"
)

func EditCfg() {
	log.Debug().Msg("EditCfg")
	oldCfg := make(map[string]any)
	if goutils.FileExists(common.HomeFileCfg) {
		err := goutils.ReadToml(common.HomeFileCfg, &oldCfg)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}
	log.Debug().Interface("oldCfg", oldCfg).Send()

}
