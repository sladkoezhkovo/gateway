package main

import (
	"flag"
	"fmt"
	admingrpc "github.com/sladkoezhkovo/gateway/api/admin"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/router"
	"github.com/sladkoezhkovo/gateway/internal/service/admin"
	"github.com/sladkoezhkovo/gateway/internal/service/auth"
	"github.com/sladkoezhkovo/lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	authConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Auth.Host, cfg.Auth.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	authClientGrpc := api.NewAuthServiceClient(authConn)

	adminConn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Admin.Host, cfg.Auth.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	adminClientGrpc := admingrpc.NewAdminServiceClient(adminConn)

	userService, err := auth.NewUserService(authClientGrpc)
	if err != nil {
		panic(err)
	}
	roleService, err := auth.NewRoleService(authClientGrpc)
	if err != nil {
		panic(err)
	}

	cityService, err := admin.NewCityService(adminClientGrpc)

	gateway := router.New(
		&cfg,
		userService,
		userService,
		roleService,
		cityService,
	)
	if err := gateway.Start(); err != nil {
		panic(err)
	}
}
