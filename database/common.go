package database

import (
	"reflect"

	. "github.com/gyf841010/pz-infra-new/logging"
)

func SaveModelWithChan(t interface{}, isCreate bool, ch chan error) error {
	db := GetDB()
	if isCreate {
		if err := db.Create(t).Error; err != nil {
			Log.Error("Failed to Create Model in Batch Mode", With("model", t), WithError(err))
			ch <- err
			return err
		}
	} else {
		if err := db.Save(t).Error; err != nil {
			Log.Error("Failed to Save Model in Batch Mode", With("model", t), WithError(err))
			ch <- err
			return err
		}
	}
	ch <- nil
	return nil
}

func EmptyResponseChannel(ch chan error) {
	ch <- nil
}

func BatchSaveModels(t interface{}, isCreate bool) error {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		list := reflect.ValueOf(t)

		Log.Info("BatchSaveModels, len:", With("len", list.Len()))
		i := 0
		maxProcess := 10
		ch := make(chan error, maxProcess)
		for {
			for j := 0; j < maxProcess; j++ {
				if i+j >= list.Len() {
					go EmptyResponseChannel(ch)
				} else {
					go SaveModelWithChan(list.Index(i+j).Interface(), true, ch)
				}
			}

			for j := 0; j < maxProcess; j++ {
				err := <-ch
				if err != nil {
					return err
				}
			}

			i += maxProcess
			if i >= list.Len() {
				break
			}
		}
	}
	return nil
}

func ExecSql(sql string, args ...interface{}) error {
	db := GetDB()
	if err := db.Exec(sql, args...).Error; err != nil {
		Log.Error("Failed to Exec Sql", With("sql", sql), With("args", args), WithError(err))
		return err
	}
	return nil
}

func AddModel(t interface{}) error {
	db := GetDB()
	if err := db.Create(t).Error; err != nil {
		Log.Error("Failed to AddModel", With("model", t), WithError(err))
		return err
	}
	return nil
}
