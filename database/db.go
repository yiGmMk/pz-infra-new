package database

import (
	"github.com/gyf841010/pz-infra-new/log"
	"github.com/gyf841010/pz-infra-new/logging"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var globalDB *gorm.DB

var (
	ErrRecordNotFound       = gorm.ErrRecordNotFound
	ErrInvalidSQL           = gorm.ErrInvalidSQL
	ErrInvalidTransaction   = gorm.ErrInvalidTransaction
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	ErrUnaddressable        = gorm.ErrUnaddressable
)

func InitDB(connectString string, logger logging.Logger) error {
	gormDB, err := gorm.Open("mysql", connectString)
	if err != nil {
		log.Errorf("init db error with url %s failed: %s", connectString, err.Error())
		panic(err)
	}
	gormDB.LogMode(true)
	if logger != nil {
		gormDB.SetLogger(logger.GetPrintLogger())
	} else {
		gormDB.SetLogger(log.CurrentLogger())
	}
	SetDB(gormDB)
	return err
}

func LogMode(enable bool) {
	if globalDB != nil {
		globalDB.LogMode(enable)
	}
}

func GetDB() *gorm.DB {
	return globalDB
}

func SetDB(db *gorm.DB) {
	globalDB = db
}
