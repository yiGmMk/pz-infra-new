package database

import (
	. "github.com/gyf841010/pz-infra-new/logging"

	"gorm.io/gorm"
)

const IN_BATCH_SIZE = 1000

func GetTransactionDatabases(dbs []*gorm.DB) (*gorm.DB, bool) {
	if len(dbs) > 0 && dbs[0] != nil {
		return dbs[0], false
	}
	db := GetDB()
	return db.Begin(), true
}

func GetNonTransactionDatabases(dbs []*gorm.DB) *gorm.DB {
	if len(dbs) > 0 && dbs[0] != nil {
		return dbs[0]
	}
	return GetDB()
}

func RollbackIfError(isNewTran bool, tran *gorm.DB, err error) error {
	if isNewTran {
		if err == nil {
			if err = tran.Commit().Error; err != nil {
				return Log.Error(err.Error())
			}
		} else {
			tran.Rollback()
		}
	}
	return err
}

type TransactionFunc func(db *gorm.DB) error

func TransactionRun(method TransactionFunc, dbs ...*gorm.DB) error {
	tx, isNew := GetTransactionDatabases(dbs)
	err := method(tx)
	return RollbackIfError(isNew, tx, err)
}
