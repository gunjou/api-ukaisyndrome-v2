package absensi

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router, handler *Handler) {

	absensi := r.Group("/absensi/peserta/me")

	absensi.Post("/check-in", handler.CheckInPeserta)
	absensi.Get("", handler.GetAbsensiPeserta)
	absensi.Get("/:id_absensi", handler.GetAbsensiPesertaByID)
	absensi.Get("/jadwal/:id_jadwal", handler.GetAbsensiStatusByJadwal)
}