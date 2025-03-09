package main

import (
	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/pkg/ots"
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
	if err != nil {
		return err
	}
	return nil
}

func main() {
	goutils.InitZeroLog()
	CliLoad()
}
