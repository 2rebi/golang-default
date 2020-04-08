package def

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagNameDefault = "def"
)

type structInitSelector func(v reflect.Value) error

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
	return initialize(v.Elem(), maybeInit)
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
	return initialize(v.Elem(), justInit)
}

func initialize(v reflect.Value, selector structInitSelector) error {
	defer callInit(v)
	if !v.CanSet() {
		return nil
	}

	return selector(v)
}

func justInit(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val, ok := t.Field(i).Tag.Lookup(tagNameDefault); val != "-" && ok {
			if err := initField(v.Field(i), val, justInit); err != nil {
				// 묶음 후 리턴
			}
		}
	}

	return nil
}

func maybeInit(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val:= t.Field(i).Tag.Get(tagNameDefault); val != "-" {
			if err := initField(v.Field(i), val, maybeInit); err != nil {
				return err
			}
		}
	}

	return nil
}


func initField(v reflect.Value, val string, selector structInitSelector) error {
	defer callInit(v)
	if !v.CanSet() {
		return nil
	}

	switch k := v.Kind(); k {
	case reflect.Invalid:
		return nil
	case reflect.Ptr:
		elem := v.Elem()
		if elem.Kind() == reflect.Invalid {
			v.Set(reflect.New(v.Type().Elem()))
			elem = v.Elem()
		}
		return initField(elem, val, selector)
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
	case reflect.Array, reflect.Slice,
			reflect.Interface, reflect.Map, reflect.Struct:
		if k == reflect.Struct && val == "" {
			return initialize(v, selector)
		}
		ref := reflect.New(v.Type())
		if err := json.Unmarshal([]byte(val), ref.Interface()); err != nil {
			return err
		}
		v.Set(ref.Elem())
	case reflect.Chan:
		if val == "" {
			return nil
		} else if strings.HasPrefix(val, "-") {
			return errors.New("negative buffer size, param must be 0 or more")
		} else if i, err := strconv.Atoi(val); err != nil {
			return err
		} else {
			v.Set(reflect.MakeChan(v.Type(), i))
		}
	}

	//TODO
	//Uintptr == uintptr
	//Func
	//UnsafePointer == unsafe.Pointer

	return nil
}

func callInit(i interface{}) {
	if init, ok := i.(Initializer); ok {
		init.Init()
	}
}