package main

import (
	"flag"
	"fmt"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/router"
	"github.com/sladkoezhkovo/gateway/internal/service/auth"
	"github.com/sladkoezhkovo/lib"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/.yml", "set path to config")
}

func main() {
	flag.Parse()

	var cfg config.Config
	if err := lib.SetupConfig(configPath, &cfg); err != nil {
		panic(fmt.Errorf("SetupConfig: %s", err))
	}

	authGRPC, err := auth.New(cfg.Auth)
	if err != nil {
		panic(err)
	}

	gateway := router.New(
		&cfg,
		authGRPC,
		authGRPC,
	)
	if err := gateway.Start(); err != nil {
		panic(err)
	}
}
