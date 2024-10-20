package wallet

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestString2Cent(t *testing.T) {
	assert.Equal(t, uint(123), String2Cent("1.23"), "they should be equal")
}

func TestCent2String(t *testing.T) {
	assert.Equal(t, "1.23", Cent2String(123), "they should be equal")
}
