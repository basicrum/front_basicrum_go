package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func archiveDay(backupRootDir string, day time.Time, factory CompressionWriterFactory) error {
	files, err := os.ReadDir(backupRootDir)
	if err != nil {
		return fmt.Errorf("cannot read dir[%v] err[%w]", backupRootDir, err)
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		if err := archiveHost(backupRootDir, file.Name(), day, factory); err != nil {
			return err
		}
	}
	return nil
}

func archiveHost(backupRootDir, host string, day time.Time, factory CompressionWriterFactory) error {
	datePath := makeDayPath(backupRootDir, host, day)
	err := validateDirExist(datePath)
	if err != nil {
		return err
	}

	daySummary, err := collectHours(datePath, day)
	if err != nil {
		return err
	}

	err = writeDayFile(backupRootDir, host, day, daySummary.dayContent, factory)
	if err != nil {
		return err
	}

	err = writeHourlySummary(backupRootDir, host, day, daySummary.linesPerHour)
	if err != nil {
		return err
	}

	return os.RemoveAll(datePath)
}

type daySummary struct {
	dayContent   string
	linesPerHour map[int]int
}

func collectHours(datePath string, day time.Time) (*daySummary, error) {
	var dayContent string
	linesPerHour := map[int]int{}
	for hour := 0; hour < 24; hour++ {
		hourPath := makeHourPath(datePath, dayWithHour(day, hour))
		fileContent, err := readAll(hourPath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error read file[%v] err[%w]", hourPath, err)
		}
		lines := string(fileContent)
		dayContent += lines
		linesPerHour[hour] = len(strings.Split(lines, "\n")) - 1
	}
	return &daySummary{
		dayContent:   dayContent,
		linesPerHour: linesPerHour,
	}, nil
}

func validateDirExist(datePath string) error {
	dateDir, err := os.Stat(datePath)
	if err != nil {
		if os.IsNotExist(err) {
			// nothing to do
			return nil
		}
		return err
	}
	if !dateDir.IsDir() {
		return fmt.Errorf("expected directory[%v]", datePath)
	}
	return nil
}

// nolint: revive
func writeDayFile(backupRootDir string, host string, day time.Time, dayContent string, factory CompressionWriterFactory) error {
	archiveDayPath := makeArchiveDayPath(backupRootDir, host, day)
	if err := writeToFile(archiveDayPath, dayContent, factory); err != nil {
		return fmt.Errorf("cannot write to file[%v] err[%w]", archiveDayPath, err)
	}
	return nil
}

func writeHourlySummary(backupRootDir string, host string, day time.Time, linesPerHour any) error {
	hourMeta, err := json.Marshal(linesPerHour)
	if err != nil {
		return err
	}

	archiveDayMetaPath := makeArchiveDayMetaPath(backupRootDir, host, day)
	if err := writeToFile(archiveDayMetaPath, string(hourMeta), newNoneFactory()); err != nil {
		return fmt.Errorf("cannot write to file[%v] err[%w]", archiveDayMetaPath, err)
	}
	return nil
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

func dayWithHour(day time.Time, hour int) time.Time {
	return day.Truncate(time.Hour).Add(time.Hour * time.Duration(hour))
}
