package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

// GetStartupConfig reads StartupConfig from environment variables
// nolint: revive
func GetStartupConfig() (*StartupConfig, error) {
	var cfg StartupConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Backup.Enabled {
		if len(cfg.Backup.Directory) == 0 {
			return nil, fieldError(cfg.Backup, "Directory")
		}
	}
	if cfg.Server.SSL {
		switch cfg.Server.SSLType {
		case SSLTypeFile:
			if len(cfg.Server.SSLFile.SSLFileCertFile) == 0 {
				return nil, fieldError(cfg.Server.SSLFile, "SSLFileCertFile")
			}
			if len(cfg.Server.SSLFile.SSLFileKeyFile) == 0 {
				return nil, fieldError(cfg.Server.SSLFile, "SSLFileKeyFile")
			}
		case SSLTypeLetsEncrypt:
			if len(cfg.Server.SSLLetsEncrypt.Domain) == 0 {
				return nil, fieldError(cfg.Server.SSLLetsEncrypt, "Domain")
			}
		default:
			return nil, fmt.Errorf("unsupported ssl type[%v]", cfg.Server.SSLType)
		}
	}
	return &cfg, nil
}

func getStructTag(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}

func fieldError(parentElement any, fieldName string) error {
	el := reflect.TypeOf(parentElement).Elem()
	field, ok := el.FieldByName(fieldName)
	if !ok {
		return fmt.Errorf("internal error: field %v is not found", fieldName)
	}
	return fmt.Errorf("required environment variable[%v]", getStructTag(field, "envconfig"))
}
