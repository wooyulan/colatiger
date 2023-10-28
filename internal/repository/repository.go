package repository

import (
	"colatiger/config"
	"colatiger/internal/model"
	"colatiger/pkg/helper/path"
	"context"
	"fmt"
	"github.com/jassue/gin-wire/app/service"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
	sf  *sonyflake.Sonyflake
}

func NewRepository(logger *zap.Logger, db *gorm.DB, rdb *redis.Client, sf *sonyflake.Sonyflake) (*Repository, func(), error) {
	cleanup := func() {
		logger.Info("closing the data resources")
	}
	return &Repository{
		db:  db,
		rdb: rdb,
		sf:  sf,
	}, cleanup, nil
}

func NewDB(conf *config.Configuration, gLog *zap.Logger) *gorm.DB {
	if conf.Database.Driver != "mysql" {
		panic(conf.Database.Driver + " driver is not supported")
	}

	var writer io.Writer
	var logMode logger.LogLevel

	// 是否启用日志文件
	if conf.Database.EnableFileLogWriter {
		logFileDir := conf.Log.RootDir
		if !filepath.IsAbs(logFileDir) {
			logFileDir = filepath.Join(path.RootPath(), logFileDir)
		}
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   filepath.Join(logFileDir, conf.Database.LogFilename),
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxAge,
			Compress:   conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}

	switch conf.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	newLogger := logger.New(
		log.New(writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                        // 慢查询 SQL 阈值
			Colorful:                  !conf.Database.EnableFileLogWriter, // 禁用彩色打印
			IgnoreRecordNotFoundError: false,                              // 忽略ErrRecordNotFound（记录未找到）错误
			LogLevel:                  logMode,                            // Log lever
		},
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Database.UserName,
		conf.Database.Password,
		conf.Database.Host,
		strconv.Itoa(conf.Database.Port),
		conf.Database.Database,
		conf.Database.Charset,
	)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Database.TablePrefix,
			//SingularTable: true,

		},
		DisableForeignKeyConstraintWhenMigrating: true,      // 禁用自动创建外键约束
		Logger:                                   newLogger, // 使用自定义 Logger
	}); err != nil {
		gLog.Error("failed opening connection to err:", zap.Any("err", err))
		panic("failed to connect database")
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
		return db
	}
}

func NewRedis(c *config.Configuration, gLog *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + c.Redis.Port,
		Password: c.Redis.Password, // no password set
		DB:       c.Redis.DB,       // use default DB
	})
	//client.AddHook(redisotel.TracingHook{})
	if err := client.Ping(context.Background()).Err(); err != nil {
		gLog.Fatal("redis connect failed, err:", zap.Any("err", err))
		panic("failed to connect redis")
	}

	return client
}

type contextTxKey struct{}

func (d *Repository) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Repository) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

// NewTransaction .
func NewTransaction(d *Repository) service.Transaction {
	return d
}

// 数据库表初始化
func initMySqlTables(db *gorm.DB) {
	log.Print("init table start...")
	err := db.AutoMigrate(
		model.User{},
	)
	if err != nil {
		log.Printf("migrate table failed err:%v", zap.Any("err", err))
		os.Exit(0)
	}
	log.Print("init table end...")
}
