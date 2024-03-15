package router

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/handler/auth"
)

const (
	ADMIN         = 1
	MOD           = 2
	SHOP_OWNER    = 3
	FACTORY_OWNER = 4
)

type router struct {
	app         *fiber.App
	cfg         *config.Config
	authHandler *auth.Handler
}

func New(cfg *config.Config, authService auth.Service) *router {
	app := fiber.New(fiber.Config{
		AppName:       "mail-client-api",
		CaseSensitive: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			err = ctx.Status(code).JSON(fiber.Map{
				"message": e.Message,
			})

			return nil
		},
		BodyLimit: 10 << 20,
	})

	r := &router{
		app:         app,
		cfg:         cfg,
		authHandler: auth.New(authService),
	}

	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	r.app.Use(logger.New())

	api := r.app.Group("/api")
	api.Post("/sign-in", r.authHandler.SignIn())
	api.Post("/sign-up", r.authHandler.Auth(MOD), r.authHandler.SignUp())

	api.Get("/refresh", r.authHandler.Refresh())

	api.Get("/users", r.authHandler.Auth(ADMIN), r.authHandler.List())

	return r
}

func (r *router) Start() error {
	return r.app.Listen(fmt.Sprintf(":%d", r.cfg.Http.Port))
}
