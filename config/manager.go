package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func GetStartupConfig() StartupConfig {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	confPath := path + "/startup_config.yaml"

	f, err := os.Open(confPath)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	var cfg StartupConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		log.Println(err)
	}

	return cfg
}
