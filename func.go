package def

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	funcMap = make(map[string]interface{})
)

func SetFunc(key string, fun interface{}) error {
	t := reflect.TypeOf(fun)
	if t.Kind() != reflect.Func {
		return errors.New("it is not function")
	} else if t.NumIn() != 1 {
		return errors.New("function in count must be 1")
	} else if t.In(0).Kind() != reflect.Ptr {
		return errors.New("function in param must pointer")
	} else if _, ok := funcMap[key]; ok {
		return fmt.Errorf("key \"%s\" is already has", key)
	}

	funcMap[key] = fun
	return nil
}