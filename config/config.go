package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

// Config is a struct that contains all configuration for the application
// NOTE: struct name should be in lowercase and field name should be in uppercase
// you can group the configuration by adding new struct
// Example:
//
//	type Config struct {
//			...
//			GCP gcp  // no need to add tag `env` for struct here.
//	}
//
// then create gcp struct with tag `env` for each field
//
//	type gcp struct {
//		ProjectID string `env:"GCP_PROJECT_ID"`
//	}
//
// you can add field without grouping them by adding new field with tag `env`
// Example:
//
//	type Config struct {
//		...
//		AppName string `env:"APP_NAME"`
//	}
type Config struct {
	Server                    Server
	AccessControl             AccessControl
	Database                  Database
	Header                    Header
	ServiceCoreDltAccountUrl  string `env:"SERVICE_CORE_DLT_ACCOUNT_URL"`
	ServiceCoreDltInterPermit string `env:"SERVICE_CORE_DLT_INTER_PERMIT"`
	// add more configuration here below
}

type Server struct {
	Hostname string `env:"HOSTNAME"`
	Port     string `env:"PORT,notEmpty"`
}

type AccessControl struct {
	AllowOrigin string `env:"ACCESS_CONTROL_ALLOW_ORIGIN"`
}

type Database struct {
	MongoURL    string `env:"MONGO_URL"`    // example: mongodb://localhost:27017
	PostgresURL string `env:"POSTGRES_URL"` // example: postgres://postgres:password@localhost:5432/dbname?sslmode=disable
}

type Header struct {
	RefIDHeaderKey string `env:"REF_ID_HEADER_KEY,notEmpty"`
}

var once sync.Once
var config Config

func prefix(e string) string {
	if e == "" {
		return ""
	}

	return fmt.Sprintf("%s_", e)
}

func C(envPrefix string) Config {
	once.Do(func() {
		opts := env.Options{
			Prefix: prefix(envPrefix),
		}

		var err error
		config, err = parseEnv[Config](opts)
		if err != nil {
			log.Fatal(err)
		}
	})

	return config
}

// TODO: read config from xxx.yaml file that contains ${ENV} variable e.g. serviceDLTUrl: ${SERVICE_CORE_DLT_ACCOUNT_URL}
