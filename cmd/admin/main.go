package main

import (
	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
)

func (c *cmdCreateTable) Run() error {
	log.Info().Msg("Admin")

	return nil
}

func main() {
	goutils.InitZeroLog()
	cfgLoad()
	// log.Debug().Strs("args", os.Args).Send()
	// remove os.args[1], remain os.args[0], os.args[2], os.args[3] ...
	// os.Args = slices.Delete(os.Args, 1, 2)
	// cfg.CfgLoad("admin")

}
