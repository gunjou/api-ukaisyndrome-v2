package jadwal

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	jadwal := r.Group("/jadwal")
	jadwal.Get("/", handler.GetJadwalPeserta)
	jadwal.Get("/:id", handler.GetJadwalPesertaByID)
}