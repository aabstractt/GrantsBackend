package routes

import (
    "github.com/Mides-Projects/kyro"
    "github.com/gofiber/fiber/v3"
)

// Lookup handles the player lookup.
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
    } else if t, err := kyro.Service().HandleLookup(v, src == "id"); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    } else if t == nil {
        return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
            "message": "Player not found",
        })
    } else {
        expired := make(map[string]interface{})
        if exp == "true" {
            for _, gi := range t.Expired() {
                expired[gi.ID()] = gi.Marshal()
            }
        }

        actives := make(map[string]interface{})
        for _, gi := range t.Actives() {
            actives[gi.ID()] = gi.Marshal()
        }

        return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
            "expired": expired,
            "actives": actives,
        })
    }
}
