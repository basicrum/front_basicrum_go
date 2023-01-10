package config

// StartupConfig contains application configuration
type StartupConfig struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"database_name"`
		Host         string `yaml:"host"`
		Port         int16  `yaml:"port"`
		TablePrefix  string `yaml:"table_prefix"`
	} `yaml:"database"`
	Persistance struct {
		DatabaseStrategy string `yaml:"database_strategy"`
		TableStrategy    string `yaml:"table_strategy"`
	} `yaml:"persistance"`
	Backup struct {
		Enabled         bool   `yaml:"enabled"`
		Directory       string `yaml:"directory"`
		IntervalSeconds uint32 `yaml:"interval_seconds"`
	} `yaml:"backup"`
}
