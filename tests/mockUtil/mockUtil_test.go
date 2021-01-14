package mockUtil

import (
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type FooInterface interface {
	DoSomething(number int) (bool, error)
	DoSomething2(number int) (*int, error)
	DoSomething3(number int) ([]int, error)
}

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) DoSomething(number int) (bool, error) {
	args := m.Called(number)
	return args.Bool(0), args.Error(1)
}
func (m *MyMockedObject) DoSomething2(number int) (*int, error) {
	args := m.Called(number)
	return args.Get(0).(*int), args.Error(1)
}
func (m *MyMockedObject) DoSomething3(number int) ([]int, error) {
	args := m.Called(number)
	return args.Get(0).([]int), args.Error(1)
}

func TestSetupAllMethods(t *testing.T) {
	Convey("TestSetupAllMethods", t, func() {
		mockObj := SetupAllMethods(&MyMockedObject{}, reflect.TypeOf((*FooInterface)(nil))).(*MyMockedObject)

		boolValue, err := mockObj.DoSomething(0)
		So(err, ShouldBeNil)
		So(boolValue, ShouldBeFalse)

		intPtr, err := mockObj.DoSomething2(1)
		So(err, ShouldBeNil)
		So(intPtr, ShouldBeNil)

		sliceValue, err := mockObj.DoSomething3(2)
		So(err, ShouldBeNil)
		So(sliceValue, ShouldBeNil)
	})
}

func TestReplaceMockMethods(t *testing.T) {
	Convey("TestReplaceMockMethods", t, func() {
		mockObj := SetupAllMethods(&MyMockedObject{}, reflect.TypeOf((*FooInterface)(nil))).(*MyMockedObject)

		ReplaceMockMethods(mockObj, mockObj.DoSomething).Return(true, nil)

		boolValue, err := mockObj.DoSomething(2)
		So(boolValue, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}
