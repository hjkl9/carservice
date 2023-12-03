package svc

import (
	"carservice/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	DBC    *sqlx.DB
	RDBC   *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DBC:    NewSqlx(c),
		RDBC:   NewRedis(c),
	}
}

func NewSqlx(c config.Config) *sqlx.DB {
	return sqlx.MustOpen(c.MysqlConf.Driver, c.MysqlConf.Source)
}

func NewRedis(c config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.RedisConf.Addr,
		DB:       c.RedisConf.DB,
		Password: c.RedisConf.Requirepass,
	})
}
