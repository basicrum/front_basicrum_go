package config

// StartupConfig contains application configuration
type StartupConfig struct {
	Server struct {
		Host string `required:"true" envconfig:"SERVER_HOST"`
		Port string `required:"true" envconfig:"SERVER_PORT"`
	}
	Database struct {
		Username     string `required:"true" envconfig:"DATABASE_USERNAME"`
		Password     string `required:"true" envconfig:"DATABASE_PASSWORD"`
		DatabaseName string `required:"true" envconfig:"DATABASE_NAME"`
		Host         string `required:"true" envconfig:"DATABASE_HOST"`
		Port         int16  `required:"true" envconfig:"DATABASE_PORT"`
		TablePrefix  string `envconfig:"DATABASE_TABLE_PREFIX"`
	}
	Persistance struct {
		DatabaseStrategy string `envconfig:"PERSISTANCE_DATABASE_STRATEGY"`
		TableStrategy    string `envconfig:"PERSISTANCE_TABLE_STRATEGY"`
	}
	Backup struct {
		Enabled         bool   `envconfig:"BACKUP_ENABLED"`
		Directory       string `envconfig:"BACKUP_DIRECTORY"`
		IntervalSeconds uint32 `envconfig:"BACKUP_INTERVAL_SECONDS"`
	}
}
