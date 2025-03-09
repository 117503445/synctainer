package ots

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/rs/zerolog/log"
)

const tableName = "tasks"

// TableManager 管理表格存储操作
type TableManager struct {
	client tablestore.TableStoreClient
}

// NewTableManager 创建一个新的 TableManager
func NewTableManager(endpoint, instanceName, accessKeyID, accessKeySecret string) (*TableManager, error) {
	client := tablestore.NewClient(endpoint, instanceName, accessKeyID, accessKeySecret)
	return &TableManager{
		client: *client,
	}, nil
}

// CreateTable 创建一个新表
func (tm *TableManager) CreateTable() {
	tableMeta := &tablestore.TableMeta{
		TableName: tableName,
	}
	tableMeta.AddPrimaryKeyColumn("id", tablestore.PrimaryKeyType_STRING)

	tableOption := &tablestore.TableOption{
		TimeToAlive: -1, // 数据永不过期
		MaxVersion:  1,  // 属性列值最多保留1个版本
	}
	tableOption.DeviationCellVersionInSec = 86400 // 写入数据的时间戳偏差允许最大值为1天

	reservedThroughput := &tablestore.ReservedThroughput{
		Readcap:  0, // 容量型实例下的数据表只能设置为0
		Writecap: 0,
	}

	request := &tablestore.CreateTableRequest{
		TableMeta:          tableMeta,
		TableOption:        tableOption,
		ReservedThroughput: reservedThroughput,
	}

	if resp, err := tm.client.CreateTable(request); err != nil {
		log.Fatal().Err(err).Msg("create table")
	} else {
		log.Info().Interface("resp", resp).Msg("create table done")
	}
}

func (tm *TableManager) PutRow(id string, column map[string]any) error {
	putRowRequest := new(tablestore.PutRowRequest)
	putRowChange := new(tablestore.PutRowChange)
	putRowChange.TableName = tableName
	putPk := new(tablestore.PrimaryKey)
	putPk.AddPrimaryKeyColumn("id", id)
	putRowChange.PrimaryKey = putPk

	for k, v := range column {
		putRowChange.AddColumn(k, v)
	}
	putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	putRowRequest.PutRowChange = putRowChange
	resp, err := tm.client.PutRow(putRowRequest)
	if err != nil {
		return err
	}
	log.Info().Discard().Interface("response", resp).Msg("put row done")

	return nil
}
