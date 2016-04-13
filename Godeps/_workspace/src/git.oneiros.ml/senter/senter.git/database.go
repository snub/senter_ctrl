package senter

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var dbHandle *gorm.DB = nil

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Unix(time.Now().Unix(), 0).UTC()
	}
}

func InitDatabase(config *DatabaseConfig) error {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&autocommit=false", config.Username, config.Password, config.Host, config.Port, config.Database)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("opening database handle failed: %s\n", err)
	}

	err = db.DB().Ping()
	if err != nil {
		logger.Fatalf("cannot connect to database: %s\n", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.SingularTable(true)
	dbHandle = &db
	return nil
}

func CloseDatabase() {
	db := getDb()
	if db != nil {
		err := db.Close()
		if err != nil {
			logger.Printf("closing database handle failed: %s\n", err)
		}
	}
}

func EnableDatabaseLogger() {
	if dbHandle == nil {
		logger.Fatalln("no active database handle")
	}
	dbHandle.LogMode(true)
}

func DisableDatabaseLogger() {
	if dbHandle == nil {
		logger.Fatalln("no active database handle")
	}
	dbHandle.LogMode(false)
}

func getDb() *gorm.DB {
	if dbHandle == nil {
		logger.Fatalln("no active database handle")
	}
	err := dbHandle.DB().Ping()
	if err != nil {
		logger.Fatalf("cannot connect to database: %s\n", err)
	}
	return dbHandle
}
