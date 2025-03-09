package main

import (
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	"github.com/rs/zerolog/log"
)

type cmdCreateTable struct {
}

type cmdPutRow struct {
}

type cmdUpdateRow struct {
}

type cmdGetRow struct {
}

var cli struct {
	CreateTable cmdCreateTable `cmd:"" help:"create table"`
	PutRow      cmdPutRow      `cmd:"" help:"put row"`
	UpdateRow   cmdUpdateRow   `cmd:"" help:"update row"`
	GetRow      cmdGetRow      `cmd:"" help:"get row"`

	TablestoreAk       string
	TablestoreSk       string
	TablestoreEndpoint string
	TablestoreName     string
}

func CliLoad() {
	ctx := kong.Parse(&cli, kong.Configuration(kongtoml.Loader, "/workspace/config.toml"))
	log.Info().Interface("cli", cli).Send()
	err := ctx.Run()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
