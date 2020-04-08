package def

import "reflect"

func New(i interface{}) (interface{}, error) {
	ref := reflect.New(reflect.TypeOf(i))
	if err := initialize(ref, maybeInit); err != nil {
		return nil, err
	}

	return ref.Interface(), nil
}

func MustNew(i interface{}) interface{} {
	ref := reflect.New(reflect.TypeOf(i))
	if err := initialize(ref, maybeInit); err != nil {
		panic(err)
	}

	return ref.Interface()
}

func JustNew(i interface{}) (interface{}, error) {
	ref := reflect.New(reflect.TypeOf(i))
	err := initialize(ref, justInit)
	return ref.Interface(), err
}
