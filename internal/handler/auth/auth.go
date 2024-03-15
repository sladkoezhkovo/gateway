package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/entity"
	"github.com/sladkoezhkovo/gateway/internal/handler"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Service interface {
	SignUp(ctx context.Context, user *entity.User) (*api.TokenResponse, error)
	SignIn(ctx context.Context, user *entity.User) (*api.TokenResponse, error)
	Refresh(ctx context.Context, refresh string) (*api.TokenResponse, error)
	Auth(ctx context.Context, access string, roleId int64) (bool, error)
	Logout(ctx context.Context, access string) error

	List(ctx context.Context, limit, offset int32) (*api.ListUserResponse, error)
	FindById(ctx context.Context, id int64) (*api.UserDetails, error)
}

var (
	ACCESS_TOKEN = "accessToken"
)

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

func (h *Handler) SignUp() fiber.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		RoleId   int64  `json:"roleId"`
	}

	return func(ctx *fiber.Ctx) error {
		var req request
		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if req.RoleId <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "role id must be set & be greater than zero")
		}

		tokens, err := h.service.SignUp(ctx.Context(), &entity.User{
			Email:    req.Email,
			Password: req.Password,
			Role:     entity.Role{Id: req.RoleId},
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, tokens)
	}
}

func (h *Handler) Auth(roleId int64) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		header := ctx.Get("Authorization", "")
		if header == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "doesnt have authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		approve, err := h.service.Auth(ctx.Context(), parts[1], roleId)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unauthenticated:
					return fiber.NewError(fiber.StatusUnauthorized, e.Message())
				}
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if !approve {
			return fiber.NewError(fiber.StatusForbidden, "insufficient permission")
		}

		ctx.Set(ACCESS_TOKEN, parts[1])

		return ctx.Next()
	}
}

func (h *Handler) Refresh() fiber.Handler {
	type request struct {
		RefreshToken string `json:"refreshToken"`
	}
	return func(ctx *fiber.Ctx) error {
		var req request
		if err := ctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		tokens, err := h.service.Refresh(ctx.Context(), req.RefreshToken)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, tokens)
	}
}

func (h *Handler) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		access := ctx.Get(ACCESS_TOKEN)

		if err := h.service.Logout(ctx.Context(), access); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return ctx.SendStatus(200)
	}
}

func (h *Handler) List() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var bounds api.Bounds
		if err := ctx.BodyParser(&bounds); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		entries, err := h.service.List(ctx.Context(), bounds.Limit, bounds.Offset)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, entries)
	}
}

func (h *Handler) FindUserById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "bad id")
		}

		user, err := h.service.FindById(ctx.Context(), int64(id))
		if err != nil {
			// TODO process rpc error
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, user)
	}
}
