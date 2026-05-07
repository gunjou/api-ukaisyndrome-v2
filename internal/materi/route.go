package materi

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	materi := r.Group("/materi")

	materi.Get("/peserta/private", handler.GetMateriPrivatePeserta)
	materi.Get("/peserta/:id_modul", handler.GetMateriPeserta)
}