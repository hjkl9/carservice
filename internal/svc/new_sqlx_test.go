package svc

import (
	"carservice/internal/config"
	"fmt"
	"testing"
)

func TestNewSqlx(t *testing.T) {
	dbc := NewSqlx(config.MysqlConf{
		Driver: "mysql",
		Source: "root:@tcp(127.0.0.1:3306)/carservicedb",
	})
	rows, err := dbc.Query("SELECT 1")
	if err != nil {
		t.Errorf("数据库查询数据时发生了错误, err: %s\n", err.Error())
		return
	}
	cols, _ := rows.Columns()
	fmt.Printf("%#v\n", cols)
}
