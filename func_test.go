package def

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestSetFunc(t *testing.T) {
	assert.NoError(t, SetFunc("key", func(self *sample) interface{} {
		return func() {}
	}))
	assert.Error(t, SetFunc("key", func(self *sample) interface{} {
		return func() {}
	}))
	assert.Error(t, SetFunc("error1", func(self interface{}) interface{} {
		return func() {}
	}))
	assert.Error(t, SetFunc("error2", func(self1 *sample, self2 *sample) interface{} {
		return func() {}
	}))
}
