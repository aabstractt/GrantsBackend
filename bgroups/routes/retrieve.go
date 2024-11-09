package routes

import (
    "github.com/Mides-Projects/Kyro/bgroups"
    "github.com/Mides-Projects/Operator/helper"
    "github.com/gofiber/fiber/v3"
)

// Retrieve handles the retrieval of all groups.
func Retrieve(ctx fiber.Ctx) error {
    body := map[string]interface{}{}
    for _, g := range bgroups.Service().Values() {
        body[g.ID()] = g.Marshal()
    }

    if len(body) == 0 {
        return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
            "message":    "No groups found",
            "service_id": helper.ServiceId,
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(body)
}
