package model

import (
	"fmt"
	"github.com/Zhenghao-Liu/OAuth_demo/config"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	OAuthDemoDB    *gorm.DB
	OAuthDemoCache *redis.Client
)

func InitDatabase() {
	OAuthDemoDB = openMysql(config.ConfigInstance.OAuthDemoDB)
}

func InitRedis() {
	OAuthDemoCache = openRedis(config.ConfigInstance.OAuthDemoCache)
}

func openMysql(mysqlConfig config.MysqlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", mysqlConfig.UserName, mysqlConfig.Password,
		mysqlConfig.DefaultHost, mysqlConfig.Database, mysqlConfig.Settings)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func openRedis(redisConfig config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DBIdx,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	}
	return rdb
}
