package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/handler"
	"net/http"
)

type server struct {
	e   *echo.Echo
	cfg *config.HttpConfig

	auth *handler.AuthHandler
}

func New(cfg *config.HttpConfig, auth handler.AuthService) *server {
	return &server{
		e:    echo.New(),
		cfg:  cfg,
		auth: handler.NewAuth(auth),
	}
}

func (s *server) setup() {
	s.e.Use(echomw.CORS())
	s.e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
		Skipper:          echomw.DefaultSkipper,
		Format:           "method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))

	s.e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echo.Map{
			"message": "pong",
		})
	})

	s.auth.InitRoutes(s.e)
}

func (s *server) Start() {
	s.setup()
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%d", s.cfg.Port)))
}
