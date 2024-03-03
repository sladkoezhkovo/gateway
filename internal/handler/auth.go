package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthService interface {
	//SignIn(email, password string) (*entity.Tokens, error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuth(svc AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) InitRoutes(e *echo.Echo) {
	e.POST("/sign-in", h.signUp())
}

func (h *AuthHandler) signUp() echo.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c echo.Context) error {
		var req request

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, &echo.Map{
				"error": err.Error(),
			})
		}

		//h.svc.

		return c.JSON(http.StatusOK, req)
	}
}
