package mqdb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var User string
var Pass string
var Host string
var Port string
var Name string

var Debug bool
var PingEachMinute int
var MaxIdleConns int
var MaxOpenConns int

var DB *gorm.DB
var lastPing time.Time
var err error

func New() (*gorm.DB, error) {
	if DB == nil {
		DB, err = Connect()
		if err != nil {
			return DB, err
		}
		return New()
	}
	if PingEachMinute > 0 && time.Now().After(lastPing.Add(time.Duration(PingEachMinute)*time.Minute)) {
		lastPing = time.Now()
		err := DB.DB().Ping()
		if err != nil {
			err = DB.Close()
			if err != nil {
				return DB, err
			}
			DB, err = Connect()
			if err != nil {
				return DB, err
			}
			return New()
		}
	}

	return DB, nil
}

func Connect() (*gorm.DB, error) {
	dbLink := GetLInk()
	var err error
	DB, err = gorm.Open("mysql", dbLink)
	DB.LogMode(Debug)
	if err != nil {
		return DB, err
	}
	DB.DB().SetMaxIdleConns(MaxIdleConns)
	DB.DB().SetMaxOpenConns(MaxOpenConns)

	return DB, nil
}

func Close() error {
	err = DB.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetLInk() string {
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", User, Pass, Host, Port, Name)

	return dbLink
}
