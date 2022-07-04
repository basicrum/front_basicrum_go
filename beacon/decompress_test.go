package beacon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompress_ScrollLog_Sample1(t *testing.T) {
	res := DecompressBucketLog("00.3bl._.3e..2v..1t.k000D.3x..5s..8n.P.2c..6u..2h.*7*0H.4c..2m.2")

	assert.Equal(
		t,
		[]uint64{0, 4305, 62, 122, 103, 65, 20, 0, 0, 0, 39, 141, 208, 311, 51, 84, 246, 89, 0, 0, 0, 0, 0, 0, 0, 43, 156, 94, 2},
		res)
}

func TestDecompress_ScrollLog_Sample2(t *testing.T) {
	res := DecompressBucketLog("00.43..54..28..3z..29.")

	assert.Equal(t, []uint64{0, 147, 184, 80, 143, 81}, res)
}
