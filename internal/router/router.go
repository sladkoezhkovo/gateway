package router

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/handler/auth"
	"github.com/sladkoezhkovo/gateway/internal/handler/city"
	"github.com/sladkoezhkovo/gateway/internal/handler/role"
	"github.com/sladkoezhkovo/gateway/internal/handler/user"
)

const (
	ADMIN         = 1
	MOD           = 2
	SHOP_OWNER    = 3
	FACTORY_OWNER = 4

	AUTHORIZED = 1000
)

type router struct {
	app *fiber.App
	cfg *config.Config

	// Auth service
	authHandler *auth.Handler
	userHandler *user.Handler
	roleHandler *role.Handler

	// Admin service
	cityHandler *city.Handler
}

func New(cfg *config.Config, authService auth.Service, userService user.Service, roleService role.Service, cityService city.Service) *router {
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
		BodyLimit:         10 << 20,
		EnablePrintRoutes: true,
	})

	r := &router{
		app:         app,
		cfg:         cfg,
		authHandler: auth.New(authService),
		userHandler: user.New(userService),
		roleHandler: role.New(roleService),
		cityHandler: city.New(cityService),
	}

	r.app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost:5173, http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	r.app.Use(logger.New())

	api := r.app.Group("/api")
	api.Post("/sign-in", r.authHandler.SignIn())
	api.Post("/sign-up", r.authHandler.Auth(MOD), r.authHandler.SignUp())
	api.Post("/logout", r.authHandler.Logout())
	api.Post("/auth", r.authHandler.CheckAuth())
	api.Get("/refresh", r.authHandler.Refresh())

	users := api.Group("/users", r.authHandler.Auth(ADMIN))
	users.Get("/", r.userHandler.List())
	users.Get("/:id", r.userHandler.FindUserById())
	users.Delete("/:id", r.userHandler.DeleteUser())

	roles := api.Group("/roles")
	roles.Post("/", r.authHandler.Auth(ADMIN), r.roleHandler.Create())
	roles.Get("/", r.roleHandler.List())
	roles.Get("/:id", r.roleHandler.FindById())

	cities := api.Group("/cities")
	cities.Get("/", r.cityHandler.List())
	cities.Get("/:id", r.cityHandler.FindById())

	return r
}

func (r *router) Start() error {
	return r.app.Listen(fmt.Sprintf(":%d", r.cfg.Http.Port))
}
