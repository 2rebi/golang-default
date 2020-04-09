package def

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type initSample struct {
	number int
	name string
}

func (i *initSample) Init() {
	i.name = "init test"
	i.number = 1111
}

func Test_checkPtr(t *testing.T) {
	var typeInt int
	typeInt = 123
	v, err := checkPtr(typeInt)
	if assert.Error(t, err) {
		valueOfInt := reflect.ValueOf(typeInt)
		assert.Equal(t, v.Type(), valueOfInt.Type())
		assert.Equal(t, v.Int(), valueOfInt.Int())
	}

	v, err = checkPtr(&typeInt)
	if assert.NoError(t, err) {
		valueOfInt := reflect.ValueOf(&typeInt)
		assert.Equal(t, v.Elem().Type(), valueOfInt.Elem().Type())
		assert.Equal(t, v.Pointer(), valueOfInt.Pointer())
	}
}

func Test_callInit(t *testing.T) {
	i := initSample{}
	callInit(reflect.ValueOf(&i))
	assert.Equal(t, i.name, "init test")
	assert.Equal(t, i.number, 1111)
}

func Test_initField(t *testing.T) {
	assert.NoError(t, initField(reflect.Value{}, "", false, nil, nil))
	assert.NoError(t, initField(reflect.ValueOf(0), "", true, nil, nil))
}


func TestInit(t *testing.T) {

	{
		sample := sample{}
		err := Init(&sample)
		if assert.NoError(t, err) {
			checkSample(t, &sample)
		}
	}

	{
		sample := sample{}
		err := Init(sample)
		assert.Error(t, err)
	}

	{
		sample := nestedSample{}
		err := Init(&sample)
		if assert.NoError(t, err) {
			assert.Equal(t, sample.Name, "this is nested sample")
			if assert.NotNil(t, sample.Psample) {
				checkSample(t, sample.Psample)
			}
			checkSample(t, &sample.Sample)
		}
	}

	{
		sample := jsonSample{}
		err := Init(&sample)
		if assert.NoError(t, err) {
			assert.Equal(t, sample.Name, "this is json struct sample")

			assert.Equal(t, sample.Json.Name, "rebirth lee")
			assert.Equal(t, sample.Json.Age, 25)

			if assert.NotNil(t, sample.Pjson) {
				assert.Equal(t, sample.Pjson.Name, "lee rebirth")
				assert.Equal(t, sample.Pjson.Age, 52)
			}
		}
	}

	{
		errorCheckList := []interface{}{
			cycleErrorSample{},
			boolErrorSample{},
			intErrorSample{},
			uintErrorSample{},
			floatErrorSample{},
			complexErrorSample{},
			complexFailParseError1Sample{},
			complexFailParseError2Sample{},
			chanErrorSample{},
			chanFailParseErrorSample{},
		}

		for i := range errorCheckList {
			err := Init(errorCheckList[i])
			assert.Error(t, err)
		}
	}
}


func TestMustInit(t *testing.T) {

	{
		sample := sample{}
		MustInit(&sample)
		checkSample(t, &sample)
	}

	{
		sample := sample{}
		assert.Panics(t, func() {
			MustInit(sample)
		})
	}

	{
		sample := nestedSample{}
		MustInit(&sample)
		assert.Equal(t, sample.Name, "this is nested sample")
		if assert.NotNil(t, sample.Psample) {
			checkSample(t, sample.Psample)
		}
		checkSample(t, &sample.Sample)
	}

	{
		sample := jsonSample{}
		MustInit(&sample)
		assert.Equal(t, sample.Name, "this is json struct sample")

		assert.Equal(t, sample.Json.Name, "rebirth lee")
		assert.Equal(t, sample.Json.Age, 25)

		if assert.NotNil(t, sample.Pjson) {
			assert.Equal(t, sample.Pjson.Name, "lee rebirth")
			assert.Equal(t, sample.Pjson.Age, 52)
		}
	}

	{
		errorCheckList := []interface{}{
			cycleErrorSample{},
			boolErrorSample{},
			intErrorSample{},
			uintErrorSample{},
			floatErrorSample{},
			complexErrorSample{},
			complexFailParseError1Sample{},
			complexFailParseError2Sample{},
			chanErrorSample{},
			chanFailParseErrorSample{},
		}

		for i := range errorCheckList {
			assert.Panics(t, func() {
				MustInit(errorCheckList[i])
			})
		}
	}
}