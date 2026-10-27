package jadwal

import (
	"errors"

	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

// GetJadwalPeserta godoc
// @Summary Get jadwal peserta
// @Description Get daftar jadwal aktif peserta. ONLINE berdasarkan paket kelas, OFFLINE berdasarkan kelompok belajar.
// @Tags Jadwal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /jadwal [get]
func (h *Handler) GetJadwalPeserta(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	data, err := h.Service.GetJadwalPeserta(c.Context(), userID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}

// GetJadwalPesertaByID godoc
// @Summary Get detail jadwal peserta
// @Description Get detail jadwal aktif milik peserta (ONLINE via paket kelas, OFFLINE via kelompok belajar)
// @Tags Jadwal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Jadwal"
// @Router /jadwal/{id} [get]
func (h *Handler) GetJadwalPesertaByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	jadwalID, err := c.ParamsInt("id")
	if err != nil {
		return response.Error(c, 400, "invalid jadwal id", "BAD_REQUEST", nil)
	}

	data, err := h.Service.GetJadwalPesertaByID(c.Context(), userID, jadwalID)
	if err != nil {
		if errors.Is(err, ErrJadwalNotFound) {
			return response.Error(c, 404, "jadwal tidak ditemukan", "NOT_FOUND", nil)
		}
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}