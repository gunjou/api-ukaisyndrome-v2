package materi

import (
	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}


// GetMateriPeserta godoc
// @Summary Get materi peserta
// @Description Get list materi berdasarkan modul
// @Tags Materi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_modul path int true "ID Modul"
// @Param type query string false "Filter tipe materi (video | document)"
// @Router /materi/peserta/{id_modul} [get]
func (h *Handler) GetMateriPeserta(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	modulID, err := c.ParamsInt("id_modul")
	if err != nil {
		return response.Error(c, 400, "invalid modul id", "BAD_REQUEST", nil)
	}

	// query param
	materiType := c.Query("type")

	var typePtr *string
	if materiType != "" {
		typePtr = &materiType
	}

	data, err := h.Service.GetMateriPeserta(c.Context(), userID, modulID, typePtr)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}

// GetMateriPrivatePeserta godoc
// @Summary Get materi private peserta
// @Description Get list materi private berdasarkan modul
// @Tags Materi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter tipe materi (video | document)"
// @Router /materi/peserta/private [get]
func (h *Handler) GetMateriPrivatePeserta(c *fiber.Ctx) error {

	userID := c.Locals("sub").(int)

	// optional query
	materiType := c.Query("type")

	var typePtr *string
	if materiType != "" {
		typePtr = &materiType
	}

	data, err := h.Service.GetMateriPrivatePeserta(
		c.Context(),
		userID,
		typePtr,
	)

	if err != nil {
		return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
	}

	return response.Success(c, data)
}