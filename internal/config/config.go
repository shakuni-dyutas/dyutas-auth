// Package config loads overall auth app configurations.
package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Port string

	AllowOrigins  []string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string

	RDbConfig  RDbConfig
	AuthConfig AuthConfig
}

const confKeyPort = "AUTH_APP_PORT"

var allowOrigins = []string{"https://local.dyutas.com:8010", "https://local-api.dyutas.com:8010"}
var allowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
var allowHeaders = []string{"Origin", "Content-Type", "Authorization"}
var exposeHeaders = []string{"Content-Length"}

func LoadAppConfigs() (AppConfig, error) {
	port, ok := LoadConfigOf(confKeyPort)
	if !ok {
		return AppConfig{}, fmt.Errorf("AUTH_APP_PORT is not set")
	}

	unconfigureds := []string{}

	rdbConfig, rdbUnconfigureds := loadRDBConfig()
	authConfig, authUnconfigureds := loadAuthConfig()

	unconfigureds = append(unconfigureds, rdbUnconfigureds...)
	unconfigureds = append(unconfigureds, authUnconfigureds...)

	if len(unconfigureds) > 0 {
		return AppConfig{}, fmt.Errorf("required configurations are not set: %v", unconfigureds)
	}

	return AppConfig{
		Port:          port,
		AllowOrigins:  allowOrigins,
		AllowMethods:  allowMethods,
		AllowHeaders:  allowHeaders,
		ExposeHeaders: exposeHeaders,
		RDbConfig:     rdbConfig,
		AuthConfig:    authConfig,
	}, nil
}

// LoadConfigOf loads the configuration value of given key in basic manner.
// It returns false if the configuration is not set.
func LoadConfigOf(key string) (v string, ok bool) {
	v, ok = os.LookupEnv(key)

	return v, ok
}
