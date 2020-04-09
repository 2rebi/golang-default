package def

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagNameDefault = "def"
)

type structInitSelector func(v reflect.Value, visitedStruct map[reflect.Type]bool) error

func checkPtr(i interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return v, errors.New("param must be pointer")
	}

	return v, nil
}

func Init(i interface{}) error {
	v, err := checkPtr(i)
	if err != nil {
		return err
	}
	return initStruct(v.Elem(), maybeInit, make(map[reflect.Type]bool))
}

func MustInit(i interface{}) {
	if err := Init(i); err != nil {
		panic(err)
	}
}

func JustInit(i interface{}) error {
	v, err := checkPtr(i)
	if err != nil {
		return err
	}
	return initStruct(v.Elem(), justInit, make(map[reflect.Type]bool))
}

func initStruct(v reflect.Value, selector structInitSelector, visitedStruct map[reflect.Type]bool) error {
	defer callInit(v)
	if !v.CanSet() {
		return nil
	}

	t := v.Type()
	if visitedStruct[t] {
		return fmt.Errorf("struct type \"%s\" is cycle", t.Name())
	}

	visitedStruct[t] = true
	defer delete(visitedStruct, t)
	return selector(v, visitedStruct)
}

func justInit(v reflect.Value, visitedStruct map[reflect.Type]bool) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val, ok := t.Field(i).Tag.Lookup(tagNameDefault); val != "-" {
			if err := initField(v, v.Field(i), val, ok, justInit, visitedStruct); err != nil {
				// TODO 묶음 후 리턴
			}
		}
	}

	return nil
}

func maybeInit(v reflect.Value, visitedStruct map[reflect.Type]bool) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val, ok := t.Field(i).Tag.Lookup(tagNameDefault); val != "-" {
			if err := initField(v, v.Field(i), val, ok, maybeInit, visitedStruct); err != nil {
				return err
			}
		}
	}

	return nil
}


func initField(structVal reflect.Value, v reflect.Value, val string, lookupOk bool, selector structInitSelector, visitedStruct map[reflect.Type]bool) error {
	if !v.CanSet() || !lookupOk {
		return nil
	}

	switch k := v.Kind(); k {
	case reflect.Invalid:
		return nil
	case reflect.Ptr:
		elem := v.Elem()
		//if val == "nil" || val == "null" {
		//	v.Set(reflect.Zero(v.Type()))
		//	return nil
		//} else
		if elem.Kind() == reflect.Invalid {
			v.Set(reflect.New(v.Type().Elem()))
			elem = v.Elem()
		}
		defer callInit(v)
		return initField(structVal, elem, val, lookupOk, selector, visitedStruct)
	case reflect.String:
		v.SetString(val)
	case reflect.Bool:
		if b, err := strconv.ParseBool(val); err != nil {
			return err
		} else {
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(getBaseValue(val)); err != nil {
			return err
		} else {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(getBaseValue(val)); err != nil {
			return err
		} else {
			v.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(val, 0); err != nil {
			return err
		} else {
			v.SetFloat(f)
		}
	case reflect.Complex64, reflect.Complex128:
		vals := strings.Split(val,",")
		if len(vals) != 2 {
			return errors.New("only two args")
		}
		if c, err := getStringToComplex(vals[0], vals[1]); err != nil {
			return err
		} else {
			v.SetComplex(c)
		}
		//TODO Fix: 구조체형태의 배열, 슬라이스 처리 필요
	case reflect.Array, reflect.Slice,
			reflect.Interface, reflect.Map, reflect.Struct:
		if k == reflect.Struct && val == "dive" {
			return initStruct(v, selector, visitedStruct)
		}
		ref := reflect.New(v.Type())
		if err := json.Unmarshal([]byte(val), ref.Interface()); err != nil {
			return err
		}
		v.Set(ref.Elem())
	case reflect.Chan:
		if strings.HasPrefix(val, "-") {
			return errors.New("negative buffer size, param must be 0 or more")
		} else if i, err := strconv.Atoi(val); err != nil {
			return err
		} else {
			v.Set(reflect.MakeChan(v.Type(), i))
		}
	case reflect.Func:
		srcFunc, ok := funcMap[val]
		if !ok {
			return fmt.Errorf("don't setup function for key : \"%s\"", val)
		} else if !structVal.CanAddr() {
			return errors.New("function initialize failed, because can't access address of struct")
		}

		funcIface := reflect.ValueOf(srcFunc)
		self := structVal.Addr().Convert(structVal.Addr().Type())
		srcVal := funcIface.Call([]reflect.Value{self})[0].Elem()

		srcType := srcVal.Type()
		if srcType.Kind() != reflect.Func {
			return errors.New("return value must be function type")
		}
		vType := v.Type()
		if vType.NumIn() != srcType.NumIn() {
			return errors.New("args count not equal")
		} else if vType.NumOut() != srcType.NumOut() {
			return errors.New("returns count not equal")
		}

		for i, cnt := 0, vType.NumIn(); i < cnt; i++ {
			if vType.In(i) != srcType.In(i) {
				return fmt.Errorf("(argument at %d) wrong argument type, dest func arg type \"%s::%s\", src func arg type \"%s::%s\"", i,
					vType.In(i).PkgPath(), vType.In(i).Name(), srcType.In(i).PkgPath(), srcType.In(i).Name())
			}
		}

		for i, cnt := 0, vType.NumOut(); i < cnt; i++ {
			if vType.Out(i) != srcType.Out(i) {
				return fmt.Errorf("(argument at %d) wrong argument type, dest func arg type \"%s::%s\", src func arg type \"%s::%s\"", i,
					vType.Out(i).PkgPath(), vType.Out(i).Name(), srcType.Out(i).PkgPath(), srcType.Out(i).Name())
			}
		}

		v.Set(srcVal)
	}


	return nil
}

func callInit(v reflect.Value) {
	if init, ok := v.Interface().(Initializer); ok && init != nil {
		init.Init()
	}
}