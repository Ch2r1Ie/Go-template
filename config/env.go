package config

import (
	"os"
	"strings"
)

var Env string

func init() {
	Env = os.Getenv("ENV")
}

const (
	Local string = "LOCAL"
	Dev   string = "DEV"
	UAT   string = "UAT"
	Prod  string = "PROD"
)

func IsLocalEnv() bool {
	return strings.ToUpper(Env) == Local
}

func IsDevEnv() bool {
	return strings.ToUpper(Env) == Dev
}

func IsUATEnv() bool {
	return strings.ToUpper(Env) == UAT
}

func IsProdEnv() bool {
	return strings.ToUpper(Env) == Prod
}
