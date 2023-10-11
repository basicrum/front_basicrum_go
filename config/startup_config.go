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
			Port   string `envconfig:"BRUM_SERVER_SSL_LETS_ENCRYPT_PORT" default:"80"`
			Domain string `envconfig:"BRUM_SERVER_SSL_LETS_ENCRYPT_DOMAIN"`
		}
	}
	PrivateAPI struct {
		Token string `envconfig:"BRUM_PRIVATE_API_TOKEN"`
	}
	Database struct {
		Username     string `required:"true" envconfig:"BRUM_DATABASE_USERNAME"`
		Password     string `required:"true" envconfig:"BRUM_DATABASE_PASSWORD"`
		DatabaseName string `required:"true" envconfig:"BRUM_DATABASE_NAME"`
		Host         string `required:"true" envconfig:"BRUM_DATABASE_HOST"`
		Port         int16  `required:"true" envconfig:"BRUM_DATABASE_PORT"`
		TablePrefix  string `envconfig:"BRUM_DATABASE_TABLE_PREFIX"`
	}
	Persistance struct {
		DatabaseStrategy string `envconfig:"BRUM_PERSISTANCE_DATABASE_STRATEGY"`
		TableStrategy    string `envconfig:"BRUM_PERSISTANCE_TABLE_STRATEGY"`
	}
	Backup struct {
		Enabled          bool   `envconfig:"BRUM_BACKUP_ENABLED"`
		Directory        string `envconfig:"BRUM_BACKUP_DIRECTORY"`
		IntervalSeconds  uint32 `envconfig:"BRUM_BACKUP_INTERVAL_SECONDS"`
		CompressionType  string `envconfig:"BRUM_COMPRESSION_TYPE" default:"GZIP"`
		CompressionLevel string `envconfig:"BRUM_COMPRESSION_LEVEL"`
	}
}
