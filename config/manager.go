package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

// GetStartupConfig reads StartupConfig from environment variables
func GetStartupConfig() (*StartupConfig, error) {
	var cfg StartupConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Backup.Enabled {
		el := reflect.TypeOf(cfg.Backup).Elem()
		fieldDirectory, ok := el.FieldByName("Directory")
		if !ok {
			return nil, fmt.Errorf("internal error: field Directory is not found")
		}
		if len(cfg.Backup.Directory) == 0 {
			return nil, fmt.Errorf("required environment variable[%v]", getStructTag(fieldDirectory, "envconfig"))
		}
	}
	return &cfg, nil
}

func getStructTag(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}
