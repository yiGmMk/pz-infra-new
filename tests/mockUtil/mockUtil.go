package mockUtil

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/stretchr/testify/mock"
)

type mockObjectInterface interface {
	On(string, ...interface{}) *mock.Call
}

func SetupAllMethods(mockObj interface{}, interfaceType reflect.Type) interface{} {
	if mockObj == nil || interfaceType == nil {
		panic("mockObj is nil or interfaceType is nil")
	}
	mockTypeObject := mockObj.(mockObjectInterface)

	if interfaceType.Kind() == reflect.Ptr {
		interfaceType = interfaceType.Elem()
	}
	if interfaceType.Kind() == reflect.Ptr {
		panic("cannot get interface type object from interfaceType")
	}
	if interfaceType.Kind() != reflect.Interface {
		panic("interfaceType is pointer to pointer, not supported")
	}

	for i := 0; i < interfaceType.NumMethod(); i++ {
		method := interfaceType.Method(i)
		arguments := make([]interface{}, 0, method.Type.NumIn())
		for j := 0; j < method.Type.NumIn(); j++ {
			arguments = append(arguments, mock.Anything)
		}
		returns := make([]interface{}, 0, method.Type.NumOut())
		for j := 0; j < method.Type.NumOut(); j++ {
			returns = append(returns, reflect.Zero(method.Type.Out(j)).Interface())
		}
		mockTypeObject.On(method.Name, arguments...).Return(returns...)
	}
	return mockObj
}

// remove all setups of THE method and add the new one
func ReplaceMockMethods(mockObj, method interface{}, arguments ...interface{}) *mock.Call {
	if mockObj == nil || method == nil {
		panic("mockObj or method is nil")
	}
	mockValue := reflect.ValueOf(mockObj)
	if mockValue.Kind() == reflect.Ptr {
		mockValue = mockValue.Elem()
	}
	if mockValue.Kind() == reflect.Ptr {
		panic("mockObj is pointer to pointer, not supported")
	}
	mockFieldValue := mockValue.FieldByName("Mock")
	if !mockFieldValue.IsValid() {
		panic("not found Mock field in mockObj, make sure mockObj type inherit mock.Mock")
	}
	mockFieldPtr := mockFieldValue.Addr().Interface().(*mock.Mock)

	methodName := getFunctionName(method)
	argNum := getFunctionArgumentCount(method)
	for i := len(arguments); i < argNum; i++ {
		arguments = append(arguments, mock.Anything)
	}
	for i := len(mockFieldPtr.ExpectedCalls) - 1; i >= 0; i-- {
		call := mockFieldPtr.ExpectedCalls[i]
		if call.Method == methodName {
			mockFieldPtr.ExpectedCalls = append(mockFieldPtr.ExpectedCalls[:i], mockFieldPtr.ExpectedCalls[i+1:]...)
		}
	}
	return mockFieldPtr.On(methodName, arguments...)
}

func getFunctionName(method interface{}) string {
	methodName := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	if index := strings.LastIndex(methodName, "."); index != -1 {
		methodName = methodName[index+1:]
	}
	if index := strings.LastIndex(methodName, "-"); index != -1 {
		methodName = methodName[:index]
	}
	return methodName
}

func getFunctionArgumentCount(method interface{}) int {
	return reflect.ValueOf(method).Type().NumIn()
}
