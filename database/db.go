package database

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yiGmMk/pz-infra-new/log"
	"github.com/yiGmMk/pz-infra-new/logging"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	globalDB                 *gorm.DB
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

// gorm V2 日志输出需要实现logger.Interface接口
type GormLogger struct {
	logger        logging.Logger // logrus封装,日志记录
	SlowThreshold time.Duration  // 慢查询阈值,用于慢查询日志记录
	SourceField   string         //
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
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

func InitDB(connectString string, infraLogger logging.Logger) error {
	var gormlog logger.Interface
	if infraLogger != nil {
		gormlog = &GormLogger{
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
