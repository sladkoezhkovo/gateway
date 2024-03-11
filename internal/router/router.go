package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sladkoezhkovo/gateway/internal/config"
)

type router struct {
	app *fiber.App
	cfg *config.Config
}

func New(cfg *config.Config) *router {
	return &router{
		app: fiber.New(),
		cfg: cfg,
	}
}

func (r *router) Start() error {
	return r.app.Listen(fmt.Sprintf(":%d", r.cfg.Http.Port))
}
