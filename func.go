package def

import (
	"fmt"
)

var (
	funcMap = make(map[string]func(self interface{}) interface{})
)

func SetFunc(key string, fun func(self interface{}) interface{}) error {
	//if reflect.TypeOf(fun).Out(0).Kind() != reflect.Func {
	//	return errors.New("it is not function")
	//} else
	if _, ok := funcMap[key]; ok {
		return fmt.Errorf("key \"%s\" is already has", key)
	}

	funcMap[key] = fun
	return nil
}