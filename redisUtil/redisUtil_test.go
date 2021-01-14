package redisUtil

import (
	"testing"
	"time"

	"github.com/gyf841010/pz-infra-new/tests/base"

	"math/rand"

	"github.com/astaxie/beego"
	"github.com/pborman/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	base.InitConfigFile("roav/api/conf/app.conf")
}

type TestStruc struct {
	Id        string `orm:"column(uuid);pk"`
	OrderId   string `orm:"column(order_id);not null"`
	DealType  int    `orm:"column(deal_type);"`
	Status    int    `orm:"column(status);null"`
	UserId    string `orm:"column(user_id);"`
	DeviceId  string `orm:"column(device_id);"`
	SlotId    string `orm:"column(slot_id);"`
	BatteryId string `orm:"column(battery_id);"`
}

//func TestAddGeoIndex(t *testing.T) {
//	err := AddGeoIndex("testGeo","geo1",17.333,35.666)
//
//	Convey("err should be nil", func() {
//		So(err, ShouldBeNil)
//	})
//}

func TestSetAndSetStruct(t *testing.T) {
	v1 := TestStruc{
		Id:       "mockTranSactionID",
		DeviceId: "mockDeviceID",
		DealType: 1,
		UserId:   "userID",
	}

	err := SetObject("mockTranSactionID", &v1)
	Convey("redis set object with no err", t, func() {
		So(err, ShouldBeNil)
	})

	v2 := TestStruc{}
	err = GetObject("mockTranSactionID", &v2)
	if err == nil {
		beego.Trace("-------------redis value------------")
		beego.Trace(v2)
	}
	Convey("redis get object with no error and every field value are correct.", t, func() {
		Convey("err should be nil", func() {
			So(err, ShouldBeNil)
		})
		Convey("Id is correct", func() {
			So(v2.Id, ShouldEqual, "mockTranSactionID")
		})
		Convey("DeviceID is correct", func() {
			So(v2.DeviceId, ShouldEqual, "mockDeviceID")
		})
		Convey("DealType is correct", func() {
			So(v2.DealType, ShouldEqual, 1)
		})
		Convey("UserId is correct", func() {
			So(v2.UserId, ShouldEqual, "userID")
		})
	})
}

func TestSetObjectWithExpire(t *testing.T) {
	v := TestStruc{Id: "abc"}
	Convey("test set object with expire", t, func() {
		key := "testKey"
		err := SetObjectWithExpire(key, &v, 3)
		So(err, ShouldBeNil)

		result := TestStruc{}
		err = GetObject(key, &result)
		So(err, ShouldBeNil)
		So(result.Id, ShouldEqual, v.Id)

		time.Sleep(4 * time.Second)
		So(Exists(key), ShouldBeFalse)
	})
}

func TestGetNotExistObject(t *testing.T) {
	result := TestStruc{}
	err := GetObject("NotExistingKey", &result)
	Convey("test Get Not Existing Object", t, func() {
		So(err, ShouldBeNil)
		So(result.DeviceId, ShouldEqual, "")
		So(result.BatteryId, ShouldEqual, "")
		So(result.DealType, ShouldEqual, 0)
	})
}

func TestGetNotExistString(t *testing.T) {
	result, err := GetString("NotExistingSring")
	Convey("Test Get Not Exist String", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldEqual, "")
	})
}

func TestGetNotExistingObject(t *testing.T) {
	v := TestStruc{}
	err := GetObject("NOT_EXISTING_KEY", &v)
	Convey("redis get with not existing key should not return error and result is nil", t, func() {
		So(err, ShouldBeNil)
		So(v.Id, ShouldBeEmpty)
	})

}

func TestGetSetString(t *testing.T) {
	Delete("TestGetSetString_key")
	e := Exists("TestGetSetString_key")
	Convey("not exists before set", t, func() {
		So(e, ShouldEqual, false)
	})
	err := SetStringWithExpire("TestGetSetString_key", "TestGetSetString_value", 2)
	Convey("redis set string with no err", t, func() {
		So(err, ShouldBeNil)
	})
	e = Exists("TestGetSetString_key")
	Convey("exists after set", t, func() {
		So(e, ShouldEqual, true)
	})
	v, err := GetString("TestGetSetString_key")
	Convey("redis get string with no err", t, func() {
		So(err, ShouldBeNil)
	})
	Convey("get string value should be correct", t, func() {
		So(v, ShouldEqual, "TestGetSetString_value")
	})

	time.Sleep(time.Second * 4)
	e = Exists("TestGetSetString_key")
	Convey("not exists after timeout", t, func() {
		So(e, ShouldEqual, false)
	})
}

func TestGetKeysByPrefix(t *testing.T) {
	prefix := "testPrefix_"
	SetString(prefix+"1", "1")
	SetString(prefix+"2", "2")
	Convey("should return two keys with existed prefix", t, func() {
		keys, err := GetKeysByPrefix(prefix)
		So(err, ShouldBeNil)
		So(len(keys), ShouldEqual, 2)
	})

	Delete(prefix + "1")
	Delete(prefix + "2")

	notExistPrefix := "mockPrefix_balabala"
	Convey("should return zero with not existed prefix", t, func() {
		keys, err := GetKeysByPrefix(notExistPrefix)
		So(err, ShouldBeNil)
		So(len(keys), ShouldBeZeroValue)
	})
}

func TestSetNotExist(t *testing.T) {
	prefix := "test_set_not_exist"

	Delete(prefix + "test")

	result, err := SetStringIfNotExist(prefix+"test", "haha", 2)
	Convey("first set should return 1 and no error", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldEqual, true)
	})

	result, err = SetStringIfNotExist(prefix+"test", "hahahaha", 2)
	afterString, _ := GetString(prefix + "test")
	Convey("second set should return 0 and no error", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldEqual, false)
		So(afterString, ShouldEqual, "haha")
	})
}

type mockStruct struct {
	Id   int
	Name string
}

func TestSetAndGetValue(t *testing.T) {
	mockKey := "mock_key"
	Convey("test SetValue and GetValue for int", t, func() {
		Delete(mockKey) // clear key

		mockValue := 1024
		err := SetValue(mockKey, mockValue)
		So(err, ShouldBeNil)

		var v1 int
		err = GetValue(mockKey, &v1)
		So(err, ShouldBeNil)
		So(v1, ShouldEqual, mockValue)

		err = SetValue(mockKey, &mockValue) // pointer
		So(err, ShouldBeNil)

		var v2 int
		err = GetValue(mockKey, &v2)
		So(err, ShouldBeNil)
		So(v2, ShouldEqual, mockValue)
	})

	Convey("test SetValue and GetValue for string", t, func() {
		Delete(mockKey) // clear key

		mockValue := "1024"
		err := SetValue(mockKey, mockValue)
		So(err, ShouldBeNil)

		var v1 string
		err = GetValue(mockKey, &v1)
		So(err, ShouldBeNil)
		So(v1, ShouldEqual, mockValue)

		err = SetValue(mockKey, &mockValue) // pointer
		So(err, ShouldBeNil)

		var v2 string
		err = GetValue(mockKey, &v2)
		So(err, ShouldBeNil)
		So(v2, ShouldEqual, mockValue)
	})

	Convey("test SetValue and GetValue for bool", t, func() {
		Delete(mockKey) // clear key

		mockValue := true
		err := SetValue(mockKey, mockValue)
		So(err, ShouldBeNil)

		var v1 bool
		err = GetValue(mockKey, &v1)
		So(err, ShouldBeNil)
		So(v1, ShouldEqual, mockValue)

		err = SetValue(mockKey, &mockValue) // pointer
		So(err, ShouldBeNil)

		var v2 bool
		err = GetValue(mockKey, &v2)
		So(err, ShouldBeNil)
		So(v2, ShouldEqual, mockValue)
	})

	Convey("test SetValue and GetValue for struct", t, func() {
		Delete(mockKey) // clear key

		mockValue := mockStruct{Id: rand.Int()} // struct
		err := SetValue(mockKey, mockValue)
		So(err, ShouldBeNil)

		var v1 mockStruct
		err = GetValue(mockKey, &v1)
		So(err, ShouldBeNil)
		So(v1.Id, ShouldEqual, mockValue.Id)

		err = SetValue(mockKey, &mockValue) // pointer
		So(err, ShouldBeNil)

		var v2 mockStruct
		err = GetValue(mockKey, &v2)
		So(err, ShouldBeNil)
		So(v2.Id, ShouldEqual, mockValue.Id)
	})

	Convey("test SetValue and GetValue for nil-pointer", t, func() {
		Delete(mockKey) // clear key

		var v *mockStruct
		err := SetValue(mockKey, v) // nil-pointer
		So(err, ShouldEqual, errValueIsNil)
	})

	Convey("test SetValue and GetValue for nil", t, func() {
		Delete(mockKey) // clear key

		err := SetValue(mockKey, nil)
		So(err, ShouldEqual, errValueIsNil)

		err = GetValue(mockKey, nil)
		So(err, ShouldEqual, errValueIsNotPointer)
	})

	Convey("test SetValue and GetValue for blank key", t, func() {
		blankKey := ""
		err := SetValue(blankKey, nil)
		So(err, ShouldEqual, errKeyIsBlank)

		err = GetValue(blankKey, nil)
		So(err, ShouldEqual, errKeyIsBlank)
	})

	Convey("test GetValue with not exist key", t, func() {
		notExistKey := uuid.New()
		Delete(notExistKey)

		var v interface{}
		err := GetValue(notExistKey, &v)
		So(err, ShouldEqual, ErrKeyNotFound)
	})

	Convey("test SetValue with expire seconds", t, func() {
		Delete(mockKey) // clear key

		expireSeconds := 3
		sleepSeconds := time.Duration(expireSeconds+2) * time.Second
		mockValue := 1024
		err := SetValue(mockKey, mockValue, expireSeconds)
		So(err, ShouldBeNil)

		var v1 int
		err = GetValue(mockKey, &v1)
		So(err, ShouldBeNil)
		So(v1, ShouldEqual, mockValue)

		time.Sleep(sleepSeconds)

		var v2 int
		err = GetValue(mockKey, &v2)
		So(err, ShouldEqual, ErrKeyNotFound)
		So(v2, ShouldEqual, 0)
	})

}

var setAndGetStringsTest = []struct {
	key    string
	ssSet  []string
	ex     []int
	errSet error
	ssGet  []string
	errGet error
}{
	{"", []string{}, nil, errKeyIsBlank, nil, errKeyIsBlank},
	{"foo", nil, nil, nil, []string{}, nil},
	{"foo", []string{}, nil, nil, []string{}, nil},
	{"foo", []string{""}, nil, nil, []string{""}, nil},
	{"foo", []string{"a"}, nil, nil, []string{"", "a"}, nil},
	{"foo", []string{"a", "b"}, nil, nil, []string{"", "a", "b"}, nil},
	{"foo", []string{"c", "d", "d"}, nil, nil, []string{"", "a", "b", "c", "d"}, nil},
	{"foo", []string{"d"}, []int{10}, nil, []string{"", "a", "b", "c", "d"}, nil},
}

func TestSetAndGetStrings(t *testing.T) {
	for _, tt := range setAndGetStringsTest {
		Convey("Test SetStrings", t, func() {
			errSet := SetStrings(tt.key, tt.ssSet, tt.ex...)
			So(errSet, ShouldEqual, tt.errSet)

			ss, errGet := GetStrings(tt.key)
			So(errGet, ShouldEqual, tt.errGet)
			So(len(ss), ShouldEqual, len(tt.ssGet))
		})
	}
}
