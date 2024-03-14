package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/entity"
	"github.com/sladkoezhkovo/gateway/internal/handler"
)

type Service interface {
	SignUp(ctx context.Context, user *entity.User) (*api.TokenResponse, error)
	SignIn(ctx context.Context, user *entity.User) (*api.TokenResponse, error)
	Refresh(ctx context.Context, refresh string) (*api.TokenResponse, error)
	Auth(ctx context.Context, access string, roleId int64) (bool, error)
	Logout(ctx context.Context, access string) error

	List(ctx context.Context, limit, offset int32) (*api.ListUserResponse, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SignIn() fiber.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(ctx *fiber.Ctx) error {
		var req request

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tokens, err := h.service.SignIn(ctx.Context(), &entity.User{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, tokens)
	}
}
