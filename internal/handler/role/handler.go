package role

import (
	"context"
	"github.com/gofiber/fiber/v2"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/handler"
)

type Service interface {
	List(ctx context.Context, limit, offset int32) (*api.ListRoleResponse, error)
	FindById(ctx context.Context, id int64) (*api.Role, error)
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
		if err := ctx.QueryParser(&bounds); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		entries, err := h.service.List(ctx.Context(), bounds.Limit, bounds.Offset)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, entries)
	}
}

func (h *Handler) FindById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "bad id")
		}

		role, err := h.service.FindById(ctx.Context(), int64(id))
		if err != nil {
			// TODO process rpc error
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, role)
	}
}
