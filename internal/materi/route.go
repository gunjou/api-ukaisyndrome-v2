package materi

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	materi := r.Group("/materi")

	materi.Get("/peserta/private", handler.GetMateriPrivatePeserta)
	materi.Post("/peserta/:id_materi/progress", handler.OpenMateri)
	materi.Get("/peserta/progress", handler.GetMateriProgress)
	materi.Get("/peserta/progress/monitoring", handler.GetMateriProgressMonitoring)
	materi.Get("/peserta/:id_modul", handler.GetMateriPeserta)
}