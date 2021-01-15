package database

import (
	"context"
	"errors"
	"time"

	"github.com/gyf841010/pz-infra-new/log"
	"github.com/gyf841010/pz-infra-new/logging"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var globalDB *gorm.DB

type MyLogger struct {
	//logger logrus.Logger
	logger        logging.Logger
	SlowThreshold time.Duration
	SourceField   string
}

func (l *MyLogger) LogMode(logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *MyLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	//l.logger.Info(msg, logging.With("data", append([]interface{}{}, data...)))
	l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *MyLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	//l.logger.Warn(msg, logging.With(msg, append([]interface{}{}, data...)))
	l.logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *MyLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	//l.logger.Error(msg, logging.With("data", append([]interface{}{}, data...)))
	l.logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *MyLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// elapsed := time.Since(begin)
	// sql, rows := fc()
	// if rows == -1 {
	// 	l.logger.Info("trace info",
	// 		logging.With("err", err),
	// 		logging.With("time", float64(elapsed.Nanoseconds())/1e6),
	// 		logging.With("sql", sql))
	// } else {
	// 	l.logger.Info("trace info",
	// 		logging.With("err", err),
	// 		logging.With("time", float64(elapsed.Nanoseconds())/1e6),
	// 		logging.With("row", rows),
	// 		logging.With("sql", sql))
	// }

	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		fields[logrus.ErrorKey] = err
		l.logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	l.logger.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}

var (
	ErrRecordNotFound        = gorm.ErrRecordNotFound
	ErrInvalidTransaction    = gorm.ErrInvalidTransaction
	ErrNotImplemented        = gorm.ErrNotImplemented
	ErrMissingWhereClause    = gorm.ErrMissingWhereClause
	ErrUnsupportedRelation   = gorm.ErrUnsupportedRelation
	ErrPrimaryKeyRequired    = gorm.ErrPrimaryKeyRequired
	ErrModelValueRequired    = gorm.ErrModelValueRequired
	ErrInvalidData           = gorm.ErrInvalidData
	ErrUnsupportedDriver     = gorm.ErrUnsupportedDriver
	ErrRegistered            = gorm.ErrRegistered
	ErrInvalidField          = gorm.ErrInvalidField
	ErrEmptySlice            = gorm.ErrEmptySlice
	ErrDryRunModeUnsupported = gorm.ErrDryRunModeUnsupported
)

func InitDB(connectString string, infraLogger logging.Logger) error {
	var gormlog logger.Interface
	if infraLogger != nil {
		gormlog = &MyLogger{
			SlowThreshold: 200 * time.Millisecond,
			logger:        infraLogger,
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
