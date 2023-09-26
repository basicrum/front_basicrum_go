package backup

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/basicrum/front_basicrum_go/testhelper"
	"github.com/stretchr/testify/require"
)

func Test_archiveDay(t *testing.T) {
	testDay := day(2023, 9, 20)
	tests := []struct {
		name        string
		day         time.Time
		compression Compression
	}{
		{
			name:        "none1",
			day:         testDay,
			compression: NoneCompression,
		},
		{
			name:        "gzip1",
			day:         testDay,
			compression: GZIPCompression,
		},
		{
			name:        "zstd1",
			day:         testDay,
			compression: ZStandardCompression,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			tempDir := copySourceToTempDir(t, tt.name)
			factory := NewCompressionWriterFactory(true, tt.compression, DefaultCompressionLevel)

			// when
			err := archiveDay(tempDir, tt.day, factory)
			require.NoError(t, err)

			// then
			wantDir := path.Join(testhelper.GetProjectRoot(), "testdata", tt.name, "target")
			testhelper.AssertDirEqual(t, wantDir, tempDir)
		})
	}
}

func copySourceToTempDir(t *testing.T, name string) string {
	tempDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	backupRootDir := path.Join(testhelper.GetProjectRoot(), "testdata", name, "source")
	err = testhelper.CopyDir(backupRootDir, tempDir)
	require.NoError(t, err)
	return tempDir
}

func day(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
