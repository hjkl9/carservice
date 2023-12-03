package svc

import (
	"carservice/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ServiceContext struct {
	Config config.Config
	DBC    *sqlx.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DBC:    NewSqlx(c.MysqlConf),
	}
}

func NewSqlx(c config.MysqlConf) *sqlx.DB {
	return sqlx.MustOpen(c.Driver, c.Source)
}
