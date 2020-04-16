package def

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type sample struct {
	B bool `def:"true"`

	I int `def:"123321"`
	I8 int8 `def:"-125"`
	I16 int16 `def:"-31000"`
	I32 int32 `def:"-123"`
	I64 int64 `def:"-1"`

	Ui uint `def:"123456790"`
	Ui8 uint8 `def:"0b10011011"`
	Ui16 uint16 `def:"65000"`
	Ui32 uint32 `def:"0xFFFFFFFF"`
	Ui64 uint64 `def:"0o77777777"`

	F32 float32 `def:"-3.141592"`
	F64 float64 `def:"3.141592653589"`

	C64 complex64 `def:"321,123"`
	C128 complex128 `def:"123,321"`

	Str string `def:"Hello World"`


	Pb *bool `def:"true"`

	Pi *int `def:"123321"`
	Pi8 *int8 `def:"-125"`
	Pi16 *int16 `def:"-31000"`
	Pi32 *int32 `def:"-123"`
	Pi64 *int64 `def:"-1"`

	Pui *uint `def:"123456790"`
	Pui8 *uint8 `def:"0b10011011"`
	Pui16 *uint16 `def:"65000"`
	Pui32 *uint32 `def:"0xFFFFFFFF"`
	Pui64 *uint64 `def:"0o77777777"`

	Pf32 *float32 `def:"-3.141592"`
	Pf64 *float64 `def:"3.141592653589"`

	Pc64 *complex64 `def:"321,123"`
	Pc128 *complex128 `def:"123,321"`

	Pstr *string `def:"Hello World"`

	Pnull *int `def:"-"`
	Pnil *int `def:"-"`

	ArrInt [3]int `def:"[1,2,3]"`
	SliInt []int `def:"[1,2,3,4,5,6,7,8,9,10]"`

	ChanInt chan int `def:"0"`

	Map map[string]int `def:"{\"math\":100,\"english\":30,\"some\":999}"`
}

type nestedSample struct {
	Name    string  `def:"this is nested sample"`
	Sample  sample  `def:"dive"`
	Psample *sample `def:"dive"`
}

type jsonSample struct {
	Name string `def:"this is json struct sample"`

	Json  nestedJsonSample  `def:"{\"age\":25,\"name\":\"rebirth lee\"}"`
	Pjson *nestedJsonSample `def:"{\"age\":52,\"name\":\"lee rebirth\"}"`
}

type nestedJsonSample struct {
	Age int `json:"age"`
	Name string `json:"name"`
}

type cycleErrorSample struct {
	Name string             `def:"this is cycle sample"`
	Cycle *cycleErrorSample `def:"dive"`
}

type boolErrorSample struct {
	Bool bool `def:"fail parse"`
}

type intErrorSample struct {
	Number int `def:"fail parse"`
}

type uintErrorSample struct {
	Number uint `def:"-1"`
}

type floatErrorSample struct {
	Number float64 `def:"fail parse"`
}

type complexErrorSample struct {
	Number complex128 `def:"1,2,3,4,5,6"`
}

type complexFailParseError1Sample struct {
	Number complex128 `def:"fail parse,321"`
}

type complexFailParseError2Sample struct {
	Number complex128 `def:"123,fail parse"`
}

type chanErrorSample struct {
	Chan chan int `def:"-1"`
}

type chanFailParseErrorSample struct {
	Chan chan int `def:"fail parse"`
}


func TestNew(t *testing.T) {
	iface, err := New(sample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*sample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			checkSample(t, sample)
		}
	}


	iface, err = New(nestedSample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*nestedSample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is nested sample")
			if assert.NotNil(t, sample.Psample) {
				checkSample(t, sample.Psample)
			}
			checkSample(t, &sample.Sample)
		}
	}

	iface, err = New(jsonSample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*jsonSample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is json struct sample")

			assert.Equal(t, sample.Json.Name, "rebirth lee")
			assert.Equal(t, sample.Json.Age, 25)

			if assert.NotNil(t, sample.Pjson) {
				assert.Equal(t, sample.Pjson.Name, "lee rebirth")
				assert.Equal(t, sample.Pjson.Age, 52)
			}
		}
	}

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
		iface, err = New(errorCheckList[i])
		assert.Nil(t, iface)
		assert.Error(t, err)
	}
}

func TestMustNew(t *testing.T) {

	{
		sample := MustNew(sample{}).(*sample)
		if assert.NotNil(t, sample) {
			checkSample(t, sample)
		}
	}

	{
		sample := MustNew(nestedSample{}).(*nestedSample)
		if assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is nested sample")
			if assert.NotNil(t, sample.Psample) {
				checkSample(t, sample.Psample)
			}
			checkSample(t, &sample.Sample)
		}
	}

	{
		sample := MustNew(jsonSample{}).(*jsonSample)
		if assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is json struct sample")

			assert.Equal(t, sample.Json.Name, "rebirth lee")
			assert.Equal(t, sample.Json.Age, 25)

			if assert.NotNil(t, sample.Pjson) {
				assert.Equal(t, sample.Pjson.Name, "lee rebirth")
				assert.Equal(t, sample.Pjson.Age, 52)
			}
		}
	}

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
			MustNew(errorCheckList[i])
		})
	}
}

func TestJustNew(t *testing.T) {
	iface, err := JustNew(sample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*sample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			checkSample(t, sample)
		}
	}


	iface, err = JustNew(nestedSample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*nestedSample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is nested sample")
			if assert.NotNil(t, sample.Psample) {
				checkSample(t, sample.Psample)
			}
			checkSample(t, &sample.Sample)
		}
	}

	iface, err = JustNew(jsonSample{})
	if assert.NoError(t, err) {
		sample, ok := iface.(*jsonSample)
		if assert.True(t, ok) && assert.NotNil(t, sample) {
			assert.Equal(t, sample.Name, "this is json struct sample")

			assert.Equal(t, sample.Json.Name, "rebirth lee")
			assert.Equal(t, sample.Json.Age, 25)

			if assert.NotNil(t, sample.Pjson) {
				assert.Equal(t, sample.Pjson.Name, "lee rebirth")
				assert.Equal(t, sample.Pjson.Age, 52)
			}
		}
	}

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
		iface, err = JustNew(errorCheckList[i])
		assert.NotNil(t, iface)
		assert.Error(t, err)
	}
}

func checkSample(t *testing.T, sample *sample) {
	assert.True(t, sample.B)

	assert.Equal(t, sample.I, 123321)
	assert.Equal(t, sample.I8, int8(-125))
	assert.Equal(t, sample.I16, int16(-31000))
	assert.Equal(t, sample.I32, int32(-123))
	assert.Equal(t, sample.I64, int64(-1))

	assert.Equal(t, sample.Ui, uint(123456790))
	assert.Equal(t, sample.Ui8, uint8(0b10011011))
	assert.Equal(t, sample.Ui16, uint16(65000))
	assert.Equal(t, sample.Ui32, uint32(0xFFFFFFFF))
	assert.Equal(t, sample.Ui64, uint64(0o77777777))

	assert.Equal(t, sample.F32, float32(-3.141592))
	assert.Equal(t, sample.F64, 3.141592653589)

	assert.Equal(t, sample.C64, complex64(321+123i))
	assert.Equal(t, sample.C128, 123+321i)

	assert.Equal(t, sample.Str, "Hello World")

	//pointer
	assert.NotNil(t, sample.Pb)

	assert.NotNil(t, sample.Pi)
	assert.NotNil(t, sample.Pi8)
	assert.NotNil(t, sample.Pi16)
	assert.NotNil(t, sample.Pi32)
	assert.NotNil(t, sample.Pi64)

	assert.NotNil(t, sample.Pui)
	assert.NotNil(t, sample.Pui8)
	assert.NotNil(t, sample.Pui16)
	assert.NotNil(t, sample.Pui32)
	assert.NotNil(t, sample.Pui64)

	assert.NotNil(t, sample.Pf32)
	assert.NotNil(t, sample.Pf64)

	assert.NotNil(t, sample.Pc64)
	assert.NotNil(t, sample.Pc128)

	assert.NotNil(t, sample.Pstr)

	assert.Nil(t, sample.Pnull)
	assert.Nil(t, sample.Pnil)

	// equal pointer value
	assert.Equal(t, sample.B, *sample.Pb)

	assert.Equal(t, sample.I, *sample.Pi)
	assert.Equal(t, sample.I8, *sample.Pi8)
	assert.Equal(t, sample.I16, *sample.Pi16)
	assert.Equal(t, sample.I32, *sample.Pi32)
	assert.Equal(t, sample.I64, *sample.Pi64)

	assert.Equal(t, sample.Ui, *sample.Pui)
	assert.Equal(t, sample.Ui8, *sample.Pui8)
	assert.Equal(t, sample.Ui16, *sample.Pui16)
	assert.Equal(t, sample.Ui32, *sample.Pui32)
	assert.Equal(t, sample.Ui64, *sample.Pui64)

	assert.Equal(t, sample.F32, *sample.Pf32)
	assert.Equal(t, sample.F64, *sample.Pf64)

	assert.Equal(t, sample.C64, *sample.Pc64)
	assert.Equal(t, sample.C128, *sample.Pc128)

	assert.Equal(t, sample.Str, *sample.Pstr)

	// special type
	assert.NotNil(t, sample.ArrInt)
	assert.Equal(t, sample.ArrInt, [3]int{1,2,3})

	assert.NotNil(t, sample.SliInt)
	assert.Equal(t, sample.SliInt, []int{1,2,3,4,5,6,7,8,9,10})

	assert.NotNil(t, sample.ChanInt)
	assert.Len(t, sample.ChanInt, 0)

	assert.NotNil(t, sample.Map)
	assert.Equal(t, sample.Map, map[string]int{
		"math": 100,
		"english": 30,
		"some": 999,
	})
}