package main

import (
	"github.com/117503445/goutils"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/rs/zerolog/log"
)

func (c *cmdCreateTable) Run() error {
	client := tablestore.NewClient(cli.TablestoreEndpoint, cli.TablestoreName, cli.TablestoreAk, cli.TablestoreSk)

	tableMeta := new(tablestore.TableMeta)
	tableMeta.TableName = "tasks"
	tableMeta.AddPrimaryKeyColumn("id", tablestore.PrimaryKeyType_STRING)

	tableOption := new(tablestore.TableOption)
	//数据的过期时间，-1表示永不过期。
	tableOption.TimeToAlive = -1
	//最大版本数，属性列值最多保留1个版本，即保存最新的版本。
	tableOption.MaxVersion = 1
	//有效版本偏差，即写入数据的时间戳与系统当前时间的偏差允许最大值为86400秒（1天）。
	tableOption.DeviationCellVersionInSec = 86400

	//设置预留读写吞吐量，容量型实例下的数据表只能设置为0，高性能型实例下的数据表可以设置为非零值。
	reservedThroughput := new(tablestore.ReservedThroughput)
	reservedThroughput.Readcap = 0
	reservedThroughput.Writecap = 0

	request := new(tablestore.CreateTableRequest)
	request.TableMeta = tableMeta
	request.TableOption = tableOption
	request.ReservedThroughput = reservedThroughput

	response, err := client.CreateTable(request)
	if err != nil {
		log.Fatal().Err(err).Msg("creating table")
	}

	log.Info().Interface("response", response).Msg("created table")

	return nil
}

func main() {
	goutils.InitZeroLog()
	cfgLoad()
}
