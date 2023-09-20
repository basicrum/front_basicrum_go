package testhelper

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// AssertDirEqual compares two directories if they are the same recursively
func AssertDirEqual(
	t *testing.T,
	expected, actual string,
) {
	doAssertDirEqual(t, expected, actual, "")
}

func doAssertDirEqual(
	t *testing.T,
	expected, actual string,
	relativePath string,

) {
	sourceFiles := assertNamesAreTheSame(t, true, expected, actual, relativePath)
	sourceDirs := assertNamesAreTheSame(t, false, expected, actual, relativePath)

	for fileName := range sourceFiles {
		assertFileContentIsTheSame(t, expected, actual, filepath.Join(relativePath, fileName))
	}

	for dirName := range sourceDirs {
		doAssertDirEqual(t, expected, actual, filepath.Join(relativePath, dirName))
	}
}

// nolint: revive
func assertNamesAreTheSame(
	t *testing.T,
	collectFiles bool,
	expected, actual string,
	relativePath string,
) map[string]bool {
	expectedNames := requireCollectNameMap(t, collectFiles, expected, relativePath)
	actualNames := requireCollectNameMap(t, collectFiles, actual, relativePath)
	require.Equal(t, expectedNames, actualNames, "collectFiles[%v] relativePath[%v]", collectFiles, relativePath)
	return actualNames
}

func requireCollectNameMap(
	t *testing.T,
	collectFiles bool,
	parent string,
	relativePath string,
) map[string]bool {
	dirPath := filepath.Join(parent, relativePath)
	expectedNames, err := collectNameMap(collectFiles, dirPath)
	require.NoError(t, err, "dirPath[%v]", dirPath)
	return expectedNames
}

func assertFileContentIsTheSame(
	t *testing.T,
	expected, actual string,
	relativePath string,
) {
	expectedBytes := requireReadAll(t, expected, relativePath)
	actualBytes := requireReadAll(t, actual, relativePath)
	require.Equal(t, string(expectedBytes), string(actualBytes), "relativePath[%v]", relativePath)
}

func requireReadAll(t *testing.T, parent string, relativePath string) []byte {
	filePath := filepath.Join(parent, relativePath)
	expectedBytes, err := readAll(filePath)
	require.NoError(t, err, "filePath[%v]", filePath)
	return expectedBytes
}

func readAll(
	filePath string,
) ([]byte, error) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()
	return io.ReadAll(sourceFile)
}

// nolint: revive
func collectNameMap(
	collectFiles bool,
	dirPath string,
) (map[string]bool, error) {
	result := map[string]bool{}
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range files {
		if skipFile(fileInfo) {
			continue
		}
		if collectFiles == isFile(fileInfo) {
			result[fileInfo.Name()] = true
		}
	}
	return result, nil
}

func skipFile(fileInfo fs.DirEntry) bool {
	if fileInfo.Name() == ".gitkeep" {
		return true
	}
	if fileInfo.Name() == ".git" {
		return true
	}
	return false
}

func isFile(fileInfo fs.DirEntry) bool {
	return !fileInfo.IsDir()
}
