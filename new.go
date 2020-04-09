package def

import "reflect"

func New(i interface{}) (interface{}, error) {
	ref := reflect.New(reflect.TypeOf(i))
	defer callInit(ref)
	if err := initStruct(ref.Elem(), maybeInit, make(map[reflect.Type]bool)); err != nil {
		return nil, err
	}

	return ref.Interface(), nil
}

func MustNew(i interface{}) interface{} {
	ref := reflect.New(reflect.TypeOf(i))
	defer callInit(ref)
	if err := initStruct(ref.Elem(), maybeInit, make(map[reflect.Type]bool)); err != nil {
		panic(err)
	}

	return ref.Interface()
}

func JustNew(i interface{}) (interface{}, error) {
	ref := reflect.New(reflect.TypeOf(i))
	defer callInit(ref)
	err := initStruct(ref.Elem(), justInit, make(map[reflect.Type]bool))
	return ref.Interface(), err
}
