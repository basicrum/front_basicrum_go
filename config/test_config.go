package config

import (
	"fmt"
	"os"
)

// SetTestDefaultConfig sets default test configuration
func SetTestDefaultConfig() {
	vars := map[string]string{
		"BRUM_SERVER_HOST":                   "localhost",
		"BRUM_SERVER_PORT":                   "8087",
		"BRUM_DATABASE_HOST":                 "localhost",
		"BRUM_DATABASE_PORT":                 "9000",
		"BRUM_DATABASE_NAME":                 "default",
		"BRUM_DATABASE_USERNAME":             "default",
		"BRUM_DATABASE_PASSWORD":             "",
		"BRUM_DATABASE_TABLE_PREFIX":         "local_test_",
		"BRUM_PERSISTANCE_DATABASE_STRATEGY": "all_in_one_db",
		"BRUM_PERSISTANCE_TABLE_STRATEGY":    "all_in_one_table",
		"BRUM_BACKUP_ENABLED":                "false",
		"BRUM_BACKUP_DIRECTORY":              "/home/basicrum_backup/archive",
		"BRUM_BACKUP_EXPIRED_DIRECTORY":      "/home/basicrum_backup/expired",
		"BRUM_BACKUP_UNKNOWN_DIRECTORY":      "/home/basicrum_backup/unknown",
		"BRUM_BACKUP_INTERVAL_SECONDS":       "5",
	}
	for k, v := range vars {
		setIfEmpty(k, v)
	}
}

func setIfEmpty(key, defaultValue string) {
	value := os.Getenv(key)
	if value != "" {
		return
	}
	err := os.Setenv(key, defaultValue)
	if err != nil {
		panic(fmt.Errorf("cannot set env var[%v]=[%v] err: %w", key, value, err))
	}
}
