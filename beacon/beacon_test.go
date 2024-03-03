package beacon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMobDlFloatToInt(t *testing.T) {
	assert.Equal(
		t,
		"1",
		RoundFloatParam("1.3"))

	assert.Equal(
		t,
		"2",
		RoundFloatParam("1.7"))

	assert.Equal(
		t,
		"",
		RoundFloatParam(""))

	assert.Equal(
		t,
		"",
		RoundFloatParam("bad val"))
}
