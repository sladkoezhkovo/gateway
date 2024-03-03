package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"net/http"
)

type server struct {
	e   *echo.Echo
	cfg *config.HttpConfig
}

func New(cfg *config.HttpConfig) *server {
	return &server{
		e:   echo.New(),
		cfg: cfg,
	}
}

func (s *server) Setup() {
	s.e.Use(echomw.CORS())
	s.e.Use(echomw.Logger())

	s.e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echo.Map{
			"message": "pong",
		})
	})
}

func (s *server) Start() {
	s.e.Logger.Fatal(s.e.Start(fmt.Sprintf(":%d", s.cfg.Port)))
}
