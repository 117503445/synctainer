package main

import (
	"github.com/117503445/synctainer/cmd/cli/common"
	"github.com/117503445/synctainer/cmd/cli/tui"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

type cmdSyncImage struct {
	Image    string `arg:"" help:"Image name"`
	Platform string `help:"Platform" default:"linux/amd64"`
}

type cmdSyncCompose struct {
}

type cmdEditConfig struct {
}

var cli struct {
	SyncImage   cmdSyncImage   `cmd:"sync-image" help:"Sync image"`
	SyncCompose cmdSyncCompose `cmd:"sync-compose" help:"Sync compose to OSS"`
	EditConfig  cmdEditConfig  `cmd:"edit-config" help:"Edit config"`

	BackendHost  string `help:"Backend host" default:"https://synctainer-api.117503445.top"`
	RegistryHost string `help:"Registry host" default:"registry.cn-hangzhou.aliyuncs.com"`
	RegistryUser string
	RegistryPass string
}

func CliLoad() {

	// const DEV_CFG = false
	// fileCfg := ""
	// if DEV_CFG {
	// 	fileCfg = "/workspace/config.toml"
	// }
	// if !goutils.FileExists(fileCfg) {
	// 	fileCfg = common.HomeFileCfg
	// 	if !goutils.FileExists(fileCfg) {
	// 		tui.EditCfg()
	// 	}
	// }

	var ctx *kong.Context

	parse := func() {
		ctx = kong.Parse(&cli, kong.Configuration(kongtoml.Loader, common.HomeFileCfg))
		log.Info().Interface("cli", cli).Send()
	}
	parse()

	checkMissing := func() bool {
		m := map[string]string{
			"registryHost": cli.RegistryHost,
			"registryUser": cli.RegistryUser,
			"registryPass": cli.RegistryPass,
		}

		fieldMissing := false
		for k, v := range m {
			if v == "" {
				log.Warn().Str("key", k).Msg("field missing")
				fieldMissing = true
			}
		}
		return fieldMissing
	}

	if checkMissing() {
		// if fileCfg != common.HomeFileCfg {
		// 	log.Fatal().Msg("field missing")
		// }
		tui.EditCfg()
		parse()
		if checkMissing() {
			log.Fatal().Msg("field missing")
		}
	}

	err := ctx.Run()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
