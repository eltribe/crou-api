package database

import (
	"crou-api/config"
	"crou-api/internal/domains"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"io"

	"log"
	"os"
	"time"
)

const (
	DefaultMaxOpenConns = 25
	DefaultMaxIdleConns = 25
)

type Persistent interface {
	DB() *gorm.DB
	REDIS() *redis.Client
}

// Database struct
type Database struct {
	Db    *gorm.DB
	Redis *redis.Client
}

// DB write DB
// return *gorm.DB
func (db *Database) DB() *gorm.DB {
	return db.Db
}

func (db *Database) REDIS() *redis.Client {
	return db.Redis
}

func NewDatabase(conf *config.Config) Persistent {

	gormLoggerConfig := gLogger.Config{
		SlowThreshold:             3 * time.Second, // Slow SQL threshold
		Colorful:                  true,            // Disable color
		IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound errors for logger
		LogLevel:                  gLogger.Silent,  // Log level
	}

	var logwriter io.Writer = os.Stdout
	if conf.Log.Type == "file" {
		logwriter = &lumberjack.Logger{
			Filename:   conf.Log.FileName,
			MaxSize:    20, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
		gormLoggerConfig.Colorful = false
	}

	newLogger := gLogger.New(
		log.New(logwriter, "\r\n", log.LstdFlags), // io writer
		gormLoggerConfig,
	)

	db, err := gorm.Open(dialector(conf.Database), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug()

	pool, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if conf.Database.MaxOpen <= 0 {
		conf.Database.MaxOpen = DefaultMaxOpenConns
	}
	if conf.Database.MaxIdle <= 0 {
		conf.Database.MaxIdle = DefaultMaxIdleConns
	}
	pool.SetMaxOpenConns(conf.Database.MaxOpen)
	pool.SetMaxIdleConns(conf.Database.MaxIdle)
	pool.SetConnMaxLifetime(5 * time.Minute)

	/* Redis */
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		DB:       conf.Redis.Db, // use default DB
		Password: "",            // no password set
	})

	return &Database{
		Db:    db,
		Redis: redisClient,
	}
}

var Alldomains = []interface{}{
	&domains.User{},
	&domains.UserDetail{},
}

func AutoMigration(db Persistent) {
	if err := db.DB().AutoMigrate(
		Alldomains...,
	); err != nil {
		log.Fatal(err)
	}
}

func dialector(conf config.DatabaseConfig) gorm.Dialector {
	log.Println(conf.Type)
	if conf.Type == "sqlite" {
		return sqlite.Open("sqlite.db")
	} else if conf.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Database,
		)
		return mysql.New(mysql.Config{
			DriverName:                "mysql",
			DSN:                       dsn,
			SkipInitializeWithVersion: false,
			DefaultStringSize:         255,  // change it if needed
			DisableDatetimePrecision:  true, // true, because datetime precision requires MySQL 5.6
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
		})
	} else if conf.Type == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require TimeZone=Asia/Seoul",
			conf.Host,
			conf.Port,
			conf.User,
			conf.Database,
			conf.Password,
		)
		return postgres.Open(dsn)
	} else {
		log.Fatal("Database type not found")
		return nil
	}
}
