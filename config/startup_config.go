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
		Port    string  `envconfig:"BRUM_SERVER_PORT"`
		SSL     bool    `envconfig:"BRUM_SERVER_SSL" default:"false"`
		SSLType SSLType `envconfig:"BRUM_SERVER_SSL_TYPE" default:"FILE"`
		SSLFile struct {
			SSLFileCertFile string `envconfig:"BRUM_SERVER_SSL_CERT_FILE"`
			SSLFileKeyFile  string `envconfig:"BRUM_SERVER_SSL_KEY_FILE"`
		}
		SSLLetsEncrypt struct {
			Domain string `envconfig:"BRUM_SERVER_SSL_LETS_ENCRYPT_DOMAIN"`
		}
	}
	Subscription struct {
		Enabled bool `envconfig:"BRUM_SUBSCRIPTION_ENABLED" default:"false"`
	}
	PrivateAPI struct {
		Token string `envconfig:"BRUM_PRIVATE_API_TOKEN"`
	}
	Database struct {
		Host         string `required:"true" envconfig:"BRUM_DATABASE_HOST"`
		Port         int16  `required:"true" envconfig:"BRUM_DATABASE_PORT" default:"9000"`
		Username     string `required:"true" envconfig:"BRUM_DATABASE_USERNAME" default:"default"`
		Password     string `required:"true" envconfig:"BRUM_DATABASE_PASSWORD"`
		DatabaseName string `required:"true" envconfig:"BRUM_DATABASE_NAME" default:"default"`
		TablePrefix  string `envconfig:"BRUM_DATABASE_TABLE_PREFIX"`
	}
	Backup struct {
		Enabled          bool   `envconfig:"BRUM_BACKUP_ENABLED" default:"false"`
		Directory        string `envconfig:"BRUM_BACKUP_DIRECTORY"`
		IntervalSeconds  uint32 `envconfig:"BRUM_BACKUP_INTERVAL_SECONDS" default:"5"`
		CompressionType  string `envconfig:"BRUM_COMPRESSION_TYPE" default:"GZIP"`
		CompressionLevel string `envconfig:"BRUM_COMPRESSION_LEVEL"`
	}
}
