package user

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/handler"
)

type Service interface {
	List(ctx context.Context, limit, offset int32) (*api.ListUserResponse, error)
	ListByRole(ctx context.Context, roleId int64, limit, offset int32) (*api.ListUserResponse, error)
	FindById(ctx context.Context, id int64) (*api.UserDetails, error)

	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List() fiber.Handler {

	type params struct {
		Limit  int32 `query:"limit"`
		Offset int32 `query:"offset"`
	}

	return func(ctx *fiber.Ctx) error {

		var bounds params
		var err error

		if err := ctx.QueryParser(&bounds); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		roleId := ctx.QueryInt("roleId", -1)
		email := ctx.Query("email", "")

		fmt.Println("roleId = ", roleId)
		fmt.Println("email = ", email)

		var entries *api.ListUserResponse

		if roleId != -1 {
			entries, err = h.service.ListByRole(ctx.Context(), int64(roleId), bounds.Limit, bounds.Offset)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		} else {
			entries, err = h.service.List(ctx.Context(), bounds.Limit, bounds.Offset)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
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

func (h *Handler) DeleteUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "bad id")
		}

		if err := h.service.Delete(ctx.Context(), int64(id)); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
