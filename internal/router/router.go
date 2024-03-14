package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/handler/auth"
)

type router struct {
	app         *fiber.App
	cfg         *config.Config
	authHandler *auth.Handler
}

func New(cfg *config.Config, authService auth.Service) *router {
	r := &router{
		app:         fiber.New(),
		cfg:         cfg,
		authHandler: auth.New(authService),
	}

	r.app.Use(cors.New())

	api := r.app.Group("/api")
	api.Post("/sign-in", r.authHandler.SignIn())

	return r
}

func (r *router) Start() error {
	return r.app.Listen(fmt.Sprintf(":%d", r.cfg.Http.Port))
}
