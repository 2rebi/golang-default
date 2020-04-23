package def

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	tagNameDefault = "def"

	valueDive = "dive"
	valueDiveLen = len(valueDive)
)

var (
	timeDurationType = reflect.TypeOf(time.Duration(0))
	//timeType = reflect.TypeOf(time.Time{})
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
	fieldErrors := make([]*ErrorJustInitField, 0)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val, ok := t.Field(i).Tag.Lookup(tagNameDefault); val != "-" && ok {
			if err := initField(v, v.Field(i), val, justInit, visitedStruct); err != nil {
				ft := t.Field(i)
				typeName := ft.Type.Name()
				if ft.Type.PkgPath() != "" {
					typeName = ft.Type.PkgPath() + "." + typeName
				}
				fieldErrors = append(fieldErrors, &ErrorJustInitField{
					StructName: t.PkgPath()+"."+t.Name(),
					FieldName:  ft.Name,
					FieldType:  typeName,
					TryValue:   val,
					Cause:      err,
					Target:     v.Field(i),
				})
			}
		}
	}

	if len(fieldErrors) > 0 {
		return &ErrorJustInit{Errors:fieldErrors}
	}

	return nil
}

func maybeInit(v reflect.Value, visitedStruct map[reflect.Type]bool) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if val, ok := t.Field(i).Tag.Lookup(tagNameDefault); val != "-" && ok {
			if err := initField(v, v.Field(i), val, maybeInit, visitedStruct); err != nil {
				return err
			}
		}
	}

	return nil
}


func initField(structVal reflect.Value, fieldVal reflect.Value, defVal string, selector structInitSelector, visitedStruct map[reflect.Type]bool) error {
	if !fieldVal.CanSet() {
		return nil
	}

	fieldType := fieldVal.Type()

	//special type
	switch fieldType {
	case timeDurationType:
		if d, err := time.ParseDuration(defVal); err != nil {
			return err
		} else {
			fieldVal.Set(reflect.ValueOf(d))
			return nil
		}
	}

	// primitive type
	switch k := fieldVal.Kind(); k {
	case reflect.Invalid:
		return nil
	case reflect.Ptr:
		elem := fieldVal.Elem()
		if elem.Kind() == reflect.Invalid {
			fieldVal.Set(reflect.New(fieldType.Elem()))
			elem = fieldVal.Elem()
		}
		defer callInit(fieldVal)
		return initField(structVal, elem, defVal, selector, visitedStruct)
	case reflect.String:
		fieldVal.SetString(defVal)
	case reflect.Bool:
		if b, err := strconv.ParseBool(defVal); err != nil {
			return err
		} else {
			fieldVal.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(getBaseValue(defVal)); err != nil {
			return err
		} else {
			fieldVal.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(getBaseValue(defVal)); err != nil {
			return err
		} else {
			fieldVal.SetUint(i)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(defVal, 0); err != nil {
			return err
		} else {
			fieldVal.SetFloat(f)
		}
	case reflect.Complex64, reflect.Complex128:
		vals := strings.Split(defVal,",")
		if len(vals) != 2 {
			return errors.New("only two args")
		}
		if c, err := getStringToComplex(vals[0], vals[1]); err != nil {
			return err
		} else {
			fieldVal.SetComplex(c)
		}
	case reflect.Interface:
		if defVal == "" {
			fieldVal.Set(reflect.Zero(fieldType))
		} else if err := jsonUnmarshalValue(fieldVal, defVal); err != nil {
			return err
		}
	case reflect.Map:
		if strings.HasPrefix(defVal, valueDive+"{") && strings.HasSuffix(defVal, "}") {
			keyType := fieldType.Key()
			valType := fieldType.Elem()

			fieldVal.Set(reflect.MakeMap(fieldType))

			tmp := defVal[valueDiveLen+1:len(defVal)-1]
			flag := byte(0x00)
			keyIndex := len(tmp)
			valIndex := 0
			set := func(targetVal reflect.Value, start int, end int) error {
				return initField(structVal, targetVal, tmp[start:end], selector, visitedStruct)
			}

			for i := len(tmp)-1; i >= 0; i-- {
				c := tmp[i]
				switch flag {
				case 0x00: // none
					if c == '"' {
						//flag ^= 0x01
						flag = 0x01
					} else if c == '}' {
						flag = 0x02
					}
				case 0x01: // isString
					if c == '"' {
						//flag ^= 0x01
						flag = 0x00
					}
				case 0x02: // isObject
					if c == '{' {
						flag = 0x00
					}
				}

				if flag > 0 {
					continue
				}

				if c == ':' {
					valIndex = i
				} else if c == ',' || i == 0 {
					var at int
					if i > valIndex {
						return errors.New("map default value malformed format")
					} else if i == 0 {
						at = i
					} else {
						at = i + 1
					}


					keyRef := reflect.New(keyType)
					valRef := reflect.New(valType)
					if keyErr, valErr := set(keyRef.Elem(), at, valIndex), set(valRef.Elem(), valIndex+1, keyIndex); keyErr != nil {
						return keyErr
					} else if valErr != nil {
						return valErr
					}
					fieldVal.SetMapIndex(keyRef.Elem(), valRef.Elem())
					keyIndex = i
				}
			}

		} else if err := jsonUnmarshalValue(fieldVal, defVal); err != nil {
			return err
		}
	case reflect.Struct:
		if defVal == valueDive {
			return initStruct(fieldVal, selector, visitedStruct)
		} else if err := jsonUnmarshalValue(fieldVal, defVal); err != nil {
			return err
		}
	case reflect.Slice:
		if strings.HasPrefix(defVal, valueDive+"(") {
			tmp := defVal[valueDiveLen+1:]
			endBracket := strings.Index(tmp, ")")
			if endBracket < 0 {
				return errors.New("you must close bracket ')'")
			}

			ln, cp := 0, 0
			sizes := strings.SplitN(tmp[:endBracket], ",", 2)
			switch l := len(sizes); l {
			default:
				return errors.New("must be \"dive(len)\" or \"dive(len,cap)\"")
			case 2:
				if strings.HasPrefix("-", sizes[0]) || strings.HasPrefix("-", sizes[1]) {
					return errors.New("negative size, param must be 0 or more")
				}

				var parseErr error
				ln, parseErr = strconv.Atoi(sizes[0])
				if parseErr != nil {
					return parseErr
				}

				cp, parseErr = strconv.Atoi(sizes[1])
				if parseErr != nil {
					return parseErr
				}

				if ln > cp {
					return errors.New("len larger than cap")
				}
			case 1:
				if strings.HasPrefix("-", sizes[0]) {
					return errors.New("negative size, param must be 0 or more")
				}

				var parseErr error
				ln, parseErr = strconv.Atoi(sizes[0])
				if parseErr != nil {
					return parseErr
				}

				cp = int(float64(ln) * 1.5)
			}

			val := tmp[endBracket+1:]
			if !strings.HasPrefix(val, ",") {
				return errors.New("not enough arguments")
			}
			val = val[1:]

			fieldVal.Set(reflect.MakeSlice(fieldType, ln, cp))
			for i := 0; i < ln; i++ {
				if err := initField(structVal, fieldVal.Index(i), val, selector, visitedStruct); err != nil {
					return err
				}
			}
		} else if err := jsonUnmarshalValue(fieldVal, defVal); err != nil {
			return err
		}
	case reflect.Array:
		if strings.HasPrefix(defVal, valueDive) {
			val := defVal[valueDiveLen:]
			if !strings.HasPrefix(val, ",") {
				return errors.New("not enough arguments")
			}
			val = val[1:]
			for i, cnt := 0, fieldVal.Len(); i < cnt; i++ {
				if err := initField(structVal, fieldVal.Index(i), val, selector, visitedStruct); err != nil {
					return err
				}
			}
		} else if err := jsonUnmarshalValue(fieldVal, defVal); err != nil {
			return err
		}
		return nil
	case reflect.Chan:
		if strings.HasPrefix(defVal, "-") {
			return errors.New("negative buffer size, param must be 0 or more")
		} else if i, err := strconv.Atoi(defVal); err != nil {
			return err
		} else {
			fieldVal.Set(reflect.MakeChan(fieldType, i))
		}
	case reflect.Func:
		srcFunc, ok := funcMap[defVal]
		if !structVal.CanAddr() {
			return errors.New("function initialize failed, because can't access address of struct")
		} else if !ok {
			return fmt.Errorf("not exists key : \"%s\"", defVal)
		}

		srcFuncVal := reflect.ValueOf(srcFunc)
		if srcFuncVal.Type().In(0) != structVal.Addr().Type() {
			return errors.New("function in type is wrong")
		}

		srcVal := srcFuncVal.Call([]reflect.Value{structVal.Addr()})[0].Elem()
		srcType := srcVal.Type()
		if srcType.Kind() != reflect.Func {
			return errors.New("return value must be function type")
		}
		vType := fieldType
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

		fieldVal.Set(srcVal)
	}


	return nil
}

func jsonUnmarshalValue(v reflect.Value, obj string) error {
	if !v.CanAddr() {
		return errors.New("json unmarshal fail, because can't access address of field")
	} else if err := json.Unmarshal([]byte(obj), v.Addr().Interface()); err != nil {
		return err
	}

	return nil
}

func callInit(v reflect.Value) {
	if init, ok := v.Interface().(Initializer); ok && init != nil {
		init.Init()
	}
}