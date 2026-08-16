package absensi

import (
	"errors"

	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}


// CheckInPeserta godoc
// @Summary Check-in peserta
// @Description Peserta melakukan check-in pada jadwal ketika mentor sedang mengajar dan peserta berada dalam radius lokasi mentor
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckInPesertaRequest true "Data check-in peserta"
// @Router /absensi/peserta/me/check-in [post]
func (h *Handler) CheckInPeserta(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	var req CheckInPesertaRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "invalid request body", "BAD_REQUEST", nil)
	}

	if req.IDJadwal <= 0 {
		return response.Error(c, 400, "id_jadwal is required", "BAD_REQUEST", nil)
	}

	data, err := h.Service.CheckInPeserta(
		c.Context(),
		userID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrJadwalNotFound):
			return response.Error(c, 404, err.Error(), "NOT_FOUND", nil)

		case errors.Is(err, ErrMentorNotStarted):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrMeetingFinished):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrCheckInExpired):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrLocationRequired):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrLocationUnavailable):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrOutsideRadius):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrAlreadyCheckedIn):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
		}

		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiStatusByJadwal godoc
// @Summary Get status absensi peserta berdasarkan jadwal
// @Description Mengecek apakah peserta yang sedang login sudah melakukan check-in pada jadwal tertentu
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_jadwal path int true "ID Jadwal"
// @Router /absensi/peserta/me/jadwal/{id_jadwal} [get]
func (h *Handler) GetAbsensiStatusByJadwal(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	jadwalID, err := c.ParamsInt("id_jadwal")
	if err != nil {
		return response.Error(c, 400, "invalid jadwal id", "BAD_REQUEST", nil)
	}

	data, err := h.Service.GetAbsensiStatusByJadwal(c.Context(), userID, jadwalID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	if data == nil {
		return response.Error(c, 404, "jadwal tidak ditemukan", "NOT_FOUND", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiPeserta godoc
// @Summary Get absensi peserta
// @Description Get seluruh riwayat absensi peserta yang sedang login
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /absensi/peserta/me [get]
func (h *Handler) GetAbsensiPeserta(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	data, err := h.Service.GetAbsensiPeserta(c.Context(), userID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiPesertaByID godoc
// @Summary Get detail absensi peserta
// @Description Get detail absensi milik peserta yang sedang login
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Absensi Peserta"
// @Router /absensi/peserta/me/{id_absensi} [get]
func (h *Handler) GetAbsensiPesertaByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	absensiID, err := c.ParamsInt("id_absensi")
	if err != nil {
		return response.Error(c, 400, "invalid absensi id", "BAD_REQUEST", nil)
	}

	data, err := h.Service.GetAbsensiPesertaByID(c.Context(), userID, absensiID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	if data == nil {
		return response.Error(c, 404, "absensi tidak ditemukan", "NOT_FOUND", nil)
	}

	return response.Success(c, data)
}