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
	"time"
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

	type response struct {
		AccessToken string `json:"accessToken"`
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
			// TODO add process for invalid creds error
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		cookie := &fiber.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
			HTTPOnly: true,
			//Domain:   "http://127.0.0.1:8000",
			//Path:     "/api/sign-in",
			SameSite: "None",
		}

		ctx.Cookie(cookie)

		response := &response{
			AccessToken: tokens.AccessToken,
		}

		return handler.Respond(ctx, response)
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
			return fiber.NewError(fiber.StatusBadRequest, "doesnt have authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid token")
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

func (h *Handler) CheckAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		roleId := int64(ctx.QueryInt("role_id", -1))
		if roleId == -1 {
			return fiber.NewError(fiber.StatusBadRequest, "enter a role id")
		}

		header := ctx.Get("Authorization", "")
		if header == "" {
			return fiber.NewError(fiber.StatusBadRequest, "no Authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid token")
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

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (h *Handler) Refresh() fiber.Handler {
	type response struct {
		AccessToken string `json:"accessToken"`
	}

	return func(ctx *fiber.Ctx) error {
		refresh := ctx.Cookies("refresh_token")
		if refresh == "" {
			return fiber.NewError(fiber.StatusBadRequest, "empty refresh token")
		}

		tokens, err := h.service.Refresh(ctx.Context(), refresh)
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Unauthenticated:
					return fiber.NewError(fiber.StatusUnauthorized, e.Message())
				}
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		cookie := &fiber.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
			HTTPOnly: true,
			//Domain:   "http://127.0.0.1:8000",
			//Path:     "/api/refresh",
			SameSite: "None",
		}

		ctx.Cookie(cookie)

		response := &response{
			AccessToken: tokens.AccessToken,
		}

		return handler.Respond(ctx, response)
	}
}

func (h *Handler) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization", "")
		if header == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "doesnt have authorization header")
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		access := parts[1]

		if err := h.service.Logout(ctx.Context(), access); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return ctx.SendStatus(200)
	}
}
