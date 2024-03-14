package handler

import "github.com/gofiber/fiber/v2"

func Respond(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(&fiber.Map{
		"data": data,
	})
}
