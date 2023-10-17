package repository

import (
	"colatiger/pkg/log"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"time"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	logger *log.Logger
}

func NewRepository(logger *log.Logger, db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
		rdb:    rdb,
	}
}

func NewDB(conf *viper.Viper, log *log.Logger) *gorm.DB {
	logger := zapgorm2.New(log.Logger)
	logger.SetAsDefault()
	db, err := gorm.Open(mysql.Open(conf.GetString("data.mysql.user")), &gorm.Config{Logger: logger})
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	return db
}

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}
