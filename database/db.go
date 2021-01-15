package database

import (
	"context"
	"time"

	"github.com/gyf841010/pz-infra-new/log"
	"github.com/gyf841010/pz-infra-new/logging"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var globalDB *gorm.DB

type MyLogger struct {
	//logger logrus.Logger
	logger logging.Logger
}

func (l *MyLogger) LogMode(logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *MyLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(msg, logging.With("data", append([]interface{}{}, data...)))
}

func (l *MyLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(msg, logging.With("data", append([]interface{}{}, data...)))
	//l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *MyLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(msg, logging.With("data", append([]interface{}{}, data...)))
	//l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *MyLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if rows == -1 {
		l.logger.Info("trace info",
			logging.With("err", err),
			logging.With("time", float64(elapsed.Nanoseconds())/1e6),
			logging.With("sql", sql))
	} else {
		l.logger.Info("trace info",
			logging.With("err", err),
			logging.With("time", float64(elapsed.Nanoseconds())/1e6),
			logging.With("row", rows),
			logging.With("sql", sql))

	}
}

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
	//ErrInvalidSQL           = gorm.ErrInvalidSQL
	ErrInvalidTransaction = gorm.ErrInvalidTransaction
	//ErrCantStartTransaction = gorm.ErrCantStartTransaction
	//ErrUnaddressable        = gorm.ErrUnaddressable
)

func InitDB(connectString string, loggers logging.Logger) error {
	var gormlog logger.Interface
	if loggers != nil {
		gormlog = &MyLogger{
			logger: loggers,
		}
	} else {
		gormlog = logger.Default
	}
	gormDB, err := gorm.Open(mysql.Open(connectString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormlog,
	})
	if err != nil {
		log.Errorf("init db error with url %s failed: %s", connectString, err.Error())
		panic(err)
	}
	gormDB.Logger.LogMode(logger.Info)
	SetDB(gormDB)
	return err
}

func LogMode(level logger.LogLevel) {
	if globalDB != nil {
		globalDB.Logger.LogMode(level)
	}
}

func GetDB() *gorm.DB {
	return globalDB
}

func SetDB(db *gorm.DB) {
	globalDB = db
}
