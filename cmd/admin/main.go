package main

import (
	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/ots"
	"github.com/rs/zerolog/log"
)

func (c *cmdCreateTable) Run() error {
	tm, err := ots.NewTableManager(cli.TablestoreEndpoint, cli.TablestoreName, cli.TablestoreAk, cli.TablestoreSk)
	if err != nil {
		return err
	}
	tm.CreateTable()
	return nil
}

func (c *cmdPutRow) Run() error {
	tm, err := ots.NewTableManager(cli.TablestoreEndpoint, cli.TablestoreName, cli.TablestoreAk, cli.TablestoreSk)
	if err != nil {
		return err
	}
	err = tm.PutRow("task1", map[string]any{
		"k1": "v1",
	})
	return err
}

func (c *cmdUpdateRow) Run() error {
	tm, err := ots.NewTableManager(cli.TablestoreEndpoint, cli.TablestoreName, cli.TablestoreAk, cli.TablestoreSk)
	if err != nil {
		return err
	}
	err = tm.UpdateRow("task1", map[string]any{
		"k1": "v2",
		"k2": "v2",
	})
	return err
}

func (c *cmdGetRow) Run() error {
	tm, err := ots.NewTableManager(cli.TablestoreEndpoint, cli.TablestoreName, cli.TablestoreAk, cli.TablestoreSk)
	if err != nil {
		return err
	}
	row, err := tm.GetRow("task1")
	if err != nil {
		return err
	}
	log.Info().Interface("row", row).Msg("get row")
	log.Info().Str("k1", ots.MapMustGetString(row, "k1")).Send()

	return nil
}

func main() {
	goutils.InitZeroLog()
	CliLoad()
}
