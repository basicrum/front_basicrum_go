package it

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func processFile(file string) ([]url.Values, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file[%v]: %w", file, err)
	}

	var rows []map[string]any
	err = json.Unmarshal(content, &rows)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json[%s]: %w", content, err)
	}
	var result []url.Values
	for i, row := range rows {
		beaconData, ok := row["beacon_data"].(string)
		if !ok {
			return nil, fmt.Errorf("unable to parse beacon data[%s], file: %s, line: %d, %w", row["beacon_data"], file, i+1, err)
		}
		var beaconDataMap map[string]string
		err = json.Unmarshal([]byte(beaconData), &beaconDataMap)

		if err != nil {
			log.Printf("Bad beacon_data in file: %v", err)
			continue
		}

		result = append(result, makeUrlValues(beaconDataMap))
	}
	return result, nil
}

func scanFile(fileName string) ([]url.Values, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to read file[%v]: %w", fileName, err)
	}

	defer file.Close()
	var row map[string]string
	scanner := bufio.NewScanner(file)

	var result []url.Values
	for scanner.Scan() {
		err = json.Unmarshal([]byte(scanner.Text()), &row)
		if err != nil {
			return nil, fmt.Errorf("unable to parse json[%s]: %w", []byte(scanner.Text()), err)
		}
		beaconData := row
		result = append(result, makeUrlValues(beaconData))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan file[%v]: %w", fileName, err)
	}
	return result, nil
}

func makeUrlValues(beaconDataMap map[string]string) url.Values {
	result := url.Values{}
	for k, v := range beaconDataMap {
		result.Set(mapKey(k), v)
	}
	return result
}

func mapKey(value string) string {
	keyPrefix := extractBeaconPrefix(value)

	if keyPrefix == "nt_" {
		return value
	}

	switch value {
	case "created_at":
		return value
	case "t_resp":
		return value
	case "t_done":
		return value
	case "t_page":
		return value
	case "t_other":
		return value
	default:
		return strings.ReplaceAll(value, "_", ".")
	}
}

func extractBeaconPrefix(value string) string {
	bKey := value
	if len(bKey) > 3 {
		return value[0:3]
	}

	return value
}
