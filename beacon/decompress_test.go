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

func TestDecompress_ScrollLog_Sample3(t *testing.T) {
	res := DecompressBucketLog("0*j*0td")

	assert.Equal(t, []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 29, 13}, res)
}

func TestDecompress_ScrollLog_Sample4(t *testing.T) {
	res := DecompressBucketLog("000.40..28..4x..54..5b..45.m*8*0.6c..27..3h..b7..e8.k*9*0O.5f..3u..3s..2l.lS.5v..3a..38..3q..3a..24.")

	assert.Equal(t, []uint64{0, 0, 144, 80, 177, 184, 191, 149, 22, 0, 0, 0, 0, 0, 0, 0, 0, 228, 79, 125, 403, 512, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 50, 195, 138, 136, 93, 21, 54, 211, 118, 116, 134, 118, 76}, res)
}
