package module

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {
	modul := r.Group("/modul")

	modul.Get("/peserta", handler.GetModul)
}