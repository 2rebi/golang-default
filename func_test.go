package def

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetFunc(t *testing.T) {
	assert.NoError(t, SetFunc("key", func(self interface{}) interface{} {
		return func() {}
	}))
	assert.Error(t, SetFunc("key", func(self interface{}) interface{} {
		return func() {}
	}))
}
