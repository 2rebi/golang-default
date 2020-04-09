package def

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getBaseValue(t *testing.T) {
	num, base, bit := getBaseValue("0xff")
	assert.Equal(t, num, "ff")
	assert.Equal(t, base, 16)
	assert.Equal(t, bit, 0)

	num, base, bit = getBaseValue("0o77")
	assert.Equal(t, num, "77")
	assert.Equal(t, base, 8)
	assert.Equal(t, bit, 0)

	num, base, bit = getBaseValue("0b11")
	assert.Equal(t, num, "11")
	assert.Equal(t, base, 2)
	assert.Equal(t, bit, 0)


	num, base, bit = getBaseValue("3000")
	assert.Equal(t, num, "3000")
	assert.Equal(t, base, 10)
	assert.Equal(t, bit, 0)
}

func Test_getStringToComplex(t *testing.T) {
	com, err := getStringToComplex("123", "321")
	if assert.NoError(t, err) {
		assert.Equal(t, com, 123+321i)
	}

	com, err = getStringToComplex("123", "")
	if assert.Error(t, err) {
		assert.Equal(t, com, complex128(0))
	}

	com, err = getStringToComplex("", "321")
	if assert.Error(t, err) {
		assert.Equal(t, com, complex128(0))
	}
}