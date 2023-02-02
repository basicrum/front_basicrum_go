package config

// SSLType is the type of http SSL configuration to use
type SSLType string

const (
	// SSLTypeFile file configuration
	SSLTypeFile SSLType = "FILE"
	// SSLTypeLetsEncrypt let's encrypt configuration
	SSLTypeLetsEncrypt SSLType = "LETS_ENCRYPT"
)

// StartupConfig contains application configuration
type StartupConfig struct {
	Server struct {
		Host    string  `required:"true" envconfig:"SERVER_HOST"`
		Port    string  `required:"true" envconfig:"SERVER_PORT"`
		SSL     bool    `envconfig:"SERVER_SSL" default:"false"`
		SSLType SSLType `envconfig:"SERVER_SSL_TYPE" default:"FILE"`
		SSLFile struct {
			SSLFileCertFile string `envconfig:"SERVER_SSL_CERT_FILE"`
			SSLFileKeyFile  string `envconfig:"SERVER_SSL_KEY_FILE"`
		}
		SSLLetsEncrypt struct {
			Port   string `envconfig:"SERVER_SSL_LETS_ENCRYPT_PORT" default:"80"`
			Domain string `envconfig:"SERVER_SSL_LETS_ENCRYPT_DOMAIN"`
		}
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
