package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/sladkoezhkovo/gateway/internal/config"
	"github.com/sladkoezhkovo/gateway/internal/handler"
	"github.com/sladkoezhkovo/gateway/pkg/colors"
	"net/http"
)

type server struct {
	e   *echo.Echo
	cfg *config.HttpConfig

	auth *handler.AuthHandler
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func New(cfg *config.HttpConfig, auth handler.AuthService) *server {
	return &server{
		e:    echo.New(),
		cfg:  cfg,
		auth: handler.NewAuth(auth),
	}
}

func (s *server) setup() {
	s.e.Use(echomw.CORS())
	//s.e.Use(echomw.LoggerWithConfig(echomw.LoggerConfig{
	//	Skipper:          echomw.DefaultSkipper,
	//	Format:           "${method}\t${uri}\t| ${status}\ttime=${time_rfc3339}\n",
	//	CustomTimeFormat: "2006-01-02 15:04:05.00000",
	//}))

	// SWITCH LOGGING LOGIC by config.env
	s.e.Use(echomw.RequestLoggerWithConfig(echomw.RequestLoggerConfig{
		LogStatus: true,
		LogMethod: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v echomw.RequestLoggerValues) error {

			methodColor := colors.White
			statusColor := colors.Green

			switch v.Method {
			case GET:
				methodColor = colors.Green
			case POST:
				methodColor = colors.Yellow
			case DELETE:
				methodColor = colors.Red
			case PUT:
				methodColor = colors.Cyan
			}

			if v.Status >= 500 {
				statusColor = colors.Red
			} else if v.Status >= 400 {
				statusColor = colors.Yellow
			} else if v.Status >= 300 {
				statusColor = colors.Cyan
			}

			fmt.Printf("%s%s%s\t", methodColor, v.Method, colors.Reset)
			fmt.Printf("%s\t", v.URI)
			fmt.Printf("| %s%d%s\t", statusColor, v.Status, colors.Reset)
			fmt.Printf("%s\n", v.StartTime)

			return nil
		},
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
