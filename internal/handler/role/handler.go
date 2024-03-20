package role

import (
	"context"
	"github.com/gofiber/fiber/v2"
	api "github.com/sladkoezhkovo/gateway/api/auth"
	"github.com/sladkoezhkovo/gateway/internal/handler"
)

type Service interface {
	Create(ctx context.Context, req *api.CreateRoleRequest) (*api.Role, error)
	List(ctx context.Context, limit, offset int32) (*api.ListRoleResponse, error)
	FindById(ctx context.Context, id int64) (*api.Role, error)

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

func (h *Handler) Create() fiber.Handler {
	type request struct {
		Name      string `json:"name"`
		Authority int32  `json:"authority"`
	}

	return func(ctx *fiber.Ctx) error {
		var dto request

		if err := ctx.BodyParser(&dto); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "bad request")
		}

		req := &api.CreateRoleRequest{
			Name:      dto.Name,
			Authority: dto.Authority,
		}

		role, err := h.service.Create(ctx.Context(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return handler.Respond(ctx, role)
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

func (h *Handler) Delete() fiber.Handler {
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
