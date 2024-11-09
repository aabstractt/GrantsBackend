package routes

import (
	"github.com/Mides-Projects/Kyro/grants"
	"github.com/Mides-Projects/Operator/helper"
	"github.com/gofiber/fiber/v3"
)

// Lookup handles the lookup of player grants.
func Lookup(ctx fiber.Ctx) error {
	if exp := ctx.Query("expired"); exp == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No expired provided",
		})
	} else if exp != "true" && exp != "false" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid state provided",
		})
	} else if src := ctx.Query("src"); src == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No source provided",
		})
	} else if src != "id" && src != "gt" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid source provided",
		})
	} else if v := ctx.Params("value"); v == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No value provided",
		})
	} else if body, err := grants.Service().HandleLookup(v, src == "id", exp == "true"); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": helper.ServiceId + ": " + err.Error(),
		})
	} else if body == nil {
		return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"message": "No such player found",
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(body)
	}
}
