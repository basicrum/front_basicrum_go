package beacon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompress_ScrollLog(t *testing.T) {
	res := DecompressBucketLog("00.43..54..28..3z..29.")

	assert.Equal(t, res, []uint64{0, 147, 184, 80, 143, 81})
}
