package svc

import (
	"carservice/internal/config"
	"carservice/internal/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	Repo   data.RepoFactory
	DBC    *sqlx.DB
	RDBC   *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	dbc := NewSqlx(c)
	return &ServiceContext{
		Config: c,
		Repo:   data.NewDatastore(dbc), // todo: testing
		DBC:    dbc,                    // ! @deprecated
		RDBC:   NewRedis(c),
	}
}

func NewSqlx(c config.Config) *sqlx.DB {
	db := sqlx.MustOpen(c.MysqlConf.Driver, c.MysqlConf.Source)
	// set connection.
	return db
}

func NewRedis(c config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.RedisConf.Addr,
		DB:       c.RedisConf.DB,
		Password: c.RedisConf.Requirepass,
	})
}
