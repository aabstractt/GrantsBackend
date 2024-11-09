package routes

import (
    "github.com/Mides-Projects/Kyro/bgroups"
    "github.com/Mides-Projects/Operator/helper"
    "github.com/gofiber/fiber/v3"
)

func Create(ctx fiber.Ctx) error {
    if name := ctx.Params("name"); name == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":    "No name provided",
            "service_id": helper.ServiceId,
        })
    } else if g := bgroups.Service().LookupByName(name); g != nil {
        return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
            "message":    "Group with name '" + name + "' already exists",
            "service_id": helper.ServiceId,
        })
    } else if id, err := bgroups.Service().Insert(name); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message":    err.Error(),
            "service_id": helper.ServiceId,
        })
    } else if id == "" {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message":    "Failed to insert group",
            "service_id": helper.ServiceId,
        })
    } else {
        return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
            "id": id,
        })
    }
}
