package binlog

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"log"
)

type BinlogSync struct {
	canal.DummyEventHandler
}

var BinLogDataChan = make(chan interface{}, 1024)
var BinLogChan = make(chan string, 1024)

// OnRow 获取 event_type 为 write_rows, update_rows, delete_rows 的数据
func (h *BinlogSync) OnRow(ev *canal.RowsEvent) error {
	//rowData := make(map[string]interface{})
	//rowList := make([]interface{}, len(ev.Rows))

	//fmt.Println("原始数据：", ev.Rows)
	fmt.Printf("sql的操作行为：%s\t", ev.Action)

	switch ev.Action {
	case canal.UpdateAction:
		BinLogDataChan <- ev.Rows[1]
		BinLogChan <- "Update"
		fmt.Println("修改數據", <-BinLogDataChan)
	case canal.InsertAction:
		//BinLogDataChan <- ev.Rows
		//BinLogChan <- "Insert"
		fmt.Println("新增數據", ev.Rows)
	case canal.DeleteAction:
		BinLogChan <- "Delete"
		fmt.Println("刪除數據")
	default:
		fmt.Println("unknown action")
	}

	for idx, _ := range ev.Table.PKColumns {
		fmt.Printf("主键为：%s\n", ev.Table.Columns[ev.Table.PKColumns[idx]].Name)
	}
	//for idxRow, _ := range ev.Rows {
	//	for columnIndex, currColumn := range ev.Table.Columns {
	//		// 字段名和对应的值
	//		row := fmt.Sprintf("%v:%v", currColumn.Name, ev.Rows[idxRow][columnIndex])
	//		fmt.Println(row)
	//		rowData[currColumn.Name] = ev.Rows[idxRow][columnIndex]
	//		rowList[idxRow] = rowData
	//	}
	//}
	//
	//rowJson, err := json.Marshal(rowList)
	//if err != nil {
	//	return fmt.Errorf("序列化错误：%s", err)
	//}
	//
	//fmt.Printf("序列化为json格式：%s\n\n", string(rowJson))
	return nil
}

func (h *BinlogSync) OnRotate(header *replication.EventHeader, r *replication.RotateEvent) error {
	fmt.Printf("下一个日志为 %s 位置为 %d EventType %s \n", string(r.NextLogName), r.Position, header.EventType)
	return nil
}

func (h *BinlogSync) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	result := fmt.Sprintf("修改了数据库%s中表%s的结构", schema, table)
	fmt.Println(result)
	return nil
}

// OnDDL query 事件中的一些信息，如执行的 sql 语句
func (h *BinlogSync) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	fmt.Println(string(queryEvent.Query))
	return nil
}

func (h *BinlogSync) OnXID(header *replication.EventHeader, m mysql.Position) error {
	fmt.Println("XID", m.Pos)
	return nil
}

func StartbinLog() {
	cfg := canal.NewDefaultConfig()
	cfg.Password = "123456"
	cfg.Dump.TableDB = "book"
	cfg.ServerID = 12
	cfg.Dump.Tables = []string{"books"}

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Register a handler to handler Events
	c.SetEventHandler(&BinlogSync{})

	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
