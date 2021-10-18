package config

import (
	"gopkg.in/gcfg.v1"
	"log"
	"time"
)

type MainConfig struct {
	UserHTTP struct {
		Port       string
		PathPrefix string
	}

	AuthHTTP struct {
		Port       string
		PathPrefix string
	}

	Mongo struct {
		URL          string
		UserDatabase string
		AuthDatabase string
	}

	Redis struct {
		Connection string
		Password   string
		DB         int
		Expiration time.Duration
	}

	GrpcClient struct {
		GRPCUserURL string
		GRPCAuthURL string
	}

	GrpcServer struct {
		AuthPort string
		UserPort string
	}
}

func ReadConfig(cfg interface{}) interface{} {
	ok := ReadModuleConfig(cfg, ".")
	if !ok {
		log.Fatalln("failed to read config.")
	}
	return cfg
}

func ReadModuleConfig(cfg interface{}, path string) bool {
	fname := path + "/" + ".env"
	err := gcfg.ReadFileInto(cfg, fname)
	if err == nil {
		return true
	}

	return false
}
