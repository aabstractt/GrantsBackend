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
            "message":    "No expired provided",
            "service_id": helper.ServiceId,
        })
    } else if exp != "true" && exp != "false" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":    "Invalid state provided",
            "service_id": helper.ServiceId,
        })
    } else if src := ctx.Query("src"); src == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":    "No source provided",
            "service_id": helper.ServiceId,
        })
    } else if src != "id" && src != "gt" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":    "Invalid source provided",
            "service_id": helper.ServiceId,
        })
    } else if v := ctx.Params("value"); v == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":    "No value provided",
            "service_id": helper.ServiceId,
        })
    } else if t, err := grants.Service().HandleLookup(v, src == "id"); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message":    err.Error(),
            "service_id": helper.ServiceId,
        })
    } else if t == nil {
        return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
            "message":    "Player not found",
            "service_id": helper.ServiceId,
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
