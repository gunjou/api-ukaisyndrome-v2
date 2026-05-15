package cdn

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	cdn := r.Group("/cdn")

	cdn.Get("/mentor", handler.GetMentorImages)
	cdn.Get("/modul", handler.GetModulImages)
	cdn.Get("/news", handler.GetNews)
}